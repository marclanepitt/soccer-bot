package soccerbot

import "os"

var (
	Token   string
	GroupId string
	BotId   string
)

func init() {
	Token = os.Getenv("TOKEN")
	GroupId = os.Getenv("GROUP_ID")
	BotId = os.Getenv("BOT_ID")
}
