package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	soccerbot "soccer-bot/m/v2"
	"time"

	"github.com/nhomble/groupme.go/groupme"
)

type ScheduleRequest struct {
	OperationName string            `json:"operationName"`
	Variables     ScheduleVariables `json:"variables"`
	Query         string            `json:"query"`
}

type ScheduleVariables struct {
	Input ScheduleInputType `json:"input"`
}

type ScheduleInputType struct {
	TeamId string `json:"teamId"`
}

type ScheduleResponse struct {
	Data struct {
		TeamSchedule struct {
			Games []struct {
				ID           string `json:"_id"`
				DateStr      string `json:"dateStr"`
				StartTimeStr string `json:"startTimeStr"`
				Rsvp         struct {
					PlayerResponse string `json:"playerResponse"`
					PlayerTeam     struct {
						Yes      int    `json:"YES"`
						No       int    `json:"NO"`
						Maybe    int    `json:"MAYBE"`
						Typename string `json:"__typename"`
					} `json:"playerTeam"`
					OpponentTeam struct {
						Yes      int    `json:"YES"`
						No       int    `json:"NO"`
						Maybe    int    `json:"MAYBE"`
						Typename string `json:"__typename"`
					} `json:"opponentTeam"`
					All []struct {
						UserID   string `json:"userId"`
						Response string `json:"response"`
						Gender   string `json:"gender"`
						Typename string `json:"__typename"`
					} `json:"all"`
					Typename string `json:"__typename"`
				} `json:"rsvp"`
				Opponent struct {
					ID    string `json:"_id"`
					Name  string `json:"name"`
					Color struct {
						Hex      string `json:"hex"`
						Typename string `json:"__typename"`
					} `json:"color"`
					Typename string `json:"__typename"`
				} `json:"opponent"`
				Location struct {
					Name     string `json:"name"`
					Typename string `json:"__typename"`
				} `json:"location"`
				FieldName     string `json:"field_name"`
				Outcome       string `json:"outcome"`
				TeamScore     int    `json:"teamScore"`
				OpponentScore int    `json:"opponentScore"`
				Typename      string `json:"__typename"`
			} `json:"games"`
			NextGame struct {
				ID           string `json:"_id"`
				StartTimeStr string `json:"startTimeStr"`
				DateStr      string `json:"dateStr"`
				Opponent     struct {
					ID       string `json:"_id"`
					Name     string `json:"name"`
					Typename string `json:"__typename"`
				} `json:"opponent"`
				Typename string `json:"__typename"`
			} `json:"nextGame"`
			Wins     int    `json:"wins"`
			Losses   int    `json:"losses"`
			Ties     int    `json:"ties"`
			Typename string `json:"__typename"`
		} `json:"teamSchedule"`
	} `json:"data"`
}

func init() {
	command := &Command{
		Name:    "Schedule",
		Action:  scheduleAction,
		Trigger: "/schedule",
	}
	nextCommand := &Command{
		Name:    "Next Game",
		Action:  nextGameAction,
		Trigger: "/nextgame",
	}
	registerCommand(command)
	registerCommand(nextCommand)
}

func scheduleAction(text string) error {
	message := "League Schedule\n---------------------\n"

	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	requestBodyStruct := ScheduleRequest{
		OperationName: "teamSchedule",
		Variables: ScheduleVariables{
			ScheduleInputType{
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
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseObj := &ScheduleResponse{}
	err = json.Unmarshal(body, responseObj)
	if err != nil {
		return err
	}

	for _, game := range responseObj.Data.TeamSchedule.Games {
		message += fmt.Sprintf("%s - %s - vs. %s \n", cleanDate(game.DateStr), cleanTime(game.StartTimeStr), game.Opponent.Name)
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func nextGameAction(text string) error {
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	requestBodyStruct := ScheduleRequest{
		OperationName: "teamSchedule",
		Variables: ScheduleVariables{
			ScheduleInputType{
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
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseObj := &ScheduleResponse{}
	err = json.Unmarshal(body, responseObj)
	if err != nil {
		return err
	}

	nextGame := responseObj.Data.TeamSchedule.NextGame
	var message string
	if nextGame.ID != "" {
		message = fmt.Sprintf("Next game will be against %s at %s on %s", nextGame.Opponent.Name, cleanTime(nextGame.StartTimeStr), cleanDate(nextGame.DateStr))
	} else {
		message = "unable to find next game"
	}

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
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
