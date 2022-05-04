package core

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/commands"
	"fmt"
)

type Ordis struct {
	Commands *commands.Module
}

func New(conf config.OrdisConfig) (Ordis, error) {
	commandsModule, err := commands.New(conf.Commands)
	if err != nil {
		return Ordis{}, fmt.Errorf("could not create the commands module: %w", err)
	}
	return Ordis{
		Commands: commandsModule,
	}, nil
}
