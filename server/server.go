package main

import (
	"log"
	"net/http"
	"soccer-bot/m/v2/commands"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)
	log.Fatal(http.ListenAndServe(":443", nil))
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	part := "/mvp"

	action, err := commands.GetActionFromCommand(part)
	if err != nil {
		return
	}

	action(part)
}
