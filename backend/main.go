package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/agajdosi/artificial_suspects/database"
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	host := flag.String("host", "localhost", "Host to run the server on, for production use 0.0.0.0")
	db_path := flag.String("db-path", "./data/artsus.db", "Path to the database file")
	flag.Parse()

	err := database.EnsureDBAvailable(*db_path)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// gameplay
	mux.HandleFunc("/new_game", enableCORS(NewGameHandler))
	mux.HandleFunc("/get_game", enableCORS(GetGameHandler))
	mux.HandleFunc("/eliminate_suspect", enableCORS(EliminateSuspectHandler))
	mux.HandleFunc("/next_round", enableCORS(NextRoundHandler))
	mux.HandleFunc("/next_investigation", enableCORS(NextInvestigationHandler))
	// scores
	mux.HandleFunc("/get_scores", enableCORS(GetScoresHandler))
	mux.HandleFunc("/save_score", enableCORS(SaveScoreHandler))
	// AI
	mux.HandleFunc("/get_models", enableCORS(GetModelsHandler))
	mux.HandleFunc("/get_or_generate_answer", enableCORS(GetOrGenerateAnswerHandler))
	// utils
	mux.HandleFunc("/status", enableCORS(statusHandler))

	url := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("🚀 Starting server on: http://%s", url)
	err = http.ListenAndServe(url, mux)
	if err != nil {
		log.Fatal(err)
	}
}

// CORS middleware for local development. TODO: Remove this for production.
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 statusHandler() request: %v", r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🎮 NewGameHandler() request: %v", r)
	playerUUID := r.URL.Query().Get("player_uuid")
	model := r.URL.Query().Get("model")
	if model == "" {
		log.Printf("NewGameHandler() error: query parameter 'model' cannot be empty!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if playerUUID == "" {
		log.Println("NewGameHandler() warning: player_uuid is empty! Creating new game without player.UUID.")
	}

	game, err := database.NewGame(playerUUID, model)
	if err != nil {
		log.Printf("NewGame() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		log.Printf("NewGame() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("🎮 NewGameHandler() completed successfully.")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Get the current game for the current player identified by required query parameter player_uuid.
func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetGameHandler() request: %v", r)
	playerUUID := r.URL.Query().Get("player_uuid")
	if playerUUID == "" {
		log.Printf("GetGameHandler() error: player_uuid is empty!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := database.GetCurrentGame(playerUUID)
	if err != nil {
		log.Printf("GetGame() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		log.Printf("GetGame() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Get the next investigation for the current game for the current player identified by required query parameter player_uuid.
func NextInvestigationHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 NextInvestigationHandler() request: %v", r)
	playerUUID := r.URL.Query().Get("player_uuid")
	if playerUUID == "" {
		log.Printf("NextInvestigationHandler() error: player_uuid is empty!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := database.GetCurrentGame(playerUUID)
	if err != nil {
		log.Printf("NextInvestigation() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	game.Investigation, err = database.NewInvestigation(game.UUID)
	if err != nil {
		log.Printf("NextInvestigation() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		log.Printf("NextInvestigation() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Get the next round for the current the current player identified by required query parameter player_uuid.
func NextRoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 NextRoundHandler() request: %v", r)
	playerUUID := r.URL.Query().Get("player_uuid")
	if playerUUID == "" {
		log.Printf("NextRoundHandler() error: player_uuid is empty!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := database.GetCurrentGame(playerUUID)
	if err != nil {
		log.Printf("NextRound() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	round, err := database.NewRound(game.Investigation.UUID)
	if err != nil {
		log.Printf("NextRound() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	game.Investigation.Rounds = append(game.Investigation.Rounds, round) // prepend
	log.Printf("New Round %d: %s", game.Level, game.Investigation.Rounds[len(game.Investigation.Rounds)-1].Question.English)

	resp, err := json.Marshal(game)
	if err != nil {
		log.Printf("NextRound() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	service, err := database.GetServiceForModel(game.Model)
	if err != nil {
		log.Printf("NextRoundHandler() could not get service for model %s: %v\n", game.Model, err)
		return
	}

	descriptions, err := database.GetDescriptionsForSuspect(
		game.Investigation.CriminalUUID,
		"openai",
		game.Model,
	)
	if err != nil {
		log.Println("NextRoundHandler() could not get descriptions for suspect")
		return
	}

	go database.GenerateAnswer(round.Question.English, descriptions[0].Description, game.Model, service)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetScoresHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetScoresHandler() request: %v", r)
	scores, err := database.GetScores()
	if err != nil {
		log.Printf("GetScores() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(scores)
	if err != nil {
		log.Printf("GetScores() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func EliminateSuspectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🎯 EliminateSuspectHandler() request: %v", r)
	suspectUUID := r.URL.Query().Get("suspect_uuid")
	roundUUID := r.URL.Query().Get("round_uuid")
	investigationUUID := r.URL.Query().Get("investigation_uuid")

	err := database.SaveElimination(suspectUUID, roundUUID, investigationUUID)
	if err != nil {
		log.Printf("EliminateSuspect() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SaveScoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("💰 SaveScoreHandler() request: %v", r)
	name := r.URL.Query().Get("player_name")
	gameUUID := r.URL.Query().Get("game_uuid")
	err := database.SaveScore(name, gameUUID)
	if err != nil {
		log.Printf("SaveScore() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Saved score: player_name: %s, game_uuid: %s", name, gameUUID)
	w.WriteHeader(http.StatusOK)
}

// Get all Models available in the database.
// WARNING: API keys must not leak in here, this goes to public frontend!
func GetModelsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetModelsHandler() request: %v", r)
	allowedOnly := false
	if r.URL.Query().Get("allowed_only") == "true" {
		allowedOnly = true
	}

	models, err := database.GetModels(allowedOnly)
	if err != nil {
		log.Printf("GetModels() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(models)
	if err != nil {
		log.Printf("GetServices() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// TODO: toto muzeme vlastne oddelat
// 1. generovat answer z newGame anebo z nextRound primo v Gocku
// 2. na frontend pak jen pockat skrze WaitForAnswer
func GetOrGenerateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetOrGenerateAnswerHandler() request: %v", r)
	playerUUID := r.URL.Query().Get("player_uuid")
	if playerUUID == "" {
		log.Printf("GetOrGenerateAnswerHandler() error: query parameter 'player_uuid' cannot be empty!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := database.GetCurrentGame(playerUUID)
	question := game.Investigation.Rounds[len(game.Investigation.Rounds)-1].Question.English
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() could not get currentGame: %v\n", err)
		return
	}

	log.Printf("===> game.Model: %s\n", game.Model)

	// TODO: get the service based on the current game's Model
	serviceName := "OpenAI"
	service, err := database.GetService(serviceName)
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() could not get service %s, %v\n", serviceName, err)
		return
	}

	descriptions, err := database.GetDescriptionsForSuspect(
		game.Investigation.CriminalUUID,
		"OpenAI",
		game.Model,
	)
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() could not get descriptions for suspect: %v\n", err)
		return
	}

	answer, err := database.GenerateAnswer(question, descriptions[0].Description, game.Model, service)
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() error generating answer: %v\n", err)
		return
	}

	// TODO: move to database.GenerateAnswer()?
	err = database.SaveAnswer(answer, game.Investigation.Rounds[len(game.Investigation.Rounds)-1].UUID)
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() error saving answer: %v\n", err)
		return
	}

	log.Printf("GetOrGenerateAnswerHandler() - generated answer: %s", answer)

	resp, err := json.Marshal(
		database.Answer{ // TODO: add UUID and Timestamp once Answer has its own table
			UUID:      "",
			Text:      answer,
			Timestamp: "",
		})
	if err != nil {
		log.Printf("GetOrGenerateAnswerHandler() error marshalling answer: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
