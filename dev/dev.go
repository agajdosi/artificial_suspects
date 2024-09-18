package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"suspects/database"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "describe",
				Aliases: []string{"d"},
				Usage:   "Describe the image of specified suspect.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "suspect-id",
						Usage:    "UUID of the Suspect",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "service",
						Usage:    "Service name for the description generation",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "model",
						Usage:    "Model name to use for description generation",
						Required: true,
					},
				},
				Action: describe,
			},
			{
				Name:  "describe-all",
				Usage: "Describe the image of specified suspect.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "service",
						Usage:    "Service name for the description generation",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "model",
						Usage:    "Model name to use for description generation",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "limit",
						Usage:    "Only describe those whose number of descriptions is below the limit",
						Required: true,
					},
				},
				Action: describeAll,
			},
			{
				Name:    "import",
				Aliases: []string{"c"},
				Usage:   "import images from ./src/input",
				Action:  renameToSha256,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func renameToSha256(cCtx *cli.Context) error {
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
		if !info.IsDir() && database.IsImage(info.Name()) {
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
	return nil
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

func describe(cCtx *cli.Context) error {
	suspectUUID := cCtx.String("suspect-id")
	serviceName := cCtx.String("service")
	modelName := cCtx.String("model")
	return database.GenerateDescription(suspectUUID, serviceName, modelName)
}

func describeAll(cCtx *cli.Context) error {
	limit := cCtx.Int("limit")
	serviceName := cCtx.String("service")
	modelName := cCtx.String("model")
	return database.GenerateDescriptionsForAll(limit, serviceName, modelName)
}
