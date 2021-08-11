package config

import (
	"encoding/json"
	"log"
	"os"
)

const Version = "0.0.0"

const StartMessage = "Send the deezer or spotify link to download the track.\nAlso send the name of the song to search it."

const AboutMessage = "Deemix Downloader Bot v" + Version + "\nDeemix: https://www.reddit.com/r/deemix/\nhttps://t.me/deemixbuildbot"

// Config is the list of configs of the bot and deemix
var Config struct {
	// Your telegram bot token
	BotToken string `json:"bot_token"`
	// Authorized users to use this bot. Use @myidbot to get your ID
	Users []int64 `json:"users"`
}

// Private map to check authorized users
var users map[int64]struct{}

// LoadConfig loads the config file from a location
func LoadConfig(location string) {
	bytes, err := os.ReadFile(location)
	if err != nil {
		log.Fatalf("Cannot read config file: %s\n", err)
	}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		log.Fatalf("Cannot parse config file: %s\n", err)
	}
	// Populate users map
	users = make(map[int64]struct{}, len(Config.Users))
	for _, user := range Config.Users {
		users[user] = struct{}{}
	}
	Config.Users = nil
}

// CheckAuthorizedUser checks if userId is available in users map and is allowed to use the bot
func CheckAuthorizedUser(userId int64) bool {
	_, exists := users[userId]
	return exists
}
