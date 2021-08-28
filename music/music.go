package music

import (
	"bytes"
	"github.com/dhowden/tag"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// maxHeightWidthSize is the max width and height of thumbnail in pixels
const maxHeightWidthSize = 300

// maxThumbSize is the max size of thumbnail in bytes
const maxThumbSize = 200 * 1000

// If the file size is bigger than this, we need to extract the duration
const durationNeededSize = 10 * 1000 * 1000

// a regex to extract numbers
var numberRegex = regexp.MustCompile(`[\d]+`)

// Metadata contains the metadata of a track
type Metadata struct {
	Artist string
	Album  string
	Name   string
	// Album picture in bytes
	Picture []byte
	// Duration of music in seconds
	DurationSeconds int
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
	// Get duration if needed
	var durationSeconds int
	if stat, err := file.Stat(); err != nil && stat.Size() >= durationNeededSize {
		var stdout bytes.Buffer
		cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
		cmd.Stdout = &stdout
		if cmd.Run() == nil {
			str := numberRegex.FindString(stdout.String())
			durationSeconds, _ = strconv.Atoi(str)
		}
	}
	// Return the data
	return Metadata{
		Artist:          m.Artist(),
		Album:           m.Album(),
		Name:            m.Title(),
		Picture:         pic,
		DurationSeconds: durationSeconds,
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
