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
