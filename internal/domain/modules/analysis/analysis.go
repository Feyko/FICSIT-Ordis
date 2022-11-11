package analysis

import (
	"FICSIT-Ordis/internal/domain/domain"
	"archive/zip"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type Config struct {
}

func New(conf Config) (*Module, error) {
	return &Module{}, nil
}

type Module struct {
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
		return m.AnalyseText(ctx, string(b))
	}

	return m.analyseZipFile(ctx, zipReader)
}

func (m *Module) AnalyseText(ctx context.Context, url string) (*domain.AnalysisResult, error) {
	return nil, nil
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
	if res1.Cl == nil {
		res1.Cl = res2.Cl
	}
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
