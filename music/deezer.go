package music

import (
	"encoding/json"
	"net/http"
	"time"
)

// httpClient is the client to do the requests with it
var httpClient = &http.Client{Timeout: 5 * time.Second}

// trackSearchEndpoint is where we should send our search requests for tracks
const trackSearchEndpoint = "https://api.deezer.com/search"

// albumSearchEndpoint is where we should send our search requests for albums
const albumSearchEndpoint = "https://api.deezer.com/search/album"

const maxSearchEntries = 10

// SearchTrack searches the deezer for a track by keyword
func SearchTrack(keyword string) ([]Track, error) {
	// Build the request and do it
	req, err := http.NewRequest("GET", trackSearchEndpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", keyword)
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var respRaw trackSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&respRaw)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// Convert the raw response to SearchResult array
	result := make([]Track, 0, maxSearchEntries)
	for i, entry := range respRaw.Data {
		if i >= maxSearchEntries { // limit entries of result
			break
		}
		result = append(result, Track{
			Title:    entry.Title,
			Link:     entry.Link,
			Artist:   entry.Artist.Name,
			Album:    entry.Album.Title,
			AlbumPic: entry.Album.Cover,
			Duration: time.Second * time.Duration(entry.Duration),
		})
	}
	return result, nil
}

// SearchAlbum searches the deezer for an album by keyword
func SearchAlbum(keyword string) ([]Album, error) {
	// Build the request and do it
	req, err := http.NewRequest("GET", albumSearchEndpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", keyword)
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result AlbumResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(result.Data) > maxSearchEntries {
		result.Data = result.Data[:maxSearchEntries]
	}
	return result.Data, nil
}
