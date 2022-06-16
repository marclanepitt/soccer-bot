package soccerbot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Token            string
	GroupId          string
	BotId            string
	Port             string
	DatabaseHost     string
	DatabasePort     uint16
	DatabaseUsername string
	DatabasePassword string
	VoloBearerToken  string
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

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	Token = os.Getenv("TOKEN")
	GroupId = os.Getenv("GROUP_ID")
	BotId = os.Getenv("BOT_ID")
	Port = os.Getenv("PORT")
	DatabaseHost = os.Getenv("DATABASE_HOST")
	DatabasePort = uint16(port)
	DatabaseUsername = os.Getenv("DATABASE_USERNAME")
	DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	VoloBearerToken = os.Getenv("BEARER_TOKEN")
}
