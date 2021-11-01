package bot

import (
	"Deemix-Bot/music"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// processMusic tries to download a music
func processMusic(text string, chatID int64, downloader music.Downloader) {
	// Process report
	msg, err := bot.Send(tgbotapi.NewMessage(chatID, "Searching and downloading..."))
	if err != nil {
		return
	}
	defer func(id int) {
		_, _ = bot.Send(tgbotapi.NewDeleteMessage(chatID, id))
	}(msg.MessageID)
	// Download the music
	path, err := downloader.Download(text)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot download the music"))
		log.Printf("cannot download music: %s\n", err)
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
	_, _ = bot.Send(tgbotapi.NewEditMessageText(chatID, msg.MessageID, "Uploading music..."))
	for _, toSend := range filenames {
		sendMusic(chatID, toSend)
	}
}

// sendMusic sends a music file in chat
func sendMusic(chatID int64, path string) {
	// Create the message
	msg := tgbotapi.NewAudio(chatID, path)
	// Get the metadata if possible
	if metadata, err := music.GetMusicMetadata(path); err == nil {
		msg.Title = metadata.Name
		msg.Performer = metadata.Artist
		msg.Duration = metadata.DurationSeconds
		if metadata.Picture != nil {
			msg.Thumb = tgbotapi.FileBytes{
				Name:  "thumb.jpg",
				Bytes: metadata.Picture,
			}
		}
	}
	// Send the message
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("cannot upload music: %s\n", err)
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot upload your music"))
	}
}

// processTrackSearch searches a keyword in deezer tracks
func processTrackSearch(text string, chatID int64) {
	// Get the result from deezer
	search, err := music.SearchTrack(text)
	// Send result
	sendSearchResult(chatID, search, err)
}

// processAlbumSearch searches the deezer for an album
func processAlbumSearch(text string, chatID int64) {
	// Get the result from deezer
	search, err := music.SearchAlbum(text)
	// Send result
	sendSearchResult(chatID, search, err)
}

// sendSearchResult sends the search result to user
func sendSearchResult(chatID int64, data []music.SearchEntry, err error) {
	// At first check the error
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Cannot search the keyword :|"))
		log.Printf("cannot search keyword: %s\n", err)
		return
	}
	// Check empty search results
	if len(data) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "No match found for this keyword!"))
		return
	}
	// Create the result string
	var resultString strings.Builder
	resultString.Grow(1024 * 2) // max chars in telegram message / 2
	for i, entry := range data {
		entry.Append(&resultString, i)
	}
	// Now send the message
	msg := tgbotapi.NewMessage(chatID, resultString.String())
	msg.ParseMode = markdownStyle
	_, _ = bot.Send(msg)
}
