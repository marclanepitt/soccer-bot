package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"soccer-bot/m/v2/commands"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)

	port := os.Getenv("PORT")

	log.Println(port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
