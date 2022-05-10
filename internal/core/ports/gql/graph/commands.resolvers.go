package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/core/ports/gql/graph/generated"
	"FICSIT-Ordis/internal/core/ports/gql/graph/model"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateCommand(ctx context.Context, command model.CommandCreation) (*model.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCommand(ctx context.Context, name string, command model.CommandUpdate) (*model.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCommand(ctx context.Context, name string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListAllCommands(ctx context.Context) ([]*model.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindCommand(ctx context.Context, name string) (*model.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TryForCommand(ctx context.Context, text string) (*model.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
