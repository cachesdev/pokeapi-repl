package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cachesdev/pokedexcli/pkg/cache"
	"github.com/cachesdev/pokedexcli/pkg/cli"
	"github.com/cachesdev/pokedexcli/pkg/commands"
)

func Run() error {
	c := cli.NewCli()
	cmds := commands.NewCommands()
	cache := cache.NewCache(10 * time.Second)

	app := &Application{
		cli:     c,
		pokeApi: http.DefaultClient,
		cmds:    cmds,
		cache:   cache,
	}

	err := app.Start()
	if err != nil {
		return fmt.Errorf("[run] Error while executing/starting app: %w", err)
	}

	return nil
}
