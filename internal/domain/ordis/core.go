package ordis

import (
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/commands"
	"FICSIT-Ordis/internal/domain/modules/crashes"
	"FICSIT-Ordis/internal/domain/modules/latestInformation"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"github.com/pkg/errors"
)

type Ordis struct {
	Commands          *commands.Module
	Crashes           *crashes.Module
	Auth              *auth.Module
	LatestInformation *latestInformation.Module
}

type Config struct {
	Arango            arango.Config
	Auth              auth.Config
	Commands          commands.Config
	Crashes           crashes.Config
	LatestInformation latestInformation.Config
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

	fillAuthConfig(authModule, &conf.Crashes.AuthedConfig)

	crashesModule, err := crashes.New(conf.Crashes, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the crashes module")
	}

	fillAuthConfig(authModule, &conf.LatestInformation.AuthedConfig)

	infoModule, err := latestInformation.New(conf.LatestInformation, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the latestInformation module")
	}

	return Ordis{
		Commands:          commandsModule,
		Crashes:           crashesModule,
		Auth:              authModule,
		LatestInformation: infoModule,
	}, nil
}

func fillAuthConfig(module *auth.Module, conf *auth.AuthedConfig) {
	if conf.AuthModule == nil {
		conf.AuthModule = module
	}
}
