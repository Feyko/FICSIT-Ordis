package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain/domain"
	"context"
)

// SetLatestInformation is the resolver for the setLatestInformation field.
func (r *mutationResolver) SetLatestInformation(ctx context.Context, text string) (bool, error) {
	err := r.O.LatestInformation.Set(ctx, text)
	return err == nil, nil
}

// RemoveLatestInformation is the resolver for the removeLatestInformation field.
func (r *mutationResolver) RemoveLatestInformation(ctx context.Context) (bool, error) {
	err := r.O.LatestInformation.Remove(ctx)
	return err == nil, nil
}

// GetLatestInformation is the resolver for the getLatestInformation field.
func (r *queryResolver) GetLatestInformation(ctx context.Context) (*domain.LatestInformation, error) {
	return r.O.LatestInformation.Get(ctx)
}
