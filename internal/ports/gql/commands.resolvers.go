package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/ports/gql/generated"
	"FICSIT-Ordis/internal/ports/gql/model"
	"context"

	"github.com/pkg/errors"
)

// CreateCommand is the resolver for the createCommand field.
func (r *mutationResolver) CreateCommand(ctx context.Context, command model.CommandCreation) (*domain.Command, error) {
	cmd := domain.Command(command)
	err := r.O.Commands.Create(ctx, cmd)
	return &cmd, err
}

// UpdateCommand is the resolver for the updateCommand field.
func (r *mutationResolver) UpdateCommand(ctx context.Context, name string, command domain.CommandUpdate) (*domain.Command, error) {
	cmd, err := r.O.Commands.Get(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find command '%v'", name)
	}
	newCommand, err := r.O.Commands.Update(ctx, cmd.ID(), command)
	if err != nil {
		return nil, errors.Wrap(err, "could not update the command")
	}

	return newCommand, nil
}

// DeleteCommand is the resolver for the deleteCommand field.
func (r *mutationResolver) DeleteCommand(ctx context.Context, name string) (bool, error) {
	err := r.O.Commands.Delete(ctx, name)
	return err == nil, err
}

// ListAllCommands is the resolver for the listAllCommands field.
func (r *queryResolver) ListAllCommands(ctx context.Context) ([]domain.Command, error) {
	return r.O.Commands.List(ctx)
}

// FindCommand is the resolver for the findCommand field.
func (r *queryResolver) FindCommand(ctx context.Context, name string) (*domain.Command, error) {
	cmd, err := r.O.Commands.Get(ctx, name)
	return cmd, err
}

// ExecuteCommand is the resolver for the executeCommand field.
func (r *queryResolver) ExecuteCommand(ctx context.Context, text string) (*domain.Response, error) {
	return r.O.Commands.Execute(ctx, text)
}

// Response is the resolver for the response field.
func (r *commandCreationResolver) Response(ctx context.Context, obj *model.CommandCreation, data *model.ResponseInput) error {
	obj.Response = domain.Response(*data)
	return nil
}

// Response is the resolver for the response field.
func (r *commandUpdateResolver) Response(ctx context.Context, obj *domain.CommandUpdate, data *model.ResponseInput) error {
	obj.Response = &domain.ResponseUpdate{
		Text:       data.Text,
		MediaLinks: data.MediaLinks,
	}
	return nil
}

// CommandCreation returns generated.CommandCreationResolver implementation.
func (r *Resolver) CommandCreation() generated.CommandCreationResolver {
	return &commandCreationResolver{r}
}

// CommandUpdate returns generated.CommandUpdateResolver implementation.
func (r *Resolver) CommandUpdate() generated.CommandUpdateResolver { return &commandUpdateResolver{r} }

type commandCreationResolver struct{ *Resolver }
type commandUpdateResolver struct{ *Resolver }
