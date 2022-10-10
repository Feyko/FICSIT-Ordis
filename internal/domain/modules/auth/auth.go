package auth

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"github.com/pkg/errors"
)

func New[T id.IDer](rep repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[Token](rep, "Auth")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}

	return &Module{
		Collection: collection,
	}, nil
}

type Module struct {
	Collection repo.Collection[Token]
}

type Token struct {
	String string
}

func (t Token) ID() string {
	return t.String
}
