package core

import (
	"FICSIT-Ordis/internal/core/modules/commands"
)

type Ordis struct {
	Commands *commands.Module
}

func New() Ordis {
	return Ordis{
		Commands: commands.New(),
	}
}
