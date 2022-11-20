package latestInformation

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"github.com/pkg/errors"
)

const latestInfoID = "latestInformation"

type latestInformation id.Wrapper[domain.LatestInformation]

func (l latestInformation) ID() string {
	return l.ID_
}

type latestInformationUpdate id.WrapperUpdate[domain.LatestInformationUpdate]

type Config struct {
	auth.AuthedConfig
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[latestInformation](repository, "LatestInformation")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}
	baseConf := base.NewDefaultConfig(conf.AuthModule)
	if conf.NoAuth {
		baseConf = base.NewDefaultConfigNoPerm(nil)
	}
	return &Module{
		*base.New[latestInformation](baseConf, collection),
	}, nil
}

type Module struct {
	base base.Module[latestInformation]
}

func (m *Module) Get(ctx context.Context) (*domain.LatestInformation, error) {
	wrapper, err := m.base.Get(ctx, latestInfoID)
	if errors.Is(err, repo.ErrElementNotFound) {
		return nil, nil
	}
	return &wrapper.Wrapped, err
}

func (m *Module) Set(ctx context.Context, text string) error {
	wrapper, err := m.base.Get(ctx, latestInfoID)
	if err != nil {

		return m.base.Create(ctx, latestInformation(id.Wrap(domain.LatestInformation{Text: text, Revision: 1}, latestInfoID)))
	}
	current := wrapper.Wrapped
	if current.Text == text {
		return nil
	}
	newRevision := current.Revision + 1
	_, _, err = m.base.Update(ctx, latestInfoID, latestInformationUpdate{Wrapped: &domain.LatestInformationUpdate{Text: &text, Revision: &newRevision}})
	return err
}

func (m *Module) Remove(ctx context.Context) error {
	return m.Set(ctx, "")
}
