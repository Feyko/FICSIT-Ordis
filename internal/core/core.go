package core

import "FICSIT-Ordis/internal/core/commands"

type Ordis struct {
	Commands commands.CommandsModule
}

func New() Ordis {
	return Ordis{
		Commands: commands.NewModule(),
	}
}
