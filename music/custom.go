package music

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

type CustomSpotify struct {
	// What is the program name which we must run
	ProgramName string
}

// Download on CustomSpotify uses your custom spotify downloader to download a music from spotify
// It expects the file to be where the working directory is
func (c CustomSpotify) Download(u string) (*TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "custom-downloader*")
	if err != nil {
		return nil, err
	}
	result := &TempDir{Address: dirName}
	// Run the command
	cmd := exec.Command(c.ProgramName, u)
	cmd.Dir = dirName
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting custom spotify downloader: %s\n", stderr.String())
		result.Delete()
		return nil, err
	}
	// Return the directory
	return result, nil
}
