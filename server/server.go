package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"soccer-bot/m/v2/commands"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	part := "/mvp"

	spew.Dump(r)
	action, err := commands.GetActionFromCommand(part)
	if err != nil {
		return
	}

	action(part)
}
