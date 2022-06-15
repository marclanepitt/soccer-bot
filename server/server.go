package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	soccerbot "soccer-bot/m/v2"
	"soccer-bot/m/v2/commands"
	"strings"

	"github.com/nhomble/groupme.go/groupme"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", soccerbot.Port), nil))
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	message := &groupme.Message{}
	json.NewDecoder(r.Body).Decode(message)

	parts := strings.Split(message.Text, " ")
	args := strings.Join(parts[1:], " ")

	action, err := commands.GetActionFromCommand(parts[0])
	if err != nil {
		log.Println(err)
		return
	}

	err = action(args)
	if err != nil {
		log.Println(err)
		token := groupme.TokenProviderFromToken(soccerbot.Token)
		client, _ := groupme.NewClient(token)
		err = client.Bots.Send(groupme.BotMessageCommand{
			BotID:   soccerbot.BotId,
			Message: "Oops, I encountered an error...",
		})
	}
}
