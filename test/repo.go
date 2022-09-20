package test

import (
	"FICSIT-Ordis/internal/ports/repos/repo"
	"os"
)

func GetRepo() (repo.Repository[any], error) {
	arangoUser := os.Getenv("ORDIS_TEST_ARANGO_USER")
	arangoUser := os.Getenv("ORDIS_TEST_ARANGO_USER")
}
