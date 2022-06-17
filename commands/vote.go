package commands

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	soccerbot "soccer-bot/m/v2"
	"sort"
	"strings"

	"github.com/nhomble/groupme.go/groupme"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Vote struct {
	gorm.Model
	Name  string
	Value int
}

var parseDbUrlRegex = regexp.MustCompile("postgres:\\/\\/(.*?):(.*)@(.*):([0-9]*)\\/(.*)")

func init() {
	command := &Command{
		Name:    "Upvote",
		Action:  upvoteAction,
		Trigger: "/upvote",
	}

	dcommand := &Command{
		Name:    "Downvote",
		Action:  downvoteAction,
		Trigger: "/downvote",
	}

	lcommand := &Command{
		Name:    "Leaderboard",
		Action:  leaderboardAction,
		Trigger: "/leaderboard",
	}

	registerCommand(command)
	registerCommand(dcommand)
	registerCommand(lcommand)

	var err error
	err = initDb()
	if err != nil {
		log.Fatal(err)
	}
}

func initDb() error {
	db, err := getDb()
	if err != nil {
		return err
	}
	db.AutoMigrate(&Vote{})
	return nil
}

func getDb() (*gorm.DB, error) {
	var dsn string
	if soccerbot.DatabaseUrl != "" {
		var (
			dbHost   string
			username string
			dbName   string
			password string
			port     string
		)
		res := parseDbUrlRegex.FindAllStringSubmatch(soccerbot.DatabaseUrl, -1)
		if len(res) > 0 {
			matches := res[0]
			username = matches[1]
			password = matches[2]
			dbHost = matches[3]
			port = matches[4]
			dbName = matches[5]
		}
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", dbHost, username, dbName, port, password)
	} else {
		dsn = "host=localhost port=5432 dbname=local sslmode=disable"
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func upvoteAction(text string) error {
	err := voteAction(text, 1)
	if err != nil {
		return err
	}
	return nil
}

func downvoteAction(text string) error {
	err := voteAction(text, -1)
	if err != nil {
		return err
	}
	return nil
}

func voteAction(text string, value int) error {
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	db, err := getDb()
	if err != nil {
		return err
	}

	group, err := client.Groups.Get(soccerbot.GroupId)
	if err != nil {
		return err
	}
	members := group.Members
	membersMap := map[string]string{}
	for _, member := range members {
		membersMap[strings.ToLower(member.Nickname)] = member.Nickname
	}

	var message string
	lowerText := strings.ToLower(text)

	// if text matches member nickname
	if name := membersMap[lowerText]; name != "" {
		var vote *Vote
		err = db.First(&vote, "name = ?", name).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			vote = &Vote{
				Name:  name,
				Value: 0,
			}
			db.Create(vote)
		}
		newValue := vote.Value + value
		db.Model(&Vote{}).Where("name = ?", name).Update("value", newValue)
		message = fmt.Sprintf("%s: %d karma", name, newValue)
	} else {
		message = fmt.Sprintf("Unable to find user with nickname \"%s\"", text)
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

func leaderboardAction(text string) error {
	message := "Turtle Yards Karma Leaderboard\n---------------------\n"
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	leaderboardMap, err := getLeaderboardMap()
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(leaderboardMap))
	for k := range leaderboardMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return leaderboardMap[keys[i]] > leaderboardMap[keys[j]]
	})

	for i, k := range keys {
		message += fmt.Sprintf("%d. %s: %d\n", i+1, k, leaderboardMap[k])
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

func getLeaderboardMap() (map[string]int, error) {
	leaderBoardMap := map[string]int{}

	db, err := getDb()
	if err != nil {
		return nil, err
	}

	votes := []*Vote{}
	db.Find(&votes)

	for _, vote := range votes {
		leaderBoardMap[vote.Name] = vote.Value
	}
	return leaderBoardMap, nil
}
