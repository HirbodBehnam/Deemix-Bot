package main

import (
	"Deemix-Bot/config"
	"Deemix-Bot/deezer"
	"Deemix-Bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)

var bot *tgbotapi.BotAPI

func main() {
	// Load config
	if len(os.Args) > 1 {
		config.LoadConfig(os.Args[1])
	} else {
		config.LoadConfig("config.json")
	}
	// Start bot
	var err error
	bot, err = tgbotapi.NewBotAPI(config.Config.BotToken)
	if err != nil {
		log.Fatal("Cannot initialize the bot:", err.Error())
	}
	log.Println("Deemix Bot v" + config.Version)
	log.Println("Bot authorized on account", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	// Get updates
	for update := range updates {
		// Only text messages are allowed
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		// Also check the userID
		if !config.CheckAuthorizedUser(update.Message.From.ID) {
			continue
		}
		// Check command
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.StartMessage))
			case "about":
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.AboutMessage))
			}
			continue
		}
		// Process the update
		go ProcessUpdate(update.Message.Text, update.Message.Chat.ID)
	}
}

// ProcessUpdate processes the text message sent to bot
func ProcessUpdate(text string, chatID int64) {
	if util.IsUrl(text) {
		processMusic(text, chatID)
	} else {
		processSearch(text, chatID)
	}
}

// processSearch searches a keyword in deezer
func processSearch(text string, chatID int64) {
	// Get the result from deezer
	search, err := deezer.SearchTrack(text)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot search the keyword :|"))
		log.Printf("cannot search keyword: %s\n", err)
		return
	}
	// Check empty search results
	if len(search) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "No match found for this keyword!"))
		return
	}
	// Create the result string
	var resultString strings.Builder
	resultString.Grow(1024 * 2) // max chars in telegram message / 2
	for i, entry := range search {
		resultString.WriteString(strconv.Itoa(i + 1))
		resultString.WriteByte('\n')
		resultString.WriteString("Title: ")
		resultString.WriteString(util.EscapeMarkdown(entry.Title))
		resultString.WriteByte('\n')
		resultString.WriteString("Album: ")
		resultString.WriteString(util.EscapeMarkdown(entry.Album))
		resultString.WriteByte('\n')
		resultString.WriteString("Artist: ")
		resultString.WriteString(util.EscapeMarkdown(entry.Artist))
		resultString.WriteByte('\n')
		resultString.WriteString("Link: `")
		resultString.WriteString(entry.Link)
		resultString.WriteString("`\n")
		resultString.WriteString("Duration: ")
		resultString.WriteString(entry.Duration.String())
		resultString.WriteString("\n\n")
	}
	// Now send the message
	msg := tgbotapi.NewMessage(chatID, resultString.String())
	msg.ParseMode = "MarkdownV2"
	_, _ = bot.Send(msg)
}

// processMusic tries to download a music using deemix
func processMusic(text string, chatID int64) {
	// Download the music
	path, err := deezer.Download(text)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot download the music"))
		return
	}
	defer path.Delete()
	// Get the filename of music
	filenames := path.GetMusics()
	if filenames == nil || len(filenames) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error on getting the music from disk"))
		return
	}
	// Upload the file
	for _, toSend := range getMusicsMessages(chatID, filenames) {
		_, err = bot.Send(toSend)
		if err != nil {
			log.Printf("cannot upload music: %s\n", err)
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot upload your music"))
		}
	}
}

func getMusicsMessages(chatID int64, filenames []string) []tgbotapi.Chattable {
	// If there is only one music, use newAudio
	if len(filenames) == 1 {
		return []tgbotapi.Chattable{tgbotapi.NewAudio(chatID, filenames[0])}
	}
	// Otherwise, create album for all of them
	result := make([]tgbotapi.Chattable, 0, len(filenames)/10+1)
	var tempFilenames []interface{}
	for i, filename := range filenames {
		if i%10 == 0 {
			if i != 0 {
				result = append(result, tgbotapi.NewMediaGroup(chatID, tempFilenames))
			}
			tempFilenames = make([]interface{}, 0, 10)
		}
		tempFilenames = append(tempFilenames, tgbotapi.NewAudio(chatID, tgbotapi.NewInputMediaAudio(filename)))
	}
	// For the last one, we should use single audio track for 1 audio
	if len(tempFilenames) == 1 {
		result = append(result, tgbotapi.NewAudio(chatID, tempFilenames[0].(tgbotapi.InputMediaAudio).Media))
	} else if len(tempFilenames) > 1 {
		result = append(result, tgbotapi.NewMediaGroup(chatID, tempFilenames))
	}
	return result
}
