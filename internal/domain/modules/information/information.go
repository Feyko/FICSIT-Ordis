package information

import (
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"github.com/pkg/errors"
)

type Config struct {
	NoAuth bool
	Auth   *auth.Module
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[id.Wrapper[string]](repository, "Information")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}
	baseConf := base.NewDefaultConfig(conf.Auth)
	if conf.NoAuth {
		baseConf = base.NewDefaultConfigNoPerm(nil)
	}
	return &Module{
		*base.New[id.Wrapper[string]](baseConf, collection),
	}, nil
}

type Module struct {
	base base.Module[id.Wrapper[string]]
}

func (m *Module) Get(ctx context.Context) (*string, error) {
	wrapper, err := m.base.Get(ctx, "latest-information")
	if errors.Is(err, repo.ErrElementNotFound) {
		return nil, nil
	}
	return &wrapper.Wrapped, err
}

func (m *Module) Set(ctx context.Context, info string) error {
	_, err := m.base.Get(ctx, "latest-information")
	if err != nil {
		return m.base.Create(ctx, id.Wrapper[string]{ID_: "latest-information", Wrapped: info})
	}
	_, err = m.base.Update(ctx, "latest-information", id.WrapperUpdate[string]{Wrapped: &info})
	return err
}
