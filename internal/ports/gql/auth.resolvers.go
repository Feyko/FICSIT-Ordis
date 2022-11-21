package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

// CreateNewToken is the resolver for the createNewToken field.
func (r *mutationResolver) CreateNewToken(ctx context.Context, roleIDs []int) (string, error) {
	token, err := r.O.Auth.NewToken(ctx, roleIDs...)
	return token.String, err
}

// RevokeToken is the resolver for the revokeToken field.
func (r *mutationResolver) RevokeToken(ctx context.Context, tokenID string) (bool, error) {
	err := r.O.Auth.RevokeTokenID(ctx, tokenID)
	return err == nil, err
}

// RevokeCurrentToken is the resolver for the revokeCurrentToken field.
func (r *mutationResolver) RevokeCurrentToken(ctx context.Context) (bool, error) {
	err := r.O.Auth.RevokeCurrentToken(ctx)
	return err == nil, err
}
