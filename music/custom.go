package music

import (
	"Deemix-Bot/types"
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

// CustomDownloader will invoke your custom downloader with the url given by user as it's first argument
type CustomDownloader struct {
	// What is the program name which we must run
	ProgramName string
	// The url prefix of the valid links
	UrlPrefix string
}

// IsValidUrl will simply check if the given url has the CustomDownloader.UrlPrefix prefix
func (c CustomDownloader) IsValidUrl(u string) bool {
	return strings.HasPrefix(u, c.UrlPrefix)
}

// Download on CustomSpotify uses your custom downloader to download a music from any music source!
// It expects the file to be where the PATH
func (c CustomDownloader) Download(u string) (types.TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "custom-downloader*")
	if err != nil {
		return types.TempDir{}, err
	}
	result := types.TempDir{Address: dirName}
	// Run the command
	cmd := exec.Command(c.ProgramName, u)
	cmd.Dir = dirName
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting custom spotify downloader: %s\n", stderr.String())
		result.Delete()
		return types.TempDir{}, err
	}
	// Return the directory
	return result, nil
}
