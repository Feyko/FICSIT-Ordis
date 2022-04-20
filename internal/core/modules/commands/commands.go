package commands

import (
	"FICSIT-Ordis/internal/core/modules/base"
	"FICSIT-Ordis/internal/core/ports/repositories"
)

type Command struct {
	Name,
	Response,
	Media string
}

func (elem Command) ID() string {
	return elem.Name
}

type CommandsModule base.BasicModule[Command]

func New() *CommandsModule {
	return &CommandsModule{
		Repository: repositories.MemoryRepository[Command]{},
	}
}
