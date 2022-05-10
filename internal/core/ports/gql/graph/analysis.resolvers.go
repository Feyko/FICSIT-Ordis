package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/core/ports/gql/graph/model"
	"context"
	"fmt"
)

func (r *queryResolver) AnalyseFile(ctx context.Context, fileURL string) (*model.AnalysisResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AnalyseText(ctx context.Context, text string) (*model.AnalysisResult, error) {
	panic(fmt.Errorf("not implemented"))
}
