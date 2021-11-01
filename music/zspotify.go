package music

import (
	"Deemix-Bot/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

type ZSpotify struct{}

// Download tries to download a spotify/deezer track from deezer
// We return a pointer to ensure that user don't recklessly call TempDir.Delete on result
func (ZSpotify) Download(u string) (*TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "zspotify*")
	if err != nil {
		return nil, err
	}
	result := &TempDir{Address: dirName}
	// Add the files to this path needed for zspotify
	configBytes, _ := json.Marshal(zspotifyConfig{
		RootPath:          dirName,
		RootPodcastPath:   dirName,
		DownloadFormat:    "mp3",
		AntiBanWaitTime:   10,
		SkipExistingFiles: true,
		ChunkSize:         50000,
	})
	_ = os.WriteFile(path.Join(dirName, "credentials.json"), config.Config.ZSpotifyCredentials, 0666)
	_ = os.WriteFile(path.Join(dirName, "zs_config.json"), configBytes, 0666)
	// Download the file
	cmd := exec.Command("zspotify", u)
	cmd.Dir = dirName
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error on excuting zspotify: %s\n", stderr.String())
		result.Delete()
		return nil, err
	}
	// Return the directory
	return result, nil
}

type zspotifyConfig struct {
	RootPath          string `json:"ROOT_PATH"`
	RootPodcastPath   string `json:"ROOT_PODCAST_PATH"`
	SkipExistingFiles bool   `json:"SKIP_EXISTING_FILES"`
	DownloadFormat    string `json:"DOWNLOAD_FORMAT"`
	ForcePremium      bool   `json:"FORCE_PREMIUM"`
	AntiBanWaitTime   int    `json:"ANTI_BAN_WAIT_TIME"`
	OverrideAutoWait  bool   `json:"OVERRIDE_AUTO_WAIT"`
	ChunkSize         int    `json:"CHUNK_SIZE"`
	SplitAlbumDiscs   bool   `json:"SPLIT_ALBUM_DISCS"`
	DownloadRealTime  bool   `json:"DOWNLOAD_REAL_TIME"`
}
