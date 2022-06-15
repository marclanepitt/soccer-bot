package commands

import (
	"context"
	"fmt"
	soccerbot "soccer-bot/m/v2"

	"github.com/densestvoid/groupme"
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
	client := groupme.NewClient(soccerbot.BotId)
	defer client.Close()

	client.CreateMessage(context.TODO(), groupme.ID(soccerbot.GroupId), &groupme.Message{
		Text: fmt.Sprintf("Congrats! Our mvp is %s", text),
	})
}
