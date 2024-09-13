package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"suspects/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	fmt.Println("App initializing")
	err := database.EnsureConfigDirAvailable()
	if err != nil {
		log.Fatal(err)
	}

	err = database.EnsureDBAvailable()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Main function starting")
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Suspects",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
