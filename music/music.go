package music

import (
	"bytes"
	"github.com/dhowden/tag"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"os"
)

// maxHeightWidthSize is the max width and height of thumbnail in pixels
const maxHeightWidthSize = 300

// maxThumbSize is the max size of thumbnail in bytes
const maxThumbSize = 200 * 1000

// Metadata contains the metadata of a track
type Metadata struct {
	Artist string
	Album  string
	Name   string
	// Album picture in bytes
	Picture []byte
}

// GetMusicMetadata gets a music's metadata from file
func GetMusicMetadata(path string) (Metadata, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return Metadata{}, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	// Read the data
	m, err := tag.ReadFrom(file)
	if err != nil {
		return Metadata{}, err
	}
	// Get pic
	var pic []byte
	if m.Picture() != nil || m.Picture().Ext == "jpeg" || m.Picture().Ext == "jpg" {
		pic = resizeThumbnail(m.Picture().Data)
	}
	return Metadata{
		Artist:  m.Artist(),
		Album:   m.Album(),
		Name:    m.Title(),
		Picture: pic,
	}, nil
}

// resizeThumbnail resizes the thumbnail to make it suitable for Telegram
// https://stackoverflow.com/a/67678654/4213397
func resizeThumbnail(imageBytes []byte) []byte {
	// Read the image
	src, err := jpeg.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil
	}
	// Check valid image (don't resize them)
	if src.Bounds().Dx() <= maxHeightWidthSize && src.Bounds().Dy() <= maxHeightWidthSize {
		if len(imageBytes) <= maxThumbSize {
			return imageBytes
		} else { // no hope to reduce the file size :(
			return nil
		}
	}
	// Resize
	dst := image.NewRGBA(image.Rect(0, 0, maxHeightWidthSize, maxHeightWidthSize))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	// Done
	var output bytes.Buffer
	_ = jpeg.Encode(&output, dst, nil)
	return output.Bytes()
}
