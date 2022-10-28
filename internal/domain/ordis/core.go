package ordis

import (
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/commands"
	"FICSIT-Ordis/internal/domain/modules/information"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"github.com/pkg/errors"
)

type Ordis struct {
	Commands    *commands.Module
	Auth        *auth.Module
	Information *information.Module
}

type Config struct {
	Arango      arango.Config
	Auth        auth.Config
	Commands    commands.Config
	Information information.Config
}

func New(conf Config) (Ordis, error) {
	authModule, err := auth.New(conf.Auth)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the auth module")
	}

	repo, err := arango.New[id.IDer](conf.Arango)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the repository")
	}

	fillAuthConfig(authModule, &conf.Commands.AuthedConfig)

	commandsModule, err := commands.New(conf.Commands, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the commands module")
	}

	fillAuthConfig(authModule, &conf.Information.AuthedConfig)

	infoModule, err := information.New(conf.Information, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the information module")
	}

	return Ordis{
		Commands:    commandsModule,
		Auth:        authModule,
		Information: infoModule,
	}, nil
}

func fillAuthConfig(module *auth.Module, conf *auth.AuthedConfig) {
	if conf.AuthModule == nil {
		conf.AuthModule = module
	}
}
