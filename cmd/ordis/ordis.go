package main

import (
	"FICSIT-Ordis/internal/core"
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/modules/commands"
	"fmt"
)

func main() {
	conf := config.OrdisConfig{
		Commands: config.CommandsConfig{
			Persistent: false,
		},
	}

	ayo := core.New(conf)
	ayo.Commands.Create(commands.Command{Name: "oua"})
	fmt.Println(ayo.Commands.List())
}
