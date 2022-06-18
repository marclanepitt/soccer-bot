package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	soccerbot "soccer-bot/m/v2"
	"soccer-bot/m/v2/commands"
	"soccer-bot/m/v2/tasks"
	"strings"

	"github.com/nhomble/groupme.go/groupme"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	for _, task := range tasks.RegisteredTasks {
		c.AddFunc(task.Cron, task.Action)

	}
	log.Println("Starting tasks...")
	c.Start()

	http.HandleFunc("/botRequest", routeRequest)
	log.Printf("Listening on port %s\n", soccerbot.Port)
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
