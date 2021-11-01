package music

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

type Deemix struct{}

// Download tries to download a spotify/deezer track from deezer
// We return a pointer to ensure that user don't recklessly call TempDir.Delete on result
func (Deemix) Download(u string) (*TempDir, error) {
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
