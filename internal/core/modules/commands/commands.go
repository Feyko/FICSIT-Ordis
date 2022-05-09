package commands

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/base"
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/core/ports/repos/translators"
	"fmt"
)

func New(conf config.CommandsConfig, repo repos.Repository) (*Module, error) {
	collection, err := repo.GetCollection("Commands")
	if err != nil {
		return nil, fmt.Errorf("could not get the collection: %w", err)
	}
	translator := translators.Wrap[Command](collection)
	return &Module{
		*base.NewSearchable[Command](translator),
	}, nil
}

type Command struct {
	Name,
	Response,
	Media string
}

func (elem Command) ID() string {
	return elem.Name
}

func (elem Command) SearchFields() []string {
	return []string{"Name", "Response"}
}

type Module struct {
	base.Searchable[Command]
}
