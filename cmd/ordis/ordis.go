package main

import (
	"FICSIT-Ordis/internal/domain/ordis"
	"FICSIT-Ordis/internal/ports/gql"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"log"
)

func main() {
	conf := ordis.Config{
		Arango: arango.Config{
			Username:      "ordis",
			Password:      "pass",
			SuperUsername: "root",
			SuperPassword: "pass",
			DBName:        "ordis",
			Endpoints:     []string{"http://localhost:8529"},
		},
	}

	ord, err := ordis.New(conf)
	if err != nil {
		log.Fatalf("Could not create an Ordis instance: %v", err)
	}

	err = gql.Server(&ord)
	if err != nil {
		log.Fatal(err)
	}
}
