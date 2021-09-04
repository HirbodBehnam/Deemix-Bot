package main

import (
	"Deemix-Bot/bot"
	"Deemix-Bot/config"
	"os"
)

func main() {
	// Load config
	if len(os.Args) > 1 {
		config.LoadConfig(os.Args[1])
	} else {
		config.LoadConfig("config.json")
	}
	// Start bot
	bot.StartBot()
}
