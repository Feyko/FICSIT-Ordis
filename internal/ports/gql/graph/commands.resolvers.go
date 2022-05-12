package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain"
	"FICSIT-Ordis/internal/ports/gql/graph/generated"
	"FICSIT-Ordis/internal/ports/gql/graph/model"
	"context"
	"fmt"
)

func (r *commandResolver) Response(ctx context.Context, obj *model.Command) (*model.Response, error) {
	resp := model.Response(obj.Response)
	return &resp, nil
}

func (r *mutationResolver) CreateCommand(ctx context.Context, command model.CommandCreation) (*model.Command, error) {
	cmd := domain.Command{
		Name:     command.Name,
		Aliases:  command.Aliases,
		Response: domain.Response(*command.Response),
	}
	err := r.o.Commands.Create(cmd)
	return model.ModelCommand(cmd), err
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

// Command returns generated.CommandResolver implementation.
func (r *Resolver) Command() generated.CommandResolver { return &commandResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type commandResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
