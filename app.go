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
	"context"
	"fmt"
	"log"

	"suspects/database"

	ollama "github.com/ollama/ollama/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// MARK: APP HANDLERS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Create a new game. Returns the suspects.
func (a *App) NewGame() database.Game {
	game, err := database.NewGame()
	if err != nil {
		fmt.Println("NewGame() error:", err)
	}
	return game
}

// Loads the last game.
func (a *App) GetGame() database.Game {
	game, err := database.GetCurrentGame()
	if err != nil {
		fmt.Println("GetGame()->getCurrentGame() error: ", err)
	}
	return game
}

func (a *App) NextInvestigation() database.Game {
	game, err := database.GetCurrentGame()
	if err != nil {
		log.Printf("NextInvestigation() could not get current game: %v\n", err)
	}
	game.Investigation, err = database.NewInvestigation(game.UUID)
	if err != nil {
		log.Printf("NextInvestigation() could not get new investigation: %v\n", err)
	}
	return game
}

// Next level is requested. Updates the question and level for the game object.
func (a *App) NextRound() database.Game {
	game, err := database.GetCurrentGame()
	if err != nil {
		log.Printf("NextRound() could not get current game: %v\n", err)
	}

	round, err := database.NewRound(game.Investigation.UUID)
	if err != nil {
		log.Printf("NextRound() could not get new Round: %v\n", err)
	}
	go database.GetAnswerFromAI(round, game.Investigation.CriminalUUID)

	game.Investigation.Rounds = append(game.Investigation.Rounds, round) // prepend

	fmt.Printf("New Round %d: %s\n", game.Level, game.Investigation.Rounds[len(game.Investigation.Rounds)-1].Question.English)
	return game
}

func (a *App) GetScores() []database.FinalScore {
	scores, err := database.GetScores()
	if err != nil {
		log.Println("GetScores()", err)
	}
	return scores
}

// Wait until the Answer from AI is present in the database.
// TODO: implement the actuall retrieval from the DB.
func (a *App) WaitForAnswer(roundUUID string) string {
	return database.WaitForAnswer(roundUUID)
}

// User selected Suspect to be eliminated from the Investigation.
// This will save the Elimination to Eliminatios table. Also it
func (a *App) EliminateSuspect(suspectUUID, roundUUID, investigationUUID string) {
	fmt.Printf(">>> Eliminating suspect (%s) on Investigation (%s) in Round (%s)\n", suspectUUID, investigationUUID, roundUUID)
	err := database.SaveElimination(suspectUUID, roundUUID, investigationUUID)
	if err != nil {
		log.Printf("EliminateSuspect() error: %v\n", err)
	}
}

func (a *App) SaveScore(name, gameUUID string) {
	fmt.Println("SaveScore:", name, gameUUID)
	err := database.SaveScore(name, gameUUID)
	if err != nil {
		fmt.Println("Could not save score:", err)
	}
}

func (a *App) GetServices() []database.Service {
	services, err := database.GetServices()
	if err != nil {
		log.Println("Could not get services:", err)
	}

	fmt.Println("Got Services:", services)

	return services
}

func (a *App) SaveService(name, text_model, visual_model, token, url string) {
	err := database.SaveService(name, text_model, visual_model, token, url)
	if err != nil {
		log.Println("Could not save the service", err)
	}
}

func (a *App) ActivateService(name string) {
	err := database.ActivateService(name)
	if err != nil {
		log.Println("Could not activate the service", err)
	}
}

func (a *App) GetDefaultModels() []database.Model {
	return database.DefaultModels
}

func (a *App) GetActiveService() database.Service {
	service, err := database.GetActiveService()
	if err != nil {
		fmt.Println("Could not GetActiveService()")
	}
	return service
}

func (a *App) ToggleFullscreen() {
	if runtime.WindowIsFullscreen(a.ctx) {
		runtime.WindowUnfullscreen(a.ctx)
	} else {
		runtime.WindowFullscreen(a.ctx)
	}
}

func (a *App) QuitApplication() {
	runtime.Quit(a.ctx)
}

func (a *App) AIServiceIsReady() database.ServiceStatus {
	return database.AIServiceIsReady()
}

func (a *App) ListModelsOllama() *ollama.ListResponse {
	return database.ListModelsOllama()
}

func (a *App) GetLastErrorMessage() ErrorMessage {
	return ErrorMessage{
		Severity: "error",
		Message:  "something went very wrong",
	}
}

type ErrorMessage struct {
	Severity string   `json:"Severity"`
	Title    string   `json:"Title"`
	Message  string   `json:"Message"`
	Actions  []string `json:"Actions"`
}
