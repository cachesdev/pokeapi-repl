package cli

import (
	"errors"
	"fmt"
)

type Command struct {
	Name        string
	Description string
	Callback    func(cli *Context) error
}

type CLI struct {
	commands map[string]Command
}

func NewCli() *CLI {
	return &CLI{
		commands: make(map[string]Command),
	}
}

func (cli *CLI) Register(command Command) {
	cli.commands[command.Name] = command
}

func (cli *CLI) Execute(c *Context) error {
	command, ok := cli.commands[c.Call]
	if !ok {
		return errors.New("[cli.Execute] Invalid command")
	}

	c.Commands = cli.commands

	err := command.Callback(c)
	if err != nil {
		return fmt.Errorf("[cli.Execute] Error executing command: %w", err)
	}

	return nil
}
