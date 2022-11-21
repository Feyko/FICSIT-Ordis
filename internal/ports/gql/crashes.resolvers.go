package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/ports/gql/generated"
	"FICSIT-Ordis/internal/ports/gql/model"
	"context"
)

// CreateCrash is the resolver for the createCrash field.
func (r *mutationResolver) CreateCrash(ctx context.Context, crash model.CrashCreation) (bool, error) {
	err := r.O.Crashes.Create(ctx, domain.Crash(crash))
	return err == nil, err
}

// UpdateCrash is the resolver for the updateCrash field.
func (r *mutationResolver) UpdateCrash(ctx context.Context, name string, crash domain.CrashUpdate) (*domain.Crash, error) {
	newCrash, err := r.O.Crashes.Update(ctx, name, crash)
	return &newCrash, err
}

// DeleteCrash is the resolver for the deleteCrash field.
func (r *mutationResolver) DeleteCrash(ctx context.Context, name string) (bool, error) {
	err := r.O.Crashes.Delete(ctx, name)
	return err == nil, err
}

// ListAllCrashes is the resolver for the listAllCrashes field.
func (r *queryResolver) ListAllCrashes(ctx context.Context) ([]domain.Crash, error) {
	return r.O.Crashes.List(ctx)
}

// FindCrash is the resolver for the findCrash field.
func (r *queryResolver) FindCrash(ctx context.Context, name string) (*domain.Crash, error) {
	crash, err := r.O.Crashes.Get(ctx, name)
	return &crash, err
}

// CrashAnalysis is the resolver for the crashAnalysis field.
func (r *queryResolver) CrashAnalysis(ctx context.Context, text string) ([]domain.CrashMatch, error) {
	return r.O.Crashes.Analyse(ctx, text)
}

// SearchCrashes is the resolver for the searchCrashes field.
func (r *queryResolver) SearchCrashes(ctx context.Context, search string) ([]domain.Crash, error) {
	return r.O.Crashes.Search(ctx, search)
}

// Response is the resolver for the response field.
func (r *crashCreationResolver) Response(ctx context.Context, obj *model.CrashCreation, data *model.ResponseInput) error {
	obj.Response = domain.Response(*data)
	return nil
}

// Response is the resolver for the response field.
func (r *crashUpdateResolver) Response(ctx context.Context, obj *domain.CrashUpdate, data *model.ResponseInput) error {
	obj.Response = (*domain.Response)(data)
	return nil
}

// CrashCreation returns generated.CrashCreationResolver implementation.
func (r *Resolver) CrashCreation() generated.CrashCreationResolver { return &crashCreationResolver{r} }

// CrashUpdate returns generated.CrashUpdateResolver implementation.
func (r *Resolver) CrashUpdate() generated.CrashUpdateResolver { return &crashUpdateResolver{r} }

type crashCreationResolver struct{ *Resolver }
type crashUpdateResolver struct{ *Resolver }
