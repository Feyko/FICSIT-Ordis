package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain/domain"
	"context"
)

// AnalyseFileURL is the resolver for the analyseFileURL field.
func (r *queryResolver) AnalyseFileURL(ctx context.Context, fileURL string) (*domain.AnalysisResult, error) {
	return r.O.Analysis.AnalyseFileURL(ctx, fileURL)
}

// AnalyseText is the resolver for the analyseText field.
func (r *queryResolver) AnalyseText(ctx context.Context, text string) (*domain.AnalysisResult, error) {
	return r.O.Analysis.AnalyseText(ctx, []byte(text))
}
