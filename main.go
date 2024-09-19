package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"suspects/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	fmt.Println("App initializing")
	database.PrintEmbededAssets(assets)
	err := database.EnsureConfigDirAvailable()
	if err != nil {
		log.Fatal(err)
	}
	err = database.EnsureDBAvailable()
	if err != nil {
		log.Fatal(err)
	}
	err = database.InitDB(assets)
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
		// Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux:   &linux.Options{},
		Windows: &windows.Options{},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Suspects",
				Message: "Conceptual game by Andreas Gajdosik.",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
