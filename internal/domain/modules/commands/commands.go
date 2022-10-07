package commands

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"github.com/pkg/errors"
	"strings"
)

func New[T id.IDer](conf config.CommandsConfig, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetCollection[domain.Command](repository, "Commands")
	notFound := errors.Is(err, repo.ErrCollectionNotFound)
	if notFound {
		collection, err = repos.CreateCollection[domain.Command](repository, "Commands")
		if err != nil {
			return nil, errors.Wrap(err, "could not create the collection")
		}
	}
	if err != nil && !notFound {
		return nil, errors.Wrap(err, "could not get the collection")
	}
	return &Module{
		*base.NewSearchable[domain.Command](collection),
	}, nil
}

type Module struct {
	base.Searchable[domain.Command]
}

func (m *Module) Execute(text string) (*domain.Response, error) {
	first, _, _ := strings.Cut(text, " ")
	cmd, err := m.Get(first)
	if err != nil {
		return nil, errors.Wrap(err, "error getting the command")
	}
	return &cmd.Response, nil
}

func (m *Module) Get(name string) (*domain.Command, error) {
	cmds, err := m.Search(name)
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
