package ordis

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/modules/commands"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"github.com/pkg/errors"
)

type Ordis struct {
	Commands *commands.Module
}

func New(conf config.OrdisConfig) (Ordis, error) {
	repo, err := arango.New[id.IDer](conf.Arango)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the repository")
	}

	commandsModule, err := commands.New(conf.Commands, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the commands module")
	}
	return Ordis{
		Commands: commandsModule,
	}, nil
}
