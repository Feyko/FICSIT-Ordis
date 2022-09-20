package commands

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"fmt"
	"github.com/mattn/go-shellwords"
)

func New[T id.IDer](conf config.CommandsConfig, repo repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetCollection[domain.Command](repo, "Commands")
	if err != nil {
		return nil, fmt.Errorf("could not get the collection: %w", err)
	}
	return &Module{
		*base.NewSearchable[domain.Command, domain.CommandUpdate](collection),
	}, nil
}

type Module struct {
	base.Searchable[domain.Command]
}

func (m *Module) Execute(text string) (*domain.Response, error) {
	args, err := shellwords.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("could not parse the input text: %w", err)
	}
	cmd, err := m.Get(args[0])
	if err != nil {
		return nil, nil
	}
	return &cmd.Response, nil
}
