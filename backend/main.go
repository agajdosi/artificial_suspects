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
	db_path := flag.String("db-path", database.GetDBPath(), "Path to the database file")
	flag.Parse()

	err := database.EnsureDBAvailable(*db_path)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/new_game", enableCORS(NewGameHandler))
	mux.HandleFunc("/get_game", enableCORS(GetGameHandler))
	mux.HandleFunc("/next_investigation", enableCORS(NextInvestigationHandler))
	mux.HandleFunc("/next_round", enableCORS(NextRoundHandler))
	mux.HandleFunc("/get_scores", enableCORS(GetScoresHandler))
	mux.HandleFunc("/eliminate_suspect", enableCORS(EliminateSuspectHandler))
	mux.HandleFunc("/save_score", enableCORS(SaveScoreHandler))
	mux.HandleFunc("/get_services", enableCORS(GetServicesHandler))
	mux.HandleFunc("/save_answer", enableCORS(saveAnswerHandler))
	mux.HandleFunc("/status", enableCORS(statusHandler))
	//	mux.HandleFunc("/wait_for_answer", enableCORS(WaitForAnswerHandler))

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
	log.Printf("🔍 NewGameHandler() request: %v", r)
	game, err := database.NewGame()
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetGameHandler() request: %v", r)
	game, err := database.GetCurrentGame()
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

func NextInvestigationHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 NextInvestigationHandler() request: %v", r)
	game, err := database.GetCurrentGame()
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

func NextRoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 NextRoundHandler() request: %v", r)
	game, err := database.GetCurrentGame()
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

func WaitForAnswerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 WaitForAnswerHandler() request: %v", r)
	roundUUID := r.URL.Query().Get("round_uuid")
	answer := database.WaitForAnswer(roundUUID)

	resp, err := json.Marshal(answer)
	if err != nil {
		log.Printf("WaitForAnswer() error: %v", err)
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

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 GetServicesHandler() request: %v", r)
	services, err := database.GetServices()
	if err != nil {
		log.Printf("GetServices() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(services)
	if err != nil {
		log.Printf("GetServices() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func saveAnswerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 saveAnswerHandler() request: %v", r)
	answer := r.URL.Query().Get("answer")
	roundUUID := r.URL.Query().Get("round_uuid")

	err := database.SaveAnswer(answer, roundUUID)
	if err != nil {
		log.Printf("SaveAnswer() error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("SaveAnswer() answer: %s, roundUUID: %s", answer, roundUUID)
	w.WriteHeader(http.StatusOK)
}
