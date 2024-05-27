package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

// DownloadImage downloads an image from a URL and returns the image.Image.
func DownloadImage(url string) (image.Image, error) {
	// Make HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check if response status code is OK
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download image, status code: " + response.Status)
	}

	// Decode image
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ExtractColorFromStrings(title string) []string {
	colors := []string{}
	colorList := []string{"black", "red", "blue", "yellow", "purple", "white", "brown", "grey", "orange", "cream", "violet", "sky blue", "gold", "teal", "peach", "pink", "silver"}
	for _, color := range colorList {
		if strings.Contains(strings.ToLower(title), color) {
			colors = append(colors, color)
		}

	}
	return colors
}

func ApplyWatermark(img image.Image, text string) image.Image {
	// Create a new image with the same dimensions as the original
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Draw the original image onto the new image
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Over)

	// Draw the text watermark
	dc := gg.NewContextForRGBA(newImg)
	dc.SetRGBA(1, 1, 1, 0.3) // Set color to white with 50% opacity
	// dc.SetColor(color.White)
	fontPath := filepath.Join(".", "DejaVuSans-Bold.ttf")
	fontSize := 250.0

	err := dc.LoadFontFace(fontPath, fontSize)
	if err != nil {
		fmt.Println("Error loading font, using default font:", err)
		// Fallback to basicfont if custom font fails to load
		face := basicfont.Face7x13
		dc.SetFontFace(face)
	}

	x := float64(bounds.Dx() / 2)               // Center horizontally
	y := float64(bounds.Dy() / 2)               // Center vertically
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5) // Anchor at the center
	dc.Stroke()

	return newImg
}

func ImageToMultipartFile(img image.Image) (*os.File, error) {
	// Create a buffer to write the image to
	buf := new(bytes.Buffer)

	// Encode the image into the buffer as JPEG
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	// Create a temporary file to act as the multipart.File
	tmpFile, err := os.CreateTemp("", "image-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}

	// Write the buffer content to the temporary file
	_, err = tmpFile.Write(buf.Bytes())
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to write to temp file: %v", err)
	}

	// Seek to the beginning of the file to mimic a file just opened
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to seek in temp file: %v", err)
	}

	return tmpFile, nil
}
