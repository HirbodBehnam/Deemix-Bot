package types

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Downloader is an interface to download a song from a provider
type Downloader interface {
	// IsValidUrl checks if given link can be downloaded using this downloader
	IsValidUrl(u string) bool
	// Download will perform a download on the given url
	// Returns the downloaded directory path and an error
	// Download does not call IsValidUrl before downloading the link
	Download(u string) (TempDir, error)
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
