package crashes

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/base"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"github.com/pkg/errors"
	"regexp"
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
		Searchable: *base.NewSearchable[domain.Crash](baseConf, collection),
	}, nil
}

type Module struct {
	base.Searchable[domain.Crash]
	cache []domain.Crash
}

func (m *Module) Analyse(ctx context.Context, s string) ([]domain.CrashMatch, error) {
	crashes, err := m.getCache(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting cache")
	}

	var matches []domain.CrashMatch

	for _, crash := range crashes {
		for _, regex := range crash.Regexes {
			newMatches, err := m.executeRegex(&crash, regex, &s)
			if err != nil {
				return nil, err
			}

			matches = append(matches, newMatches...)
		}
	}

	return matches, nil
}

func (m *Module) getCache(ctx context.Context) ([]domain.Crash, error) {
	return m.List(ctx)
}

func (m *Module) executeRegex(crash *domain.Crash, regex string, s *string) ([]domain.CrashMatch, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, errors.Wrap(err, "invalid regexp")
	}

	var matches []domain.CrashMatch

	pairs := re.FindAllStringIndex(*s, -1)
	for _, pair := range pairs {
		start := pair[0]
		end := pair[1]
		matches = append(matches, domain.CrashMatch{
			MatchedText: (*s)[start:end],
			Crash:       crash,
			CharSpan: &domain.Span{
				Start: start,
				End:   end,
			},
		})
	}
	return matches, nil
}
