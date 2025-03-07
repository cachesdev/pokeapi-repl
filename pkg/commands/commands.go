package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/cachesdev/pokeapi-repl/pkg/cache"
	"github.com/cachesdev/pokeapi-repl/pkg/cli"
)

type Commands struct {
	state *state
}

func NewCommands() *Commands {
	return &Commands{
		state: &state{
			box: make(map[string]Pokemon),
		},
	}
}

type state struct {
	currentMapPage int
	box            map[string]Pokemon
}

type CliFunc func(c *cli.Context) error

func (cmd *Commands) CommandExit() CliFunc {
	return func(c *cli.Context) error {
		fmt.Println("Closing the Pokedex... Goodbye!")
		os.Exit(0)
		return nil
	}
}

func (cmd *Commands) CommandHelp() CliFunc {
	return func(c *cli.Context) error {
		fmt.Print("Usage:\n\n")

		for _, command := range c.Commands {
			fmt.Printf("%s: %s\n", command.Name, command.Description)
		}

		return nil
	}
}

type locationAreaResp struct {
	Results []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func (cmd *Commands) Map(client *http.Client, cache *cache.Cache) CliFunc {
	return func(c *cli.Context) error {
		offset := cmd.state.currentMapPage * 20
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d", offset)

		obj, ok := cache.Get(url)
		if !ok {
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("[Map] Error getting map areas: %w", err)
			}

			dec := json.NewDecoder(resp.Body)

			var data locationAreaResp
			err = dec.Decode(&data)
			if err != nil {
				return fmt.Errorf("[Map] Error decoding json: %w", err)
			}

			for _, loc := range data.Results {
				fmt.Println(loc.Name)
			}

			cache.Set(url, data)

			cmd.state.currentMapPage++
			return nil
		}

		data := obj.(locationAreaResp)
		for _, loc := range data.Results {
			fmt.Println(loc.Name)
		}

		cmd.state.currentMapPage++
		return nil
	}
}

func (cmd *Commands) Mapb(client *http.Client, cache *cache.Cache) CliFunc {
	return func(c *cli.Context) error {
		if cmd.state.currentMapPage < 2 {
			fmt.Println("You're already at the first page!")
			return nil
		}
		cmd.state.currentMapPage -= 2

		offset := cmd.state.currentMapPage * 20
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d", offset)

		obj, ok := cache.Get(url)
		if !ok {
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("[Map] Error getting map areas: %w", err)
			}

			dec := json.NewDecoder(resp.Body)

			var data locationAreaResp
			err = dec.Decode(&data)
			if err != nil {
				return fmt.Errorf("[Map] Error decoding json: %w", err)
			}

			for _, loc := range data.Results {
				fmt.Println(loc.Name)
			}

			cache.Set(url, data)

			cmd.state.currentMapPage++
			return nil
		}

		data := obj.(locationAreaResp)
		for _, loc := range data.Results {
			fmt.Println(loc.Name)
		}
		cmd.state.currentMapPage++
		return nil
	}
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type ExploreResp struct {
	Pokemon []struct {
		Encounter Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (cmd *Commands) Explore(client *http.Client, cache *cache.Cache) CliFunc {
	return func(c *cli.Context) error {
		if len(c.SanitizedInput) < 2 {
			return errors.New("[Explore] Missing argument, expected 2 but only got 1")
		}

		fmt.Printf("Exploring %s...", c.SanitizedInput[1])

		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", c.SanitizedInput[1])

		obj, ok := cache.Get(url)
		if !ok {
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("[Explore] Error getting map information: %w", err)
			}

			dec := json.NewDecoder(resp.Body)

			var data ExploreResp
			err = dec.Decode(&data)
			if err != nil {
				return fmt.Errorf("[Explore] Error decoding json: %w", err)
			}

			fmt.Println("Found Pokemon:")
			for _, pok := range data.Pokemon {
				fmt.Printf(" - %s\n", pok.Encounter.Name)
			}

			cache.Set(url, data)

			return nil
		}

		data := obj.(ExploreResp)
		fmt.Println("Found Pokemon:")
		for _, pok := range data.Pokemon {
			fmt.Printf(" : %s", pok.Encounter.Name)
		}
		return nil
	}
}

func (cmd *Commands) Catch(client *http.Client, cache *cache.Cache) CliFunc {
	return func(c *cli.Context) error {
		if len(c.SanitizedInput) < 2 {
			return errors.New("[Explore] Missing argument, expected 2 but only got 1")
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", c.SanitizedInput[1])

		url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", c.SanitizedInput[1])

		obj, ok := cache.Get(url)
		if !ok {
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("[Catch] Error fetching pokemon data: %w", err)
			}

			dec := json.NewDecoder(resp.Body)

			var data Pokemon
			err = dec.Decode(&data)
			if err != nil {
				return fmt.Errorf("[Catch] Error decoding json: %w", err)
			}

			cache.Set(url, data)

			wasCaught := AttemptCatch(data.BaseExperience)

			if !wasCaught {
				fmt.Printf("%s escaped!\n", c.SanitizedInput[1])
				return nil
			}

			fmt.Printf("%s was caught!\n", c.SanitizedInput[1])
			cmd.state.box[data.Name] = data
			return nil
		}

		data := obj.(Pokemon)
		wasCaught := AttemptCatch(data.BaseExperience)
		if !wasCaught {
			fmt.Printf("%s escaped!\n", c.SanitizedInput[1])
			return nil
		}

		fmt.Printf("%s was caught!\n", c.SanitizedInput[1])
		cmd.state.box[data.Name] = data
		return nil
	}
}

func CalculateCatchProbability(baseXP int) float64 {
	maxXP := 680.0
	minXP := 60.0

	catchRate := 1.0 - (float64(baseXP)-minXP)/(maxXP-minXP)

	if catchRate < 0.1 {
		catchRate = 0.1
	}
	if catchRate > 0.9 {
		catchRate = 0.9
	}

	return catchRate
}

func AttemptCatch(baseXP int) bool {
	probability := CalculateCatchProbability(baseXP)
	randomValue := rand.Float64()
	return randomValue < probability
}

func (cmd *Commands) Inspect() CliFunc {
	return func(c *cli.Context) error {
		if len(c.SanitizedInput) < 2 {
			return errors.New("[Inspect] Missing argument, expected 2 but only got 1")
		}

		pok, ok := cmd.state.box[c.SanitizedInput[1]]
		if !ok {
			return fmt.Errorf("[Inspect] You don't have the pokemon %s!", c.SanitizedInput[1])
		}

		fmt.Printf("Name: %s\n", pok.Name)
		fmt.Printf("Height: %d\n", pok.Height)
		fmt.Printf("Weight: %d\n", pok.Weight)

		fmt.Println("Stats:")
		for _, s := range pok.Stats {
			fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
		}

		fmt.Println("Types:")
		for _, t := range pok.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}

		return nil
	}
}

func (cmd *Commands) Pokedex() CliFunc {
	return func(c *cli.Context) error {

		fmt.Println("Your Pokedex:")
		for _, pok := range cmd.state.box {
			fmt.Printf(" - %s\n", pok.Name)
		}

		return nil
	}
}
