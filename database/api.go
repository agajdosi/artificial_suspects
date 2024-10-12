package database

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuspectsHandler handles the HTTP request and returns suspects in JSON format.
func SuspectsHandler(w http.ResponseWriter, r *http.Request) {
	suspects, err := GetAllSuspects()
	if err != nil {
		http.Error(w, "Could not retrieve suspects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suspects); err != nil {
		http.Error(w, "Could not encode suspects to JSON", http.StatusInternalServerError)
		return
	}
}

func ConflictingSuspectsHandler(w http.ResponseWriter, r *http.Request) {
	// Query to get suspects with the highest number of wrong eliminations
	query := `
	SELECT 
		suspects.uuid, 
		suspects.image, 
		COUNT(DISTINCT eliminations.roundUUID) AS conflicting_count
	FROM eliminations
	JOIN rounds ON eliminations.roundUUID = rounds.uuid
	JOIN investigations ON rounds.investigation_uuid = investigations.uuid
	JOIN suspects ON eliminations.suspectUUID = suspects.uuid
	WHERE eliminations.suspectUUID = investigations.criminal_uuid
	GROUP BY suspects.uuid
	ORDER BY conflicting_count DESC
	LIMIT 10;
	`

	rows, err := database.Query(query)
	if err != nil {
		msg := fmt.Sprintf("Error getting conflicting suspects: %v", err)
		fmt.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var uuid, image string
		var wrongEliminations int
		if err := rows.Scan(&uuid, &image, &wrongEliminations); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		results = append(results, map[string]interface{}{
			"uuid":               uuid,
			"image":              image,
			"wrong_eliminations": wrongEliminations,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func ConflictingQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Query to get questions that caused the most wrong eliminations
	query := `
	SELECT 
		questions.uuid, 
		questions.English,
		questions.Czech,
		questions.Polish,
		COUNT(*) AS conflicting_count
	FROM rounds
	JOIN questions ON rounds.question_uuid = questions.uuid
	JOIN investigations ON rounds.investigation_uuid = investigations.uuid
	WHERE EXISTS (
		SELECT 1 
		FROM eliminations 
		WHERE eliminations.RoundUUID = rounds.uuid 
		AND eliminations.SuspectUUID = investigations.criminal_uuid
	)
	GROUP BY questions.uuid
	ORDER BY conflicting_count DESC
	LIMIT 15;
	`

	rows, err := database.Query(query)
	if err != nil {
		msg := fmt.Sprintf("Error getting conflicting questions: %v", err)
		fmt.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var uuid, english, czech, polish string
		var wrongEliminations int
		if err := rows.Scan(&uuid, &english, &czech, &polish, &wrongEliminations); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		results = append(results, map[string]interface{}{
			"uuid":               uuid,
			"english":            english,
			"czech":              czech,
			"polish":             polish,
			"wrong_eliminations": wrongEliminations,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
