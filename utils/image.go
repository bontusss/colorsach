package utils

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
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
