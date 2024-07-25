package config

import (
	"Deemix-Bot/music"
	"Deemix-Bot/types"
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

const Version = "1.5.1"

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
	CustomDownloaders []struct {
		Name      string `json:"name"`
		UrlPrefix string `json:"url_prefix"`
	} `json:"custom_downloaders"`
	// A list of websites which the bot can download musics from. (direct link)
	DirectDownloadHosts []string `json:"direct_download_hosts"`
	// Your telegram bot token
	BotToken string `json:"bot_token"`
	// Authorized users to use this bot. Use @myidbot to get your ID
	Users []int64 `json:"users"`
}

// Private map to check authorized users
var users map[int64]struct{}

// Downloaders is a list of all available downloaders for application
var Downloaders = make([]types.Downloader, 0)

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
	// Check is user has other downloaders
	for _, downloader := range Config.CustomDownloaders {
		_, err = exec.LookPath(downloader.Name)
		if err == nil {
			log.Println("Detected custom downloader:", downloader.Name)
			Downloaders = append(Downloaders, music.CustomDownloader{
				ProgramName: downloader.Name,
				UrlPrefix:   downloader.UrlPrefix,
			})
		}
	}
	if len(Config.ZSpotifyCredentials) != 0 {
		_, err = exec.LookPath("zspotify")
		if err == nil {
			log.Println("Detected zspotify!")
			Downloaders = append(Downloaders, music.ZSpotify{ZSpotifyCredentials: Config.ZSpotifyCredentials})
		}
	}
	for _, downloadHost := range Config.DirectDownloadHosts {
		Downloaders = append(Downloaders, music.DirectFile{UrlPrefix: downloadHost})
	}
	// Always add deemix at last
	Downloaders = append(Downloaders, music.Deemix{})
}

// CheckAuthorizedUser checks if userId is available in users map and is allowed to use the bot
func CheckAuthorizedUser(userId int64) bool {
	if len(users) == 0 { // public bot
		return true
	}
	_, exists := users[userId]
	return exists
}
