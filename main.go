package main

import (
	"log"

	"github.com/cachesdev/pokedexcli/pkg/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("Fatal error during execution: %s", err.Error())
	}
}
