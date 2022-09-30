package test

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"github.com/pkg/errors"
	"os"
)

func GetRepo() (repo.Repository[id.IDer], error) {
	arangoUser := os.Getenv("ORDIS_TEST_ARANGO_USER")
	arangoPassword := os.Getenv("ORDIS_TEST_ARANGO_PASSWORD")
	arangoEndpoint := os.Getenv("ORDIS_TEST_ARANGO_ENDPOINT")
	if arangoUser == "" || arangoPassword == "" || arangoEndpoint == "" {
		return memrepo.New(), nil
	}
	arangoRepo, err := arango.New[id.IDer](config.ArangoConfig{
		Username:      arangoUser,
		Password:      arangoPassword,
		SuperUsername: arangoUser,
		SuperPassword: arangoPassword,
		DBName:        "test",
		Endpoints:     []string{arangoEndpoint},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error making an Arango repo")
	}
	return arangoRepo, nil
}
