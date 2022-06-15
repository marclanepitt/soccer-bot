package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"soccer-bot/m/v2/commands"
)

func main() {
	http.HandleFunc("/botRequest", routeRequest)

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	part := "/mvp"

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	action, err := commands.GetActionFromCommand(part)
	if err != nil {
		return
	}

	action(part)
}
