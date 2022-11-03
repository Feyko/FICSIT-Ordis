package crashes

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"github.com/pkg/errors"
)

type Config struct {
	auth.AuthedConfig
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[domain.Crash](repository, "Crashes")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}
	baseConf := base.NewDefaultConfig(conf.AuthModule)
	if conf.NoAuth {
		baseConf = base.NewDefaultConfigNoPerm(nil)
	}
	return &Module{
		*base.NewSearchable[domain.Crash](baseConf, collection),
	}, nil
}

type Module struct {
	base.Searchable[domain.Crash]
}