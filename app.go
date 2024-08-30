package main

import (
	"context"
	"fmt"

	"golang.org/x/exp/rand"
)

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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Create a new game. Returns the suspects.
func (a *App) NewGame() Game {
	fmt.Println("Loading new game")
	game := Game{
		Suspects: []Suspect{
			{UUID: "1", ImageSource: "2.jpg"},
			{UUID: "2", ImageSource: "2.jpg"},
			{UUID: "3", ImageSource: "1.jpg"},
			{UUID: "4", ImageSource: "2.jpg"},
			{UUID: "5", ImageSource: "1.jpg"},
			{UUID: "6", ImageSource: "2.jpg"},
			{UUID: "7", ImageSource: "1.jpg"},
			{UUID: "8", ImageSource: "2.jpg"},
			{UUID: "9", ImageSource: "1.jpg"},
			{UUID: "10", ImageSource: "2.jpg"},
			{UUID: "11", ImageSource: "1.jpg"},
			{UUID: "12", ImageSource: "2.jpg"},
			{UUID: "13", ImageSource: "1.jpg"},
			{UUID: "14", ImageSource: "2.jpg"},
			{UUID: "15", ImageSource: "1.jpg"},
		},
		Level: 1,
	}
	return game
}

// Loads the last game.
func (a *App) GetGame() Game {
	game := a.NewGame()
	return game
}

// New round is requested. Returns question for the round.
func (a *App) NewRound() string {
	return ""
}

// Asks the AI whether it thinks the
func (a *App) GetAnswerFromAI() bool {
	return true
}

// User selected suspect to be freed.
func (a *App) FreeSuspect(suspectUUID string) bool {
	fmt.Printf("Freeing suspect: %s\n", suspectUUID)
	return rand.Intn(2) == 1
}

type Game struct {
	Suspects []Suspect `json:"suspects"`
	Level    int       `json:"level"`
	GameUUID string    `json:"gameUUID"`
}

type Suspect struct {
	UUID        string `json:"uuid"`
	ImageSource string `json:"imageSource"`
	Free        bool   `json:"free"`
}
