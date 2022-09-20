package main

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/ordis"
	"fmt"
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

	ayo, err := ordis.New(conf)
	if err != nil {
		log.Fatalf("Could not create an Ordis instance: %v", err)
	}
	//err = ayo.Commands.Create(commands.Command{Name: "oua"})
	//if err != nil {
	//	log.Fatalf("Could not create the command: %v", err)
	//}
	//got, err := ayo.Commands.List()
	//if err != nil {
	//	log.Fatalf("Could not get the command: %v", err)
	//}
	//err = ayo.Commands.Update("oua", commands.Command{Name: "yap", Response: "dang"})
	//if err != nil {
	//	log.Fatalf("Could not get the command: %v", err)
	//}
	r, err := ayo.Commands.Search("a")
	if err != nil {
		log.Fatalf("Could not search the command: %v", err)
	}
	fmt.Println(r)
}
