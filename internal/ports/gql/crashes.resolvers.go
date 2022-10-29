package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/ports/gql/model"
	"context"
	"fmt"
)

// CreateCrash is the resolver for the createCrash field.
func (r *mutationResolver) CreateCrash(ctx context.Context, crash model.CrashCreation) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

// UpdateCrash is the resolver for the updateCrash field.
func (r *mutationResolver) UpdateCrash(ctx context.Context, crash model.CrashUpdate) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

// DeleteCrash is the resolver for the deleteCrash field.
func (r *mutationResolver) DeleteCrash(ctx context.Context, name string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// ListAllCrashes is the resolver for the listAllCrashes field.
func (r *queryResolver) ListAllCrashes(ctx context.Context) ([]model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

// FindCrash is the resolver for the findCrash field.
func (r *queryResolver) FindCrash(ctx context.Context, name string) (*model.Crash, error) {
	panic(fmt.Errorf("not implemented"))
}

// CrashAnalysis is the resolver for the crashAnalysis field.
func (r *queryResolver) CrashAnalysis(ctx context.Context, text string) ([]model.CrashMatch, error) {
	panic(fmt.Errorf("not implemented"))
}
