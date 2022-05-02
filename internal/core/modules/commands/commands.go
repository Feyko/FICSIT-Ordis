package commands

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/base"
	"FICSIT-Ordis/internal/core/ports/repositories"
	"FICSIT-Ordis/internal/core/ports/repositories/memrepo"
)

func New(conf config.CommandsConfig) *Module {
	repo := newRepo(conf)
	return &Module{
		*base.New[Command](repo),
	}
}

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

type Repo repositories.Repository[Command]

func newRepo(conf config.CommandsConfig) Repo {
	if conf.Persistent {
		return newPersistentRepo()
	}
	return newMemRepo()
}

func newPersistentRepo() Repo {
	return nil
}

func newMemRepo() Repo {
	return memrepo.New[Command]()
}
