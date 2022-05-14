package commands

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/translators"
	"fmt"
)

func New(conf config.CommandsConfig, repo repos.Repository) (*Module, error) {
	collection, err := repo.GetCollection("Commands")
	if err != nil {
		return nil, fmt.Errorf("could not get the collection: %w", err)
	}
	translator := translators.Wrap[domain.Command](collection)
	return &Module{
		*base.NewSearchable[domain.Command](translator),
	}, nil
}

type Module struct {
	base.Searchable[domain.Command]
}

//func (m *Module) Execute(text string) (domain.Response, error) {
//	strings.Split(text, " ")
//	cmd, err := m.Get()
//}
//
//func textToArguments(text string) []string {
//	quotingRegexp := regexp.MustCompile(`".*"`)
//}
