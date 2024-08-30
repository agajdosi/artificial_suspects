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
func (a *App) NewGame() {

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
