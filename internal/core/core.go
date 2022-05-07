package core

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/commands"
	"FICSIT-Ordis/internal/core/ports/repos/arango"
	"fmt"
)

type Ordis struct {
	Commands *commands.Module
}

func New(conf config.OrdisConfig) (Ordis, error) {
	repo, err := arango.New(conf.Arango)
	if err != nil {
		return Ordis{}, fmt.Errorf("could not create the repository: %w", err)
	}

	commandsModule, err := commands.New(conf.Commands, repo)
	if err != nil {
		return Ordis{}, fmt.Errorf("could not create the commands module: %w", err)
	}
	return Ordis{
		Commands: commandsModule,
	}, nil
}
