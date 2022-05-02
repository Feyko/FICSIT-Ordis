package core

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/commands"
)

type Ordis struct {
	Commands *commands.Module
}

func New(conf config.OrdisConfig) Ordis {
	return Ordis{
		Commands: commands.New(conf.Commands),
	}
}
