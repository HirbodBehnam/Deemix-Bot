package music

import (
	"Deemix-Bot/types"
	"Deemix-Bot/util"
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

type Deemix struct{}

// IsValidUrl will always return true if the given string is a valid url
// Deemix will handle the invalid urls itself
func (Deemix) IsValidUrl(u string) bool {
	return util.IsUrl(u)
}

// Download tries to download a spotify/deezer track from deezer
// We return a pointer to ensure that user don't recklessly call TempDir.Delete on result
func (Deemix) Download(u string) (types.TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "deemix*")
	if err != nil {
		return types.TempDir{}, err
	}
	result := types.TempDir{Address: dirName}
	// Download the file
	cmd := exec.Command("deemix", "-p", dirName, u)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting deemix: %s\n", stderr.String())
		result.Delete()
		return types.TempDir{}, err
	}
	// Return the directory
	return result, nil
}
