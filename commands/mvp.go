package commands

import "log"

func init() {
	command := &Command{
		Name:    "MVP Command",
		Action:  mvpAction,
		Trigger: "/mvp",
	}
	registerCommand(command)
}

func mvpAction(text string) {
	log.Println("asdfsa")
}
