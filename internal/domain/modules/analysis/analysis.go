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

func (m *Module) AnalyseFileURL(ctx context.Context, url string) (domain.AnalysisResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return domain.AnalysisResult{}, errors.Wrap(err, "error downloading the file")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AnalysisResult{}, errors.Wrap(err, "error downloading the file")
	}

	buf := bytes.NewReader(b)

	return m.AnalyseFile(ctx, buf, int64(len(b)))
}

func (m *Module) AnalyseFile(ctx context.Context, file io.ReaderAt, size int64) (domain.AnalysisResult, error) {
	reader, err := zip.NewReader(file, size)
	if err != nil {
		return m.AnalyseText(ctx, string(b))
	}

	for _, file := range reader.File {
		file.Open()
	}
	return domain.AnalysisResult{}, nil
}

func (m *Module) AnalyseText(ctx context.Context, url string) (domain.AnalysisResult, error) {
	return domain.AnalysisResult{}, nil
}
