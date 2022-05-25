package music

import (
	"Deemix-Bot/types"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type ZSpotify struct {
	ZSpotifyCredentials []byte
}

// IsValidUrl simply checks if the link starts with spotify domain
func (ZSpotify) IsValidUrl(u string) bool {
	return strings.HasPrefix(u, "https://open.spotify.com")
}

// Download tries to download a spotify/deezer track from deezer
// We return a pointer to ensure that user don't recklessly call TempDir.Delete on result
func (z ZSpotify) Download(u string) (types.TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "zspotify*")
	if err != nil {
		return types.TempDir{}, err
	}
	result := types.TempDir{Address: dirName}
	// Add the files to this path needed for zspotify
	credentialPath := path.Join(dirName, "credentials.json")
	_ = os.WriteFile(credentialPath, z.ZSpotifyCredentials, 0666)
	// Download the file
	cmd := exec.Command("zspotify",
		"--download-format", "mp3",
		"--root-path", dirName,
		"--credentials-location", credentialPath,
		u)
	cmd.Dir = dirName
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting zspotify: %s\n", stderr.String())
		result.Delete()
		return types.TempDir{}, err
	}
	// Return the directory
	return result, nil
}
