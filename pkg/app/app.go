package app

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/cachesdev/pokedexcli/pkg/cache"
	"github.com/cachesdev/pokedexcli/pkg/cli"
	"github.com/cachesdev/pokedexcli/pkg/commands"
	"github.com/cachesdev/pokedexcli/pkg/repl"
)

type Application struct {
	cli     *cli.CLI
	pokeApi *http.Client
	cmds    *commands.Commands
	cache   *cache.Cache
}

func (app *Application) Start() error {
	app.register()

	fmt.Println("Welcome to the Pokedex!")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if ok := scanner.Scan(); !ok {
			continue
		}

		text := scanner.Text()
		sanitized := repl.CleanInput(text)

		c := &cli.Context{
			Input:          text,
			SanitizedInput: sanitized,
			Call:           sanitized[0],
		}

		err := app.cli.Execute(c)
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err.Error())
			continue
		}
	}
}
