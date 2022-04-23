package core

import (
	"FICSIT-Ordis/internal/core/modules/commands"
	"FICSIT-Ordis/internal/core/ports/repositories"
)

type Ordis struct {
	Commands *commands.Module
}

func New() Ordis {
	return Ordis{
		Commands: commands.New(new(repositories.MemoryRepository[commands.Command])),
	}
}
