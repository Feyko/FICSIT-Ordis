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
	"net/url"
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

func (mod *Module) Create(ctx context.Context, crash domain.Crash) error {
	err := mod.validateCrash(crash)
	if err != nil {
		return err
	}

	err = mod.Searchable.Create(ctx, crash)
	if err != nil {
		return err
	}
	err = mod.updateCache(ctx)
	return errors.Wrap(err, "error updating cache")
}

func (mod *Module) Update(ctx context.Context, name string, crashUpdate any) (domain.Crash, error) {
	oldCrash, newCrash, err := mod.Searchable.Update(ctx, name, crashUpdate)
	if err != nil {
		return domain.Crash{}, err
	}

	err = mod.validateCrash(newCrash)
	if err != nil {
		_, _, newErr := mod.Searchable.Update(ctx, name, domain.CrashUpdate{
			Name:        &oldCrash.Name,
			Description: oldCrash.Description,
			Regexes:     oldCrash.Regexes,
			Response:    &oldCrash.Response,
		})
		if newErr != nil {
			return domain.Crash{}, errors.Wrap(err, "error while reversing update of invalid crash")
		}

		return domain.Crash{}, err
	}

	err = mod.updateCache(ctx)
	return newCrash, errors.Wrap(err, "error updating cache")
}

func (m *Module) Analyse(ctx context.Context, s string) ([]domain.CrashMatch, error) {
	var matches []domain.CrashMatch

	for _, crash := range m.cache {
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

func (m *Module) updateCache(ctx context.Context) error {
	crashes, err := m.List(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting crash list")
	}
	m.cache = crashes
	return nil
}

func (m *Module) validateCrash(crash domain.Crash) error {
	if len(crash.Regexes) == 0 {
		return errors.New("a crash must have a regex")
	}

	if (crash.Response.Text == nil || *crash.Response.Text == "") && len(crash.Response.MediaLinks) == 0 {
		return errors.New("a crash's response must have at least a text or a media link")
	}

	for i, regex := range crash.Regexes {
		_, err := regexp.Compile(regex)
		if err != nil {
			return errors.Wrapf(err, "invalid regex %v '%v'", i, regex)
		}
	}

	for i, link := range crash.Response.MediaLinks {
		_, err := url.ParseRequestURI(link)
		if err != nil {
			return errors.Wrapf(err, "invalid media link %v '%v'", i, link)
		}
	}

	return nil
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
