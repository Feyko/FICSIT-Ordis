package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/core/ports/gql/graph/model"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateCrash(ctx context.Context, crash model.CrashCreation) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCrash(ctx context.Context, crash model.CrashUpdate) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCrash(ctx context.Context, name string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListAllCrashes(ctx context.Context) ([]*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindCrash(ctx context.Context, name string) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CrashAnalysis(ctx context.Context, text string) ([]*model.CrashMatch, error) {
	panic(fmt.Errorf("not implemented"))
}
