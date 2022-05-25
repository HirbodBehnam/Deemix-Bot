package util

import (
	"Deemix-Bot/util/rng"
	"mime"
	"net/http"
	"net/url"
	"path"
)

// GetFileNameFromResponseOrUrl gets the filename from response headers
// or if it does not exist, falls back to getting a path from url.
// If that also does not exist, it returns a random filename
func GetFileNameFromResponseOrUrl(resp *http.Response, u string) string {
	// Try get it from headers
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err == nil {
		filename := params["filename"]
		if filename != "" {
			return filename
		}
	}
	// Try to get it from url
	filename, err := url.PathUnescape(path.Base(u))
	if err == nil && filename != "" {
		return filename
	}
	// Otherwise, return random string
	return rng.RandomFilename()
}
