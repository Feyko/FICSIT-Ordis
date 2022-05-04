package commands

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/base"
	"FICSIT-Ordis/internal/core/ports/repositories"
	"FICSIT-Ordis/internal/core/ports/repositories/arango"
	"FICSIT-Ordis/internal/core/ports/repositories/memrepo"
	"fmt"
)

func New(conf config.CommandsConfig) (*Module, error) {
	repo, err := newRepo(conf)
	if err != nil {
		return nil, fmt.Errorf("could not create a repository: %w", err)
	}
	return &Module{
		*base.New[Command](repo),
	}, nil
}

type Command struct {
	Name,
	Response,
	Media string
}

func (elem Command) Type() string {
	return "Command"
}

func (elem Command) ID() string {
	return elem.Name
}

type Module struct {
	base.BasicModule[Command]
}

type Repo repositories.Repository[Command]

func newRepo(conf config.CommandsConfig) (Repo, error) {
	if conf.Persistent {
		return newPersistentRepo(conf)
	}
	return newMemRepo(), nil
}

func newPersistentRepo(conf config.CommandsConfig) (Repo, error) {
	return arango.New[Command](conf.Arango, "Commands")
}

func newMemRepo() Repo {
	return memrepo.New[Command]()
}
