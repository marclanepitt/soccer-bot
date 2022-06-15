package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	soccerbot "soccer-bot/m/v2"

	"github.com/nhomble/groupme.go/groupme"
)

func main() {
	userInput := ""
	localBotUrl := fmt.Sprintf("http://localhost:%s/botRequest", soccerbot.Port)
	for userInput != "exit" {
		fmt.Printf("Hooked to %s\n", localBotUrl)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		userInput = input.Text()

		message := groupme.Message{
			Text: userInput,
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(message)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = http.Post(localBotUrl, "application/json", &buf)
		if err != nil {
			fmt.Println(err)
		}
	}
}
