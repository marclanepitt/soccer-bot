package commands

import (
	"fmt"
)

type Command struct {
	Name    string
	Action  func(string)
	Trigger string
}

var RegisteredCommands []*Command

func registerCommand(c *Command) {
	RegisteredCommands = append(RegisteredCommands, c)
}

func GetActionFromCommand(text string) (func(string), error) {
	for _, c := range RegisteredCommands {
		if c.Trigger == text {
			return c.Action, nil
		}
	}

	return nil, fmt.Errorf("no command found")
}
