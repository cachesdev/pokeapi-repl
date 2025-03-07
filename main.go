package main

import (
	"log"

	"github.com/cachesdev/pokeapi-repl/pkg/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("Fatal error during execution: %s", err.Error())
	}
}
