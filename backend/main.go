package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/agajdosi/artificial_suspects/database"
)

func main() {
	err := database.EnsureDBAvailable()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/new_game", NewGameHandler)
	mux.HandleFunc("/current_game", GetCurrentGameHandler)
	mux.HandleFunc("/next_investigation", NextInvestigationHandler)
	mux.HandleFunc("/next_round", NextRoundHandler)
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	game, err := database.NewGame()
	if err != nil {
		fmt.Println("NewGame() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		fmt.Println("NewGame() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetCurrentGameHandler(w http.ResponseWriter, r *http.Request) {
	game, err := database.GetCurrentGame()
	if err != nil {
		fmt.Println("GetCurrentGame() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		fmt.Println("GetCurrentGame() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func NextInvestigationHandler(w http.ResponseWriter, r *http.Request) {
	game, err := database.GetCurrentGame()
	if err != nil {
		fmt.Println("NextInvestigation() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	game.Investigation, err = database.NewInvestigation(game.UUID)
	if err != nil {
		fmt.Println("NextInvestigation() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		fmt.Println("NextInvestigation() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func NextRoundHandler(w http.ResponseWriter, r *http.Request) {
	game, err := database.GetCurrentGame()
	if err != nil {
		fmt.Println("NextRound() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	round, err := database.NewRound(game.Investigation.UUID)
	if err != nil {
		fmt.Println("NextRound() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go database.GetAnswerFromAI(round, game.Investigation.CriminalUUID)

	game.Investigation.Rounds = append(game.Investigation.Rounds, round) // prepend

	fmt.Printf("New Round %d: %s\n", game.Level, game.Investigation.Rounds[len(game.Investigation.Rounds)-1].Question.English)

	resp, err := json.Marshal(game)
	if err != nil {
		fmt.Println("NextRound() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
