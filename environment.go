package soccerbot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Token   string
	GroupId string
	BotId   string
	Port    string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	if strings.ToLower(os.Getenv("RUN_MODE")) == "dev" {
		err := godotenv.Load(fmt.Sprintf("%s/.env", basepath))
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	Token = os.Getenv("TOKEN")
	GroupId = os.Getenv("GROUP_ID")
	BotId = os.Getenv("BOT_ID")
	Port = os.Getenv("PORT")
}
