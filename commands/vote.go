package commands

import (
	"fmt"
	"log"
	soccerbot "soccer-bot/m/v2"
	"sort"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nhomble/groupme.go/groupme"
)

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
	message := "Turtle Yards Karma Leaderboard\n---------------------\n"
	token := groupme.TokenProviderFromToken(soccerbot.Token)
	client, err := groupme.NewClient(token)
	if err != nil {
		return err
	}

	leaderboardMap, err := getLeaderboard()
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

func getLeaderboard() (map[string]int, error) {
	leaderboard := map[string]int{}
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Votes"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			value, _ := strconv.Atoi(string(v))
			leaderboard[string(k)] = value
		}
		return nil
	})
	return leaderboard, nil
}
