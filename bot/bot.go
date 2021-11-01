package bot

import (
	"Deemix-Bot/config"
	"Deemix-Bot/music"
	"Deemix-Bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// StartBot starts the telegram bot with config.Config.BotToken as argument
// It crashes the program if it couldn't start the bot
func StartBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.Config.BotToken)
	if err != nil {
		log.Fatal("Cannot initialize the bot: ", err.Error())
	}
	log.Println("Deemix Bot v" + config.Version)
	log.Println("Bot authorized on account", bot.Self.UserName)
	// Get updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	// Process each update
	for update := range updates {
		// Inline searches for bot
		if update.InlineQuery != nil {
			if config.CheckAuthorizedUser(update.InlineQuery.From.ID) && update.InlineQuery.Query != "" {
				go processInlineSearch(update.InlineQuery.ID, update.InlineQuery.Query)
			}
			continue
		}
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
			case "help":
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.HelpMessage))
			case "search":
				// I can cache this message but why? Is anyone going to spam this?
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.SearchHelpMessage)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
					NewInlineKeyboardButtonWithQueryCurrentChat("Search Track", ""),
					NewInlineKeyboardButtonWithQueryCurrentChat("Search Album", albumPrefix),
				))
				_, _ = bot.Send(msg)
			case "album":
				args := update.Message.CommandArguments()
				if args == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.AlbumMessageHelp)
					msg.ParseMode = markdownStyle
					_, _ = bot.Send(msg)
					continue
				}
				go processAlbumSearch(args, update.Message.Chat.ID)
			}
			continue
		}
		// Process the update
		go processUpdate(update.Message.Text, update.Message.Chat.ID)
	}
}

// processUpdate processes the text message sent to bot
func processUpdate(text string, chatID int64) {
	if util.IsUrl(text) {
		if strings.HasPrefix(text, "https://open.spotify.com") && config.HasZSpotify {
			processMusic(text, chatID, music.ZSpotify{})
		} else {
			processMusic(text, chatID, music.Deemix{})
		}
	} else {
		processTrackSearch(text, chatID)
	}
}
