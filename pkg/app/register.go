package app

import (
	"github.com/cachesdev/pokeapi-repl/pkg/cli"
)

func (app *Application) register() {
	app.cli.Register(cli.Command{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    app.cmds.CommandExit(),
	})

	app.cli.Register(cli.Command{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    app.cmds.CommandHelp(),
	})

	app.cli.Register(cli.Command{
		Name:        "map",
		Description: "Displays map locations, and goes to the next page",
		Callback:    app.cmds.Map(app.pokeApi, app.cache),
	})

	app.cli.Register(cli.Command{
		Name:        "mapb",
		Description: "Displays map locations, and goes to the previous page",
		Callback:    app.cmds.Mapb(app.pokeApi, app.cache),
	})

	app.cli.Register(cli.Command{
		Name:        "explore",
		Description: "Displays map information",
		Callback:    app.cmds.Explore(app.pokeApi, app.cache),
	})

	app.cli.Register(cli.Command{
		Name:        "catch",
		Description: "Tries to catch a pokemon",
		Callback:    app.cmds.Catch(app.pokeApi, app.cache),
	})

	app.cli.Register(cli.Command{
		Name:        "inspect",
		Description: "Inspects a pokemon you have",
		Callback:    app.cmds.Inspect(),
	})

	app.cli.Register(cli.Command{
		Name:        "pokedex",
		Description: "Shows the pokemon you've discovered",
		Callback:    app.cmds.Pokedex(),
	})
}
