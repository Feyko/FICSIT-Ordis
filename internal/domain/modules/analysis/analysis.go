package analysis

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/crashes"
	"FICSIT-Ordis/internal/smr"
	"archive/zip"
	"bytes"
	"context"
	"git.sr.ht/~emersion/gqlclient"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	CrashesModule *crashes.Module
}

func New(conf Config) (*Module, error) {
	return &Module{
		crashesModule: conf.CrashesModule,
		smr:           gqlclient.New("https://api.ficsit.app/v2/query", http.DefaultClient),
	}, nil
}

type Module struct {
	crashesModule *crashes.Module
	smr           *gqlclient.Client
}

func (m *Module) AnalyseFileURL(ctx context.Context, url string) (*domain.AnalysisResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "error downloading the file")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error downloading the file")
	}

	buf := bytes.NewBuffer(b)

	return m.AnalyseFile(ctx, buf)
}

func (m *Module) AnalyseFile(ctx context.Context, file io.Reader) (*domain.AnalysisResult, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}

	reader := bytes.NewReader(b)

	zipReader, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return m.AnalyseText(ctx, b)
	}

	return m.analyseZipFile(ctx, zipReader)
}

func (m *Module) AnalyseText(ctx context.Context, text []byte) (*domain.AnalysisResult, error) {
	extractor := newLogExtractor(text, m)
	return extractor.Result(ctx)
}

func (m *Module) analyseZipFile(ctx context.Context, zipFile *zip.Reader) (*domain.AnalysisResult, error) {
	var result domain.AnalysisResult

	for _, file := range zipFile.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, errors.Wrapf(err, "error opening zip subfile %q", file.Name)
		}
		newResult, err := m.AnalyseFile(ctx, fileReader)
		if newResult != nil {
			err = mergeResults(&result, newResult)
			if err != nil {
				return nil, errors.Wrap(err, "error merging results")
			}
		}
	}
	return &result, nil
}

// mergeResults merges res2 into res1
// Crash matches are merged together instead of one list taking precedence. Duplicate matches are removed
func mergeResults(res1, res2 *domain.AnalysisResult) error {
	if res1.CommandLine == nil {
		res1.CommandLine = res2.CommandLine
	}
	if res1.DesiredSMLVersion == nil {
		res1.DesiredSMLVersion = res2.DesiredSMLVersion
	}
	if res1.GameVersion == nil {
		res1.GameVersion = res2.GameVersion
	}
	if res1.ModList == nil {
		res1.ModList = res2.ModList
	}
	if res1.Path == nil {
		res1.Path = res2.Path
	}
	if res1.PiracyInfo == nil {
		res1.PiracyInfo = res2.PiracyInfo
	}
	if res1.SMLVersion == nil {
		res1.SMLVersion = res2.SMLVersion
	}

	res1.CrashMatches = mergeCrashMatches(res1.CrashMatches, res2.CrashMatches)

	return nil
}

func mergeCrashMatches(m1, m2 []domain.CrashMatch) []domain.CrashMatch {
	r := m1
	for _, match := range m2 {
		for _, other := range m1 {
			if match.Crash.Name == other.Crash.Name {
				continue
			}
		}
		r = append(r, match)
	}
	return r
}

func newLogExtractor(text []byte, module *Module) logExtractor {
	return logExtractor{
		module: module,
		text:   text,
	}
}

type logExtractor struct {
	module *Module
	text   []byte
}

func (l *logExtractor) Result(ctx context.Context) (*domain.AnalysisResult, error) {
	var result domain.AnalysisResult

	matches, err := l.module.crashesModule.Analyse(ctx, string(l.text))
	if err != nil {
		return nil, errors.Wrap(err, "error analysing crashes")
	}
	result.CrashMatches = matches

	modList, err := l.ModList()
	if err != nil {
		return nil, errors.Wrap(err, "error finding mod list")
	}
	result.ModList = modList

	l.setStringForMatch(&result.SMLVersion, `Satisfactory Mod Loader v\.?(\d(\.\d)*)`, 1)
	l.setStringForMatch(&result.GameVersion, `Net CL: (\d+)`, 1)
	l.setStringForMatch(&result.Path, `(?m)LogInit: Base Directory: (.*)$`, 1)
	l.setStringForMatch(&result.CommandLine, `(?m)LogInit: Command Line: (.*)$`, 1)
	l.setStringForMatch(&result.LauncherID, `(?m)LogInit: Launcher ID: (.*)$`, 1)
	l.setStringForMatch(&result.LauncherArtifact, `(?m)LogInit: Launcher Artifact: (.*)$`, 1)

	result.DesiredSMLVersion = l.desiredSMLVersion(result.GameVersion)

	return &result, nil
}

func (l *logExtractor) ModList() ([]string, error) {
	regex := regexp.MustCompile(`LogPakFile: New pak file ../../../FactoryGame/Mods/(.*?)/`)

	matches := regex.FindAllSubmatch(l.text, -1)
	r := make([]string, len(matches))
	for i, match := range matches {
		r[i] = string(match[1])
	}

	return r, nil
}

func (l *logExtractor) setStringForMatch(p **string, regex string, group int) {
	reg := regexp.MustCompile(regex)
	found := reg.FindSubmatch(l.text)
	if len(found) > 0 && len(found[0]) != 0 {
		s := strings.TrimSpace(string(found[group]))
		*p = &s
	}
}

func (l *logExtractor) desiredSMLVersion(gameVersion *string) *string {
	if gameVersion == nil {
		return nil
	}

	gameCL, err := strconv.ParseInt(*gameVersion, 10, 32)
	if err != nil {
		return nil
	}

	r, err := smr.QGetSMLVersions(l.module.smr, context.Background())
	if err != nil {
		return nil
	}
	gameCL = gameCL
	for _, version := range r.Sml_versions {
		if version.Satisfactory_version > int32(gameCL) {
			continue
		}
		return &version.Version
	}

	return nil
}
