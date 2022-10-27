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
	Commands    commands.Config
	Information information.Config
}

func New(conf Config) (Ordis, error) {
	authModule, err := auth.New(auth.Config{
		Secret: "test-secret",
	})
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the auth module")
	}

	repo, err := arango.New[id.IDer](conf.Arango)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the repository")
	}

	commandsConfig := conf.Commands
	if commandsConfig.Auth == nil {
		commandsConfig.Auth = authModule
	}

	commandsModule, err := commands.New(commandsConfig, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the commands module")
	}

	infoConfig := conf.Information
	if infoConfig.Auth == nil {
		infoConfig.Auth = authModule
	}

	infoModule, err := information.New(infoConfig, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the commands module")
	}

	return Ordis{
		Commands:    commandsModule,
		Auth:        authModule,
		Information: infoModule,
	}, nil
}
