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
