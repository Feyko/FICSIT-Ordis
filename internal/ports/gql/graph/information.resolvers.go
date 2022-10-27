package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

// SetLatestInformation is the resolver for the setLatestInformation field.
func (r *mutationResolver) SetLatestInformation(ctx context.Context, info string) (bool, error) {
	err := r.O.Information.Set(ctx, info)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLatestInformation is the resolver for the getLatestInformation field.
func (r *queryResolver) GetLatestInformation(ctx context.Context) (*string, error) {
	info, err := r.O.Information.Get(ctx)
	if err != nil {
		return nil, err
	}
	return info, nil
}
