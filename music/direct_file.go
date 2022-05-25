package music

import (
	"Deemix-Bot/types"
	"Deemix-Bot/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// The client to do the requests with it
var directFileDownloadClient = http.Client{Timeout: time.Minute}

// DirectFile downloader is a type of downloader which can download files directly from internet
// Downloaded files will be treated as mp3 file
type DirectFile struct {
	// The domain of the host which we should download from
	UrlPrefix string
}

// IsValidUrl will simply check if the given url has the DirectFile.UrlPrefix prefix
func (d DirectFile) IsValidUrl(u string) bool {
	return strings.HasPrefix(u, d.UrlPrefix)
}

// Download on DirectFile will simply download the file from internet to a temp dir and returns it
func (d DirectFile) Download(u string) (types.TempDir, error) {
	// Create a temp dir
	dirName, err := ioutil.TempDir("", "direct-file*")
	if err != nil {
		return types.TempDir{}, err
	}
	result := types.TempDir{Address: dirName}
	defer func() {
		if err != nil {
			result.Delete()
		}
	}()
	// Download the file
	resp, err := directFileDownloadClient.Get(u)
	if err != nil {
		return types.TempDir{}, err
	}
	defer resp.Body.Close()
	// Create a file
	filename := path.Join(dirName, util.GetFileNameFromResponseOrUrl(resp, u))
	file, err := os.Create(filename)
	if err != nil {
		return types.TempDir{}, err
	}
	defer file.Close()
	// Copy to disk
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return types.TempDir{}, err
	}
	// Done
	return result, nil
}
