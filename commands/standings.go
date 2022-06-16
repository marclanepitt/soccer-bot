package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	soccerbot "soccer-bot/m/v2"

	"github.com/nhomble/groupme.go/groupme"
)

type StandingsRequest struct {
	OperationName string             `json:"operationName"`
	Variables     StandingsVariables `json:"variables"`
	Query         string             `json:"query"`
}

type StandingsVariables struct {
	Input InputType `json:"input"`
}

type InputType struct {
	LeagueId string `json:"leagueId"`
}

type StandingsResponse struct {
	Data struct {
		LeagueStandings struct {
			Standings []struct {
				Rank               int `json:"rank"`
				Win                int `json:"WIN"`
				Lose               int `json:"LOSE"`
				Tie                int `json:"TIE"`
				Forfeit            int `json:"FORFEIT"`
				PointsFor          int `json:"pointsFor"`
				PointsAgainst      int `json:"pointsAgainst"`
				PointsDifferential int `json:"pointsDifferential"`
				Team               struct {
					ID    string `json:"_id"`
					Name  string `json:"name"`
					Color struct {
						Hex      string `json:"hex"`
						Typename string `json:"__typename"`
					} `json:"color"`
					Typename string `json:"__typename"`
				} `json:"team"`
				Typename string `json:"__typename"`
			} `json:"standings"`
			Typename string `json:"__typename"`
		} `json:"leagueStandings"`
	} `json:"data"`
}

func init() {
	command := &Command{
		Name:    "Standings",
		Action:  standingsAction,
		Trigger: "/standings",
	}
	registerCommand(command)
}

func standingsAction(text string) error {
	message := "League Standings\n---------------------\n Place - Name - Wins - Losses - Ties - Forfeits\n"

	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	requestBodyStruct := StandingsRequest{
		OperationName: "leagueStandings",
		Variables: StandingsVariables{
			InputType{
				LeagueId: "625dd0ab15c1288438cd4fa4",
			},
		},
		Query: "query leagueStandings($input: LeagueStandingsInput!) { leagueStandings(input: $input) {    standings {      rank      WIN      LOSE      TIE      FORFEIT      pointsFor      pointsAgainst      pointsDifferential      team {        _id        name        color {          hex          __typename        }        __typename      }      __typename    }    __typename  }}",
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

	responseObj := &StandingsResponse{}
	err = json.Unmarshal(body, responseObj)
	if err != nil {
		return err
	}

	for _, team := range responseObj.Data.LeagueStandings.Standings {
		message += fmt.Sprintf("%d - %s - %d - %d - %d - %d \n", team.Rank, team.Team.Name, team.Win, team.Lose, team.Tie, team.Forfeit)
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
