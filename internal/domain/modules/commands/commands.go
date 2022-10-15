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
	Auth *auth.Module
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[domain.Command](repository, "Commands")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}
	return &Module{
		*base.NewSearchable[domain.Command](base.NewDefaultConfig(conf.Auth), collection),
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
