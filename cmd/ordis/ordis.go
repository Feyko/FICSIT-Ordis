package main

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/ordis"
	"FICSIT-Ordis/internal/ports/gql"
	"log"
)

func main() {
	conf := config.OrdisConfig{
		Arango: config.ArangoConfig{
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
