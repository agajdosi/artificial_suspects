package database

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TimeFormat string = time.RFC3339Nano
)

// Get Timestamp of current time in format unified for the whole program - time.RFC3339Nano.
func TimestampNow() string {
	return time.Now().Format(TimeFormat)
}

func DescriptionsToString(descriptions []Description) string {
	var output string
	for _, description := range descriptions {
		output += description.Description + ""
	}
	return output
}

// Get the base46 string of the image at the specified imagePath.
func ImageToBase64(imagePath string) (string, error) {
	ext := filepath.Ext(imagePath)
	ext = strings.ToLower(ext)
	if ext != ".jpg" && ext != ".jpeg" {
		return "", fmt.Errorf("unsupported image format: only JPEG files are allowed")
	}

	imgFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer imgFile.Close()

	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)

	return base64Image, nil
}
