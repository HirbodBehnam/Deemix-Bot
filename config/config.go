package config

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

const Version = "1.4.0"

const StartMessage = "Apparently, you were worthy enough that you can use this bot! Use /help for more info"

const HelpMessage = "Send the deezer or spotify link to download the track.\n" +
	"Also send the name of the song to search it on deezer.\n" +
	"You might want to use the /album command to search for albums"

const AboutMessage = "Deemix Downloader Bot v" + Version + "\n" +
	"Deemix: https://www.reddit.com/r/deemix/\n" +
	"https://t.me/deemixbuildbot\n" +
	"Bot By Hirbod Behnam\n" +
	"Bot Source: https://github.com/HirbodBehnam/Deemix-Bot"

const AlbumMessageHelp = "To search for an album, use the /album command with the name you want to search as it's arguments\n" +
	"For example `/album Doom Soundtrack` \\(without mono space\\) searches for \"Doom Soundtrack\" on Deezer"

const SearchHelpMessage = "You can search using inline queries; Just press one of the buttons below and start typing the keyword you are looking for; " +
	"Then, just press on one of them and it will send the link to bot. Then the bot will start downloading it."

// Config is the list of configs of the bot and deemix
var Config struct {
	// ZSpotifyCredentials is the file of `credentials.json` in raw json
	ZSpotifyCredentials json.RawMessage `json:"zspotify_credentials"`
	// The program name of custom spotify downloader
	CustomSpotifyDownloaderName string `json:"custom_spotify_downloader_name"`
	// Your telegram bot token
	BotToken string `json:"bot_token"`
	// Authorized users to use this bot. Use @myidbot to get your ID
	Users []int64 `json:"users"`
}

// Private map to check authorized users
var users map[int64]struct{}

// HasZSpotify indicates if user has spotify
var HasZSpotify = false

// HasCustomSpotifyDownloader checks if the user has a custom downloader for spotify
var HasCustomSpotifyDownloader = false

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
	// Check is user has ZSpotify or custom downloader
	if Config.CustomSpotifyDownloaderName != "" {
		_, err = exec.LookPath(Config.CustomSpotifyDownloaderName)
		if err == nil {
			HasCustomSpotifyDownloader = true
			log.Println("Detected custom Spotify downloader!")
		}
	} else if len(Config.ZSpotifyCredentials) != 0 {
		_, err = exec.LookPath("zspotify")
		if err == nil {
			HasZSpotify = true
			log.Println("Detected zspotify!")
		}
	}
}

// CheckAuthorizedUser checks if userId is available in users map and is allowed to use the bot
func CheckAuthorizedUser(userId int64) bool {
	if len(users) == 0 { // public bot
		return true
	}
	_, exists := users[userId]
	return exists
}
