package commands

import (
	"fmt"
	soccerbot "soccer-bot/m/v2"

	"github.com/nhomble/groupme.go/groupme"
)

func init() {
	command := &Command{
		Name:    "Help",
		Action:  helpAction,
		Trigger: "/help",
	}
	registerCommand(command)
}

func helpAction(text string) error {
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	var message string
	for _, command := range RegisteredCommands {
		message += fmt.Sprintf("%s - %s\n", command.Name, command.Trigger)
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}
