package commands

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"github.com/pkg/errors"
	"strings"
)

type Config struct {
	Auth   *auth.Module
	NoAuth bool
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[domain.Command](repository, "Commands")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}
	baseConf := base.NewDefaultConfig(conf.Auth)
	if conf.NoAuth {
		baseConf = base.NewDefaultConfigNoPerm(nil)
	}
	return &Module{
		*base.NewSearchable[domain.Command](baseConf, collection),
	}, nil
}

type Module struct {
	base.Searchable[domain.Command]
}

func (m *Module) Execute(ctx context.Context, text string) (*domain.Response, error) {
	first, _, _ := strings.Cut(text, " ")
	cmd, err := m.Get(ctx, first)
	if err != nil {
		return nil, errors.Wrap(err, "error getting the command")
	}
	return &cmd.Response, nil
}

func (m *Module) Get(ctx context.Context, name string) (*domain.Command, error) {
	cmds, err := m.Search(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "error searching")
	}
	if len(cmds) == 0 {
		return nil, errors.Errorf("command '%v' not found", name)
	}
	if len(cmds) > 1 {
		return nil, errors.Errorf("command '%v' has multiple entries", name)
	}
	return &cmds[0], nil
}

func (m *Module) Create(ctx context.Context, command domain.Command) error {
	err := m.ensureDoesntExist(ctx, command)
	if err != nil {
		return err
	}

	return m.Module.Create(ctx, command)
}

func (m *Module) ensureDoesntExist(ctx context.Context, cmd domain.Command) error {
	err := m.checkCommandNameDoesntExist(ctx, cmd.Name)
	if err != nil {
		return err
	}
	for _, alias := range cmd.Aliases {
		err = m.checkCommandNameDoesntExist(ctx, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) checkCommandNameDoesntExist(ctx context.Context, name string) error {
	found, err := m.Search(ctx, name)
	if err != nil {
		return errors.Wrap(err, "error checking existing commands")
	}
	if len(found) > 0 {
		return errors.Errorf("command %v already exists", name)
	}
	return nil
}