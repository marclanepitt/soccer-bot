package commands

import (
	"fmt"
	"math/rand"
	soccerbot "soccer-bot/m/v2"

	"github.com/nhomble/groupme.go/groupme"
	"github.com/peterhellberg/giphy"
)

var mvpStatements = []string{
	"Congrats on being MVP! %s gets the broadway hat tonight",
	"Congrats on being MVP! Messi? More like %s",
	"Congrats on being MVP! First round on %s!",
}

func init() {
	command := &Command{
		Name:    "Select MVP",
		Action:  mvpAction,
		Trigger: "/mvp",
	}
	registerCommand(command)
}

func mvpAction(text string) error {
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	randomStatement := mvpStatements[rand.Intn(len(mvpStatements))]
	message := fmt.Sprintf(randomStatement, text)

	g := giphy.DefaultClient

	res, err := g.Search([]string{"mvp"})
	for i := range res.Data {
		j := rand.Intn(i + 1)
		res.Data[i], res.Data[j] = res.Data[j], res.Data[i]
	}

	var imageUrl string
	if err == nil {
		imageUrl = res.Data[0].EmbedURL
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: imageUrl,
	})

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}
