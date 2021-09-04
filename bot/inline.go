package bot

import (
	"Deemix-Bot/deezer"
	"Deemix-Bot/util/rng"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// If the inline query begins with this, we preform an album search instead of track search
const albumPrefix = "album: "

// processInlineSearch processes an inline search on deezer for user
func processInlineSearch(ID, text string) {
	if strings.HasPrefix(text, albumPrefix) {
		album := text[len(albumPrefix):]
		if album == "" {
			return
		}
		processInlineAlbumSearch(ID, album)
	} else {
		processInlineTrackSearch(ID, text)
	}
}

// processInlineTrackSearch searches the deezer for a track and answers the callback
func processInlineTrackSearch(ID, track string) {
	search, err := deezer.SearchTrack(track)
	answerInlineSearch(ID, search, err)
}

// processInlineAlbumSearch searches the deezer for an album and answers the callback
func processInlineAlbumSearch(ID, album string) {
	search, err := deezer.SearchAlbum(album)
	answerInlineSearch(ID, search, err)
}

// answerInlineSearch answers a callback query
func answerInlineSearch(ID string, data []deezer.SearchEntry, err error) {
	// At first check the error
	if err != nil {
		answerInlineQuery(ID, []interface{}{
			tgbotapi.NewInlineQueryResultArticle(rng.QueryToken(), "error", "cannot search..."),
		})
		log.Printf("cannot search keyword: %s\n", err)
		return
	}
	// Check empty results
	if len(data) == 0 {
		answerInlineQuery(ID, []interface{}{
			tgbotapi.NewInlineQueryResultArticle(rng.QueryToken(), "404", "nothing found"),
		})
		return
	}
	// Get the rows of result
	rows := make([]interface{}, len(data))
	for i, searched := range data {
		rows[i] = searched.Article()
	}
	// Send the request
	answerInlineQuery(ID, rows)
}

// answerInlineQuery is a generic function to answer a callback query
// The result must be an array of any InlineQueryResult
func answerInlineQuery(ID string, results []interface{}) {
	// This method always returns error because it can't parse it into tgbotapi.Message
	_, _ = bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: ID,
		Results:       results,
	})
}
