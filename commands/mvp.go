package commands

import (
	"fmt"
	soccerbot "soccer-bot/m/v2"

	"github.com/nhomble/groupme.go/groupme"
)

func init() {
	command := &Command{
		Name:    "MVP Command",
		Action:  mvpAction,
		Trigger: "/mvp",
	}
	registerCommand(command)
}

func mvpAction(text string) {
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: fmt.Sprintf("MVP is %s", text),
	})
}
