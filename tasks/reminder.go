package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	soccerbot "soccer-bot/m/v2"
	"soccer-bot/m/v2/commands"
	"time"

	"github.com/nhomble/groupme.go/groupme"
)

func init() {
	task := &Task{
		Name:   "Game Reminder",
		Action: reminderAction,
		Cron:   "0 12 * * SUN",
	}
	registerTask(task)
}

func reminderAction() {
	log.Println("Running reminder task...")
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		log.Println(err)
		return
	}

	requestBodyStruct := commands.ScheduleRequest{
		OperationName: "teamSchedule",
		Variables: commands.ScheduleVariables{
			Input: commands.ScheduleInputType{
				TeamId: "629e48b59af2e04836afd7ab",
			},
		},
		Query: "query teamSchedule($input: TeamScheduleInput!) { teamSchedule(input: $input) { games { _id dateStr startTimeStr rsvp { playerResponse playerTeam { YES NO MAYBE __typename } opponentTeam { YES NO MAYBE __typename } all { userId response gender __typename } __typename } opponent { _id name color { hex __typename } __typename } location { name __typename } field_name outcome teamScore opponentScore __typename } nextGame { _id startTimeStr dateStr opponent { _id name __typename } __typename } wins losses ties __typename } }",
	}

	body, _ := json.Marshal(requestBodyStruct)
	url := "https://www.volosports.com/graphql"
	var bearer = "Bearer " + soccerbot.VoloBearerToken

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	responseObj := &commands.ScheduleResponse{}
	err = json.Unmarshal(body, responseObj)
	if err != nil {
		log.Println(err)
		return
	}

	nextGame := responseObj.Data.TeamSchedule.NextGame
	var message string
	if nextGame.ID != "" {
		message = fmt.Sprintf("Next game will be against %s at %s on %s.\n\nLike this message if you can make it!", nextGame.Opponent.Name, cleanTime(nextGame.StartTimeStr), cleanDate(nextGame.DateStr))
	} else {
		message = "unable to find next game"
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func cleanTime(input string) string {
	t, err := time.Parse("15:04", input)
	if err != nil {
		return input
	}
	return t.Format(time.Kitchen)
}

func cleanDate(date string) string {
	d, err := time.Parse("06/01/02", date)
	if err != nil {
		return date
	}
	return d.Format("01/02/06")
}
