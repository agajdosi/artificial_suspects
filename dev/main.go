package main

// JUST FOR TESTING NOW
// GONNA BE A HELPER BINARY TO GENERATE AI DESCRIPTIONS ETC...

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"suspects/database"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "description",
				Aliases: []string{"d"},
				Usage:   "Generate description for image",
				Action:  generateDescription,
			},
			{
				Name:    "import",
				Aliases: []string{"c"},
				Usage:   "import images from ./src/input",
				Action: func(cCtx *cli.Context) error {
					renameToSha256()
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func renameToSha256() {
	inputDir := "./frontend/public/input"   // Input directory
	outputDir := "./frontend/public/output" // Output directory

	// Create the output directory if it doesn't exist
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// Walk through the input directory
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isJpegImage(info.Name()) {
			// Open the file
			file, err := os.Open(path)
			if err != nil {
				log.Printf("failed to open file: %v", err)
				return err
			}
			defer file.Close()

			// Calculate the SHA-256 hash
			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				log.Printf("failed to calculate hash: %v", err)
				return err
			}

			sha := hex.EncodeToString(hash.Sum(nil))

			// Create new file name
			newName := fmt.Sprintf("%s.jpeg", sha)
			newPath := filepath.Join(outputDir, newName)

			// Copy file to the output directory with the new name
			err = copyFile(path, newPath)
			if err != nil {
				log.Printf("failed to copy file: %v", err)
				return err
			}

			fmt.Printf("Copied and renamed: %s -> %s\n", path, newPath)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through the directory: %v", err)
	}
}

// copyFile copies the file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, sourceFile)
	return err
}

func isJpegImage(filename string) bool {
	extensions := []string{".jpeg", ".jpg"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

func generateDescription(cCtx *cli.Context) error {
	fmt.Println("Description for: ", cCtx.Args().First())
	database.EnsureDBAvailable()

	service, err := database.GetService("OpenAI")
	if err != nil {
		return err
	}

	fmt.Println("Token:", service.Token)

	return nil
}
