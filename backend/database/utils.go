// Copyright (C) 2024 (Andreas Gajdosik) <andreas@gajdosik.org>
// This file is part of project.
//
// project is non-violent software: you can use, redistribute,
// and/or modify it under the terms of the CNPLv7+ as found
// in the LICENSE file in the source code root directory or
// at <https://git.pixie.town/thufie/npl-builder>.
//
// project comes with ABSOLUTELY NO WARRANTY, to the extent
// permitted by applicable law. See the CNPL for details.

package database

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
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

func PrintEmbededAssets(assets embed.FS) {
	fmt.Println("Listing all files in the embedded assets:")
	err := fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to list embedded assets: %v", err)
	}
}
