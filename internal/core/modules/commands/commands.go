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

type Module struct {
	base.BasicModule[Command]
}

func New(repo repositories.Repository[Command]) *Module {
	return &Module{
		base.BasicModule[Command]{
			Repository: repo,
		},
	}
}
