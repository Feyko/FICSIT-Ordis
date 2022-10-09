package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/ports/gql/graph/generated"
	"FICSIT-Ordis/internal/ports/gql/graph/model"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

// CreateCommand is the resolver for the createCommand field.
func (r *mutationResolver) CreateCommand(ctx context.Context, command model.CommandCreation) (*domain.Command, error) {
	cmd := domain.Command(command)
	err := r.O.Commands.Create(nil, cmd)
	return &cmd, err
}

// UpdateCommand is the resolver for the updateCommand field.
func (r *mutationResolver) UpdateCommand(ctx context.Context, name string, command domain.CommandUpdate) (*domain.Command, error) {
	cmd, err := r.O.Commands.Get(name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find command '%v'", name)
	}
	newCommand, err := r.O.Commands.Update(nil, cmd.ID(), command)
	if err != nil {
		return nil, errors.Wrap(err, "could not update the command")
	}

	return &newCommand, nil
}

// DeleteCommand is the resolver for the deleteCommand field.
func (r *mutationResolver) DeleteCommand(ctx context.Context, name string) (bool, error) {
	err := r.O.Commands.Delete(nil, name)
	return err == nil, err
}

// ListAllCommands is the resolver for the listAllCommands field.
func (r *queryResolver) ListAllCommands(ctx context.Context) ([]domain.Command, error) {
	return r.O.Commands.List(nil)
}

// FindCommand is the resolver for the findCommand field.
func (r *queryResolver) FindCommand(ctx context.Context, name string) (*domain.Command, error) {
	cmd, err := r.O.Commands.Get(name)
	return cmd, err
}

// ExecuteCommand is the resolver for the executeCommand field.
func (r *queryResolver) ExecuteCommand(ctx context.Context, text string) (*domain.Response, error) {
	return r.O.Commands.Execute(text)
}

// Response is the resolver for the response field.
func (r *commandCreationResolver) Response(ctx context.Context, obj *model.CommandCreation, data *model.ResponseInput) error {
	obj.Response = domain.Response(*data)
	return nil
}

// Response is the resolver for the response field.
func (r *commandUpdateResolver) Response(ctx context.Context, obj *domain.CommandUpdate, data *model.ResponseInput) error {
	panic(fmt.Errorf("not implemented: Response - response"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// CommandCreation returns generated.CommandCreationResolver implementation.
func (r *Resolver) CommandCreation() generated.CommandCreationResolver {
	return &commandCreationResolver{r}
}

// CommandUpdate returns generated.CommandUpdateResolver implementation.
func (r *Resolver) CommandUpdate() generated.CommandUpdateResolver { return &commandUpdateResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type commandCreationResolver struct{ *Resolver }
type commandUpdateResolver struct{ *Resolver }
