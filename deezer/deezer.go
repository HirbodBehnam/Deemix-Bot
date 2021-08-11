package deezer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
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
func SearchTrack(keyword string) ([]SearchEntry, error) {
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
	var respRaw searchResponse
	err = json.NewDecoder(resp.Body).Decode(&respRaw)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// Convert the raw response to SearchResult array
	result := make([]SearchEntry, 0, maxSearchEntries)
	for i, entry := range respRaw.Data {
		if i >= maxSearchEntries { // limit entries of result
			break
		}
		result = append(result, Track{
			Title:    entry.Title,
			Link:     entry.Link,
			Artist:   entry.Artist.Name,
			Album:    entry.Album.Title,
			Duration: time.Second * time.Duration(entry.Duration),
		})
	}
	return result, nil
}

// SearchAlbum searches the deezer for an album by keyword
func SearchAlbum(keyword string) ([]SearchEntry, error) {
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
	resultInterface := make([]SearchEntry, 0, maxSearchEntries)
	for i := range result.Data {
		if i > maxSearchEntries {
			break
		}
		resultInterface = append(resultInterface, result.Data[i])
	}
	return resultInterface, nil
}

// Download tries to download a spotify/deezer track from deezer
// We return a pointer to ensure that user don't recklessly call TempDir.Delete on result
func Download(u string) (*TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "deemix*")
	if err != nil {
		return nil, err
	}
	result := &TempDir{Address: dirName}
	// Download the file
	cmd := exec.Command("deemix", "-p", dirName, u)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting deemix: %s\n", stderr.String())
		result.Delete()
		return nil, err
	}
	// Return the directory
	return result, nil
}
