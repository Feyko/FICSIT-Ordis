package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/ports/gql/model"
	"context"
	"fmt"
)

// AnalyseFile is the resolver for the analyseFile field.
func (r *queryResolver) AnalyseFile(ctx context.Context, fileURL string) (*model.AnalysisResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// AnalyseText is the resolver for the analyseText field.
func (r *queryResolver) AnalyseText(ctx context.Context, text string) (*model.AnalysisResult, error) {
	panic(fmt.Errorf("not implemented"))
}
