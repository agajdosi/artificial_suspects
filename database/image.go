package database

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
)

// loadSuspectImages reads all jpeg image files (e.g., JPEG, jpeg, JPG, jpg) from the suspects directory.
func loadSuspectImages(assets embed.FS) ([]string, error) {
	directory := "frontend/dist/suspects"
	files, err := assets.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() && IsImage(file.Name()) {
			imageFiles = append(imageFiles, filepath.Join(directory, file.Name()))
		}
	}

	return imageFiles, nil
}

// isImage checks if the file has an image extension.
func IsImage(filename string) bool {
	extensions := []string{".jpeg", ".jpg"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}
