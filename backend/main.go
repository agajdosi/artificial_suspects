package main

import (
	"encoding/json"
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
	mux.HandleFunc("/new_game", enableCORS(NewGameHandler))
	mux.HandleFunc("/get_game", enableCORS(GetGameHandler))
	mux.HandleFunc("/next_investigation", enableCORS(NextInvestigationHandler))
	mux.HandleFunc("/next_round", enableCORS(NextRoundHandler))
	mux.HandleFunc("/get_scores", enableCORS(GetScoresHandler))
	mux.HandleFunc("/wait_for_answer", enableCORS(WaitForAnswerHandler))
	mux.HandleFunc("/eliminate_suspect", enableCORS(EliminateSuspectHandler))
	mux.HandleFunc("/save_score", enableCORS(SaveScoreHandler))
	mux.HandleFunc("/get_services", enableCORS(GetServicesHandler))
	mux.HandleFunc("/save_answer", enableCORS(saveAnswerHandler))

	log.Println("Starting server on: http://localhost:8080")
	err = http.ListenAndServe("localhost:8080", mux)
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

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
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
	roundUUID := r.URL.Query().Get("round_uuid")
	log.Printf("WaitForAnswer() roundUUID: %s", roundUUID)
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
