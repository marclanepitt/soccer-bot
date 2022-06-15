package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"soccer-bot/m/v2/commands"
	"strings"

	"github.com/nhomble/groupme.go/groupme"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	message := &groupme.Message{}
	json.NewDecoder(r.Body).Decode(message)

	parts := strings.Split(message.Text, " ")
	args := strings.Join(parts[1:], " ")

	action, err := commands.GetActionFromCommand(parts[0])
	if err != nil {
		return
	}

	action(args)
}
