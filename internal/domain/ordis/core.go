package ordis

import (
	"FICSIT-Ordis/internal/domain/modules/analysis"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/commands"
	"FICSIT-Ordis/internal/domain/modules/crashes"
	"FICSIT-Ordis/internal/domain/modules/latestInformation"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"github.com/pkg/errors"
)

type Ordis struct {
	Auth              *auth.Module
	Analysis          *analysis.Module
	Commands          *commands.Module
	Crashes           *crashes.Module
	LatestInformation *latestInformation.Module
}

type Config struct {
	Analysis          analysis.Config
	Arango            arango.Config
	Auth              auth.Config
	Commands          commands.Config
	Crashes           crashes.Config
	LatestInformation latestInformation.Config
}

func New(conf Config) (Ordis, error) {
	repo, err := arango.New[id.IDer](conf.Arango)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the repository")
	}

	authModule, err := auth.New(conf.Auth, repo)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the auth module")
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

	conf.Analysis.CrashesModule = crashesModule

	analysisModule, err := analysis.New(conf.Analysis)
	if err != nil {
		return Ordis{}, errors.Wrap(err, "could not create the latestInformation module")
	}

	return Ordis{
		Auth:              authModule,
		Analysis:          analysisModule,
		Commands:          commandsModule,
		Crashes:           crashesModule,
		LatestInformation: infoModule,
	}, nil
}

func fillAuthConfig(module *auth.Module, conf *auth.AuthedConfig) {
	if conf.AuthModule == nil {
		conf.AuthModule = module
	}
}
