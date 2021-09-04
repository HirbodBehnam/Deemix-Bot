package deezer

import (
	"Deemix-Bot/util"
	"Deemix-Bot/util/rng"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SearchEntry interface {
	Append(builder *strings.Builder, index int)
	Article() tgbotapi.InlineQueryResultArticle
}

// Track is the entry of track searches
type Track struct {
	// The title (name) of the song
	Title string
	// The link to the song
	Link string
	// The artist name
	Artist string
	// The album name
	Album string
	// Url to album pic
	AlbumPic string
	// The duration of music
	Duration time.Duration
}

func (t Track) Append(builder *strings.Builder, index int) {
	builder.WriteString(strconv.Itoa(index + 1))
	builder.WriteByte('\n')
	builder.WriteString("Title: ")
	builder.WriteString(util.EscapeMarkdown(t.Title))
	builder.WriteByte('\n')
	builder.WriteString("Album: ")
	builder.WriteString(util.EscapeMarkdown(t.Album))
	builder.WriteByte('\n')
	builder.WriteString("Artist: ")
	builder.WriteString(util.EscapeMarkdown(t.Artist))
	builder.WriteByte('\n')
	builder.WriteString("Duration: ")
	builder.WriteString(t.Duration.String())
	builder.WriteByte('\n')
	builder.WriteString("Link: `")
	builder.WriteString(t.Link)
	builder.WriteString("`\n\n")
}

func (t Track) Article() tgbotapi.InlineQueryResultArticle {
	row := tgbotapi.NewInlineQueryResultArticle(rng.QueryToken(), t.Title, t.Link)
	row.Description = "Album: " + t.Album + "\nArtist: " + t.Artist
	row.ThumbURL = t.AlbumPic
	return row
}

// trackSearchResponse is the response of search
type trackSearchResponse struct {
	Data []struct {
		Title    string `json:"title"`
		Link     string `json:"link"`
		Duration int    `json:"duration"`
		Artist   struct {
			Name string `json:"name"`
		} `json:"artist"`
		Album struct {
			Title string `json:"title"`
			Cover string `json:"cover_medium"`
		} `json:"album"`
	} `json:"data"`
}

// TempDir is a simple structure which can hold the path to a temporary directory
type TempDir struct {
	// Address of the directory
	Address string
}

// Delete deletes the temporary directory
func (d TempDir) Delete() {
	if d.Address != "" {
		_ = os.RemoveAll(d.Address)
	}
}

// GetMusics gets the downloaded music filenames from temp dir
// If there is an error, returns nil
func (d TempDir) GetMusics() []string {
	result := make([]string, 0)
	err := filepath.WalkDir(d.Address, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".mp3") {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return result
}

// AlbumResponse is the album search result
type AlbumResponse struct {
	Data []Album `json:"data"`
}

// Album is a single album in album search
type Album struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Cover  string `json:"cover_medium"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
	TracksCount int `json:"nb_tracks"`
}

func (a Album) Append(builder *strings.Builder, index int) {
	builder.WriteString(strconv.Itoa(index + 1))
	builder.WriteByte('\n')
	builder.WriteString("Title: ")
	builder.WriteString(util.EscapeMarkdown(a.Title))
	builder.WriteByte('\n')
	builder.WriteString("Tracks Count: ")
	builder.WriteString(util.EscapeMarkdown(strconv.Itoa(a.TracksCount)))
	builder.WriteByte('\n')
	builder.WriteString("Artist: ")
	builder.WriteString(util.EscapeMarkdown(a.Artist.Name))
	builder.WriteByte('\n')
	builder.WriteString("Link: `")
	builder.WriteString(a.Link)
	builder.WriteString("`\n\n")
}

func (a Album) Article() tgbotapi.InlineQueryResultArticle {
	row := tgbotapi.NewInlineQueryResultArticle(rng.QueryToken(), a.Title, a.Link)
	row.Description = "Artist: " + a.Artist.Name + "\nTracks: " + strconv.Itoa(a.TracksCount)
	row.ThumbURL = a.Cover
	return row
}
