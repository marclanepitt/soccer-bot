package commands

import (
	"fmt"
	"log"
	"math/rand"
	soccerbot "soccer-bot/m/v2"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nhomble/groupme.go/groupme"
)

func init() {
	command := &Command{
		Name:    "Upvote command",
		Action:  upvoteAction,
		Trigger: "/upvote",
	}

	dcommand := &Command{
		Name:    "Upvote command",
		Action:  downvoteAction,
		Trigger: "/downvote",
	}

	lcommand := &Command{
		Name:    "Leaderboard command",
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
	count, err := vote(strings.ToLower(text), value)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("%s: %d karma", text, count)

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
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	randomStatement := mvpStatements[rand.Intn(len(mvpStatements))]
	message := fmt.Sprintf(randomStatement, text)

	err = client.Bots.Send(groupme.BotMessageCommand{
		BotID:   soccerbot.BotId,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func initDb() error {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func vote(name string, value int) (int, error) {
	var count int

	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Votes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Votes"))
		v := b.Get([]byte(name))
		count, _ = strconv.Atoi(string(v))
		return nil
	})

	if err != nil {
		return 0, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Votes"))
		err := b.Put([]byte(name), []byte(strconv.Itoa(count+value)))
		return err
	})

	if err != nil {
		return 0, err
	}

	return count + value, nil
}
