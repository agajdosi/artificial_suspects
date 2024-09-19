package database

import (
	"embed"
	"log"
	"path/filepath"
	"strings"
)

// loadSuspectImages reads all jpeg image files (e.g., JPEG, jpeg, JPG, jpg) from the suspects directory.
func loadSuspectImages(assets embed.FS) ([]string, error) {
	directory := "frontend/dist/suspects"
	files, err := assets.ReadDir(directory)
	if err != nil {
		log.Printf("FAILED TO READ DIRECTORY %s: %v, if this is outside bindings step, there is huge problem!", directory, err)
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
