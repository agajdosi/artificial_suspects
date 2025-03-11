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
	mux.HandleFunc("/new_game", enableCORS(NewGameHandler))
	mux.HandleFunc("/get_game", enableCORS(GetGameHandler))
	mux.HandleFunc("/next_investigation", enableCORS(NextInvestigationHandler))
	mux.HandleFunc("/next_round", enableCORS(NextRoundHandler))
	mux.HandleFunc("/get_scores", enableCORS(GetScoresHandler))
	mux.HandleFunc("/wait_for_answer", enableCORS(WaitForAnswerHandler))
	mux.HandleFunc("/eliminate_suspect", enableCORS(EliminateSuspectHandler))
	mux.HandleFunc("/save_score", enableCORS(SaveScoreHandler))
	mux.HandleFunc("/get_services", enableCORS(GetServicesHandler))
	mux.HandleFunc("/save_service", enableCORS(SaveServiceHandler))
	mux.HandleFunc("/activate_service", enableCORS(ActivateServiceHandler))
	mux.HandleFunc("/get_default_models", enableCORS(GetDefaultModelsHandler))
	mux.HandleFunc("/get_active_service", enableCORS(GetActiveServiceHandler))
	mux.HandleFunc("/ai_service_is_ready", enableCORS(AIServiceIsReadyHandler))
	mux.HandleFunc("/list_models_ollama", enableCORS(ListModelsOllamaHandler))

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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	game, err := database.GetCurrentGame()
	if err != nil {
		fmt.Println("GetGame() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(game)
	if err != nil {
		fmt.Println("GetGame() error:", err)
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

func GetScoresHandler(w http.ResponseWriter, r *http.Request) {
	scores, err := database.GetScores()
	if err != nil {
		fmt.Println("GetScores() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(scores)
	if err != nil {
		fmt.Println("GetScores() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func WaitForAnswerHandler(w http.ResponseWriter, r *http.Request) {
	roundUUID := r.URL.Query().Get("round_uuid")
	answer := database.WaitForAnswer(roundUUID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func EliminateSuspectHandler(w http.ResponseWriter, r *http.Request) {
	suspectUUID := r.URL.Query().Get("suspect_uuid")
	roundUUID := r.URL.Query().Get("round_uuid")
	investigationUUID := r.URL.Query().Get("investigation_uuid")

	err := database.SaveElimination(suspectUUID, roundUUID, investigationUUID)
	if err != nil {
		fmt.Println("EliminateSuspect() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SaveScoreHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	gameUUID := r.URL.Query().Get("game_uuid")

	err := database.SaveScore(name, gameUUID)
	if err != nil {
		fmt.Println("SaveScore() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := database.GetServices()
	if err != nil {
		fmt.Println("GetServices() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(services)
	if err != nil {
		fmt.Println("GetServices() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func SaveServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	text_model := r.URL.Query().Get("text_model")
	visual_model := r.URL.Query().Get("visual_model")
	token := r.URL.Query().Get("token")
	url := r.URL.Query().Get("url")

	err := database.SaveService(name, text_model, visual_model, token, url)
	if err != nil {
		fmt.Println("SaveService() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ActivateServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	err := database.ActivateService(name)
	if err != nil {
		fmt.Println("ActivateService() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetDefaultModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := database.DefaultModels

	resp, err := json.Marshal(models)
	if err != nil {
		fmt.Println("GetDefaultModels() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetActiveServiceHandler(w http.ResponseWriter, r *http.Request) {
	service, err := database.GetActiveService()
	if err != nil {
		fmt.Println("GetActiveService() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(service)
	if err != nil {
		fmt.Println("GetActiveService() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// TODO: This makes sense for AI as a service calls, but not for local Ollama.
func AIServiceIsReadyHandler(w http.ResponseWriter, r *http.Request) {
	status := database.AIServiceIsReady()

	resp, err := json.Marshal(status)
	if err != nil {
		fmt.Println("AIServiceIsReady() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// TODO: This does not make much sense, typescript should just call the Ollama API directly.
func ListModelsOllamaHandler(w http.ResponseWriter, r *http.Request) {
	models, err := database.ListModelsOllama()
	if err != nil {
		fmt.Println("ListModelsOllama() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(models)
	if err != nil {
		fmt.Println("ListModelsOllama() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
