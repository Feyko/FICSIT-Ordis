package base

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"github.com/pkg/errors"
)

type Config struct {
	Auth *auth.Module

	CreatePermission *domain.Permission
	GetPermission    *domain.Permission
	ListPermission   *domain.Permission
	DeletePermission *domain.Permission
	UpdatePermission *domain.Permission
}

func New[E id.IDer](conf Config, collection repo.Collection[E]) *Module[E] {
	return &Module[E]{
		Auth:       conf.Auth,
		Collection: collection,
		config:     conf,
	}
}

func NewDefaultConfig(authModule *auth.Module) Config {
	perm := domain.PermissionContentEditing
	return Config{
		Auth:             authModule,
		CreatePermission: &perm,
		DeletePermission: &perm,
		UpdatePermission: &perm,
	}
}

func NewDefaultConfigNoPerm(authModule *auth.Module) Config {
	return Config{
		Auth: authModule,
	}
}

type Module[E id.IDer] struct {
	Auth       *auth.Module
	Collection repo.Collection[E]
	config     Config
}

func (mod *Module[E]) Create(ctx context.Context, element E) error {
	err := mod.checkForPermission(&ctx, mod.config.CreatePermission)
	if err != nil {
		return err
	}
	err = mod.Collection.Create(ctx, element)
	if err != nil {
		return errors.Wrap(err, "could not create a new element")
	}
	return nil
}

func (mod *Module[E]) Get(ctx context.Context, ID string) (E, error) {
	err := mod.checkForPermission(&ctx, mod.config.GetPermission)
	if err != nil {
		return *new(E), err
	}
	cmd, err := mod.Collection.Get(ctx, ID)
	if err != nil {
		return *new(E), errors.Wrapf(err, "could not get the command with ID '%v'", ID)
	}
	return cmd, nil
}

func (mod *Module[E]) List(ctx context.Context) ([]E, error) {
	err := mod.checkForPermission(&ctx, mod.config.ListPermission)
	if err != nil {
		return nil, err
	}
	elems, err := mod.Collection.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get all the elements")
	}
	return elems, nil
}

func (mod *Module[E]) Delete(ctx context.Context, ID string) error {
	err := mod.checkForPermission(&ctx, mod.config.DeletePermission)
	if err != nil {
		return err
	}
	err = mod.Collection.Delete(ctx, ID)
	if err != nil {
		return errors.Wrap(err, "could not delete the element")
	}
	return nil
}

func (mod *Module[E]) Update(ctx context.Context, ID string, updateElement any) (E, error) {
	err := mod.checkForPermission(&ctx, mod.config.UpdatePermission)
	if err != nil {
		return *new(E), err
	}
	elem, err := mod.Collection.Update(ctx, ID, updateElement)
	if err != nil {
		return *new(E), errors.Wrap(err, "could not update the element")
	}
	return elem, nil
}

func (mod *Module[E]) checkForPermission(ctx *context.Context, perm *domain.Permission) error {
	if perm == nil {
		return nil
	}

	return errors.Wrap(mod.Auth.Authorize(ctx, *perm), "could not authorize")
}
