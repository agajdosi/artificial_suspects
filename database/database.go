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

package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/rand"
)

var database *sql.DB

const (
	defaultPlayerName string = "anonymous"
	appName           string = "suspects"
	numSuspect        int    = 15 // How many suspects are in one investigation - there were 12 in original board game.
)

// MARK: GENERAL DATABASE

// Get path to directory where Game will save its data.
func GetDataDirPath() string {
	return filepath.Join(xdg.ConfigHome, appName)
}

// Get path to Database located in Data Directory.
func GetDBPath() string {
	return filepath.Join(GetDataDirPath(), "default.db")
}

func EnsureConfigDirAvailable() error {
	DataDir := GetDataDirPath()
	return os.MkdirAll(DataDir, 0755)
}

// Ensure that database is ready to be used. First, check if gamesDir exists, if not create it.
// Then, check if database file exists, if not create it and initialize it.
// Returns the database connection.
func EnsureDBAvailable() error {
	gameDBPath := GetDBPath()
	fmt.Printf("Checking the database file at: %s\n", gameDBPath)
	_, err := os.Stat(gameDBPath)
	if os.IsNotExist(err) {
		fmt.Printf("Database file %s does not exist, creating it...\n", gameDBPath)
		file, err := os.Create(gameDBPath)
		if err != nil {
			return err
		}
		file.Close()
		fmt.Println("Database successfully created!")
	}

	db, err := sql.Open("sqlite3", gameDBPath)
	if err != nil {
		log.Fatal(err)
	}
	database = db
	fmt.Println("Database successfully opened!")

	return nil
}

func InitDB(assets embed.FS) error {
	var tables = []string{
		createGamesTable,
		createInvestigationsTable,
		createRoundsTable,
		createEliminationsTable,
		createServicesTable,
		createModelsTable,
		createDescriptionsTable,
	}
	for i := range tables {
		_, err := database.Exec(tables[i])
		if err != nil {
			fmt.Printf("Error initializing table: '%s', error: %v", tables[i], err)
			return err
		}
	}
	err := InitQuestionsTable()
	if err != nil {
		return err
	}

	err = InitSuspectsTable(assets)
	if err != nil {
		return err
	}
	log.Println("Database successfully initiated.")

	return nil
}

// MARK: SUSPECT

type Suspect struct {
	UUID      string `json:"UUID"`
	Image     string `json:"Image"`
	Free      bool   `json:"Free"`
	Fled      bool   `json:"Fled"`
	Timestamp string `json:"Timestamp"`
}

const createSuspectsTable = `
	CREATE TABLE IF NOT EXISTS suspects (
		uuid TEXT PRIMARY KEY,
		image TEXT,
		timestamp TEXT
	);`

func InitSuspectsTable(assets embed.FS) error {
	_, err := database.Exec(createSuspectsTable)
	if err != nil {
		return err
	}

	imagePaths, err := loadSuspectImages(assets)
	if err != nil {
		return err
	}

	for _, imagePath := range imagePaths {
		filename := filepath.Base(imagePath)
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		suspect := Suspect{
			UUID:  name,
			Image: filename,
		}
		err := SaveSuspect(suspect)
		if err != nil {
			log.Println("Cannot initialize Suspect:", err)
			return err
		}
	}

	return nil
}

func SaveSuspect(suspect Suspect) error {
	var exists bool
	if suspect.UUID == "" {
		return fmt.Errorf("suspect.UUID cannot be empty")
	}

	checkQuery := "SELECT EXISTS(SELECT 1 FROM suspects WHERE image = ?)"
	err := database.QueryRow(checkQuery, suspect.Image).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	timestamp := TimestampNow()
	query := "INSERT into suspects (uuid, image, timestamp) VALUES (?, ?, ?)"
	_, err = database.Exec(query, suspect.UUID, suspect.Image, timestamp)
	if err != nil {
		log.Printf("Could not save Suspect %s (%s): %v", suspect.Image, suspect.UUID, err)
		return err
	}

	return nil
}

// Get the basic suspect data from the Database without Suspect.Free field!
// Because Suspect.Free and Suspect.Fled needs information from table Investigation->Rounds->Eliminations.
func GetSuspect(suspectUUID string) (Suspect, error) {
	var suspect Suspect
	row := database.QueryRow("SELECT uuid, image, timestamp FROM suspects WHERE uuid = $1 LIMIT 1", suspectUUID)
	err := row.Scan(&suspect.UUID, &suspect.Image, &suspect.Timestamp)
	if err != nil {
		log.Printf("Could not load Suspect (%s): %v", suspectUUID, err)
		return suspect, err
	}

	return suspect, nil
}

// Get all Suspects and their complete data for specified Investigation.
// It needs Investigation because we need to iterate over its Rounds and Rounds' Eliminations
// to set Suspect.Free and Suspect.Fled booleans.
func getSuspectsInInvestigation(suspectUUIDs []string, investigation Investigation) ([]Suspect, error) {
	var suspects []Suspect
	eliminatedSuspectUUIDs := make(map[string]struct{})
	for i := range investigation.Rounds {
		round := investigation.Rounds[i]
		for x := range round.Eliminations {
			elimination := round.Eliminations[x]
			eliminatedSuspectUUIDs[elimination.SuspectUUID] = struct{}{}
		}
	}

	var err error
	for x := range suspectUUIDs {
		var suspect Suspect
		suspect, err = GetSuspect(suspectUUIDs[x])
		if err != nil {
			log.Printf("Error iterating over suspects: %v", err)
		}

		if _, found := eliminatedSuspectUUIDs[suspect.UUID]; found {
			if suspect.UUID == investigation.CriminalUUID {
				suspect.Fled = true
			} else {
				suspect.Free = true
			}
		}

		suspects = append(suspects, suspect)
	}

	return suspects, err
}

func GetAllSuspects() ([]Suspect, error) {
	var suspects []Suspect
	rows, err := database.Query("SELECT uuid, image, timestamp FROM suspects", numSuspect)
	if err != nil {
		log.Printf("Could not get random suspects: %v\n", err)
		return suspects, err
	}
	defer rows.Close()

	for rows.Next() {
		var suspect Suspect
		err := rows.Scan(&suspect.UUID, &suspect.Image, &suspect.Timestamp)
		if err != nil {
			log.Printf("Could not scan suspect: %v\n", err)
			return suspects, err
		}
		suspects = append(suspects, suspect)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during suspects rows iteration: %v\n", err)
		return suspects, err
	}

	return suspects, nil
}

// Get suspects and filter them by number of their descriptions by selected model.
// Only suspects with less than limit number of descriptions will be returned.
// Mostly used just for generation of descriptions in dev.go.
func GetSuspectsByDescriptions(limit int, serviceName, modelName string) ([]Suspect, error) {
	var suspects []Suspect
	allSuspects, err := GetAllSuspects()
	if err != nil {
		return suspects, err
	}
	for _, suspect := range allSuspects {
		descriptions, err := GetDescriptionsForSuspect(suspect.UUID, serviceName, modelName)
		if err != nil {
			fmt.Printf("Error getting descriptions for suspect (%s): %v", suspect.UUID, err)
			continue
		}
		if len(descriptions) >= limit {
			continue
		}
		suspects = append(suspects, suspect)
	}
	return suspects, nil
}

func randomSuspects() ([]Suspect, error) {
	var suspects []Suspect
	rows, err := database.Query("SELECT uuid, image, timestamp FROM suspects ORDER BY RANDOM() LIMIT $1", numSuspect)
	if err != nil {
		log.Printf("Could not get random suspects: %v\n", err)
		return suspects, err
	}
	defer rows.Close()

	for rows.Next() {
		var suspect Suspect
		err := rows.Scan(&suspect.UUID, &suspect.Image, &suspect.Timestamp)
		if err != nil {
			log.Printf("Could not scan suspect: %v\n", err)
			return suspects, err
		}
		suspects = append(suspects, suspect)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during suspects rows iteration: %v\n", err)
		return suspects, err
	}

	return suspects, nil
}

// MARK: GAME

// User clicks on start and plays until they make a mistake, can be several cases. This is the Game.
// TODO: add Score.
// TODO: add Name, so the player can sign their high score.
type Game struct {
	UUID          string        `json:"uuid"`
	Investigation Investigation `json:"investigation"` // TODO: actually this could be Investigations []Investigation
	Level         int           `json:"level"`         // aka number of Investigations done + 1
	Score         int           `json:"Score"`         // TODO: implement
	GameOver      bool          `json:"GameOver"`      // TODO: when true, Game is over
	Investigator  string        `json:"Investigator"`  // aka the Player's nickname
	Timestamp     string        `json:"Timestamp"`     // when game was created
}

const createGamesTable = `
	CREATE TABLE IF NOT EXISTS games (
		uuid TEXT PRIMARY KEY,
		score INT,
		investigator TEXT,
		timestamp TEXT
	);`

func NewGame() (Game, error) {
	var game Game
	game.UUID = uuid.New().String()
	game.Timestamp = TimestampNow()
	game.Score = 0
	game.Investigator = defaultPlayerName
	err := saveGame(game)
	if err != nil {
		return game, err
	}

	game.Investigation, err = NewInvestigation(game.UUID)
	if err != nil {
		return game, err
	}
	game.Level, err = getLevel(game.UUID)
	if err != nil {
		return game, err
	}

	go GetAnswerFromAI(game.Investigation.Rounds[0], game.Investigation.CriminalUUID)

	return game, err
}

func GetCurrentGame() (Game, error) {
	var game Game
	row := database.QueryRow("SELECT uuid, timestamp, score FROM games ORDER BY timestamp DESC LIMIT 1")
	err := row.Scan(&game.UUID, &game.Timestamp, &game.Score)

	// No game found - first play
	if err == sql.ErrNoRows {
		log.Println("Warning: No games in DB, creating new game")
		return NewGame()
	}
	if err != nil {
		return game, err
	}

	log.Printf("Got game: %v | %v", game.UUID, game.Timestamp)

	game.Investigation, err = getCurrentInvestigation(game.UUID)
	if err != nil {
		fmt.Println("GetGame()->getCurrentInvestigation(): ", err)
		return game, err
	}

	game.Level, err = getLevel(game.UUID)
	if err != nil {
		log.Printf("GetCurrentGame() could not get Level: %v\n", err)
		return game, err
	}

	game.GameOver = isGameOver(game)

	return game, nil
}

func saveGame(game Game) error {
	query := `INSERT INTO games (uuid, timestamp, score, investigator) VALUES (?, ?, ?, ?)`
	_, err := database.Exec(query, game.UUID, game.Timestamp, game.Score, game.Investigator)
	return err
}

func isGameOver(game Game) bool {
	for x := range game.Investigation.Rounds {
		round := game.Investigation.Rounds[x]
		for y := range round.Eliminations {
			elimination := round.Eliminations[y]
			if elimination.SuspectUUID == game.Investigation.CriminalUUID {
				return true
			}
		}
	}
	return false
}

// MARK: INVESTIGATION

// Investigation is a set of X Suspects, User needs to find a Criminal among them.
type Investigation struct {
	UUID              string    `json:"uuid"`
	GameUUID          string    `json:"game_uuid"`
	Suspects          []Suspect `json:"suspects"`
	Rounds            []Round   `json:"rounds"` // Ordered from oldest (first) to newest (last), 1st round is [0], 2nd [1] etc.
	CriminalUUID      string    `json:"CriminalUUID"`
	InvestigationOver bool      `json:"InvestigationOver"` // Last standing is the Criminal
	Timestamp         string    `json:"Timestamp"`
}

// Original has 12 suspects, for now I plan 15.
// Do not know how to make the array in more elegant way.
// This ugly shit works - for now.
const createInvestigationsTable = `
	CREATE TABLE IF NOT EXISTS investigations (
		uuid TEXT PRIMARY KEY,
		game_uuid TEXT,
		timestamp TEXT,
		criminal_uuid TEXT,
		sus1_uuid TEXT,
		sus2_uuid TEXT,
		sus3_uuid TEXT,
		sus4_uuid TEXT,
		sus5_uuid TEXT,
		sus6_uuid TEXT,
		sus7_uuid TEXT,
		sus8_uuid TEXT,
		sus9_uuid TEXT,
		sus10_uuid TEXT,
		sus11_uuid TEXT,
		sus12_uuid TEXT,
		sus13_uuid TEXT,
		sus14_uuid TEXT,
		sus15_uuid TEXT
	);`

func saveInvestigation(investigation Investigation) error {
	if len(investigation.Suspects) != 15 {
		err := fmt.Errorf("Investigation does not have 15 suspects, has %d", (len(investigation.Suspects)))
		log.Printf("Cannot save investigation: %v\n", err)
		return err
	}

	query := `INSERT OR REPLACE INTO investigations
		(uuid, game_uuid, timestamp,
		sus1_uuid,
		sus2_uuid,
		sus3_uuid,
		sus4_uuid,
		sus5_uuid,
		sus6_uuid,
		sus7_uuid,
		sus8_uuid,
		sus9_uuid,
		sus10_uuid,
		sus11_uuid,
		sus12_uuid,
		sus13_uuid,
		sus14_uuid,
		sus15_uuid,
		criminal_uuid
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := database.Exec(query, investigation.UUID, investigation.GameUUID, investigation.Timestamp,
		investigation.Suspects[0].UUID,
		investigation.Suspects[1].UUID,
		investigation.Suspects[2].UUID,
		investigation.Suspects[3].UUID,
		investigation.Suspects[4].UUID,
		investigation.Suspects[5].UUID,
		investigation.Suspects[6].UUID,
		investigation.Suspects[7].UUID,
		investigation.Suspects[8].UUID,
		investigation.Suspects[9].UUID,
		investigation.Suspects[10].UUID,
		investigation.Suspects[11].UUID,
		investigation.Suspects[12].UUID,
		investigation.Suspects[13].UUID,
		investigation.Suspects[14].UUID,
		investigation.CriminalUUID,
	)
	if err != nil {
		log.Printf("Could not save investigation: %v", err)
		return err
	}

	return nil
}

// Create a new Investigation, save it into the database and return it.
// Usage on New Game for initial first Investigation,
// or when Investigation is successfully solved and we need new one.
func NewInvestigation(gameUUID string) (Investigation, error) {
	var i Investigation
	i.UUID = uuid.New().String()
	i.GameUUID = gameUUID
	i.Timestamp = TimestampNow()

	round, err := NewRound(i.UUID)
	if err != nil {
		return i, err
	}
	i.Rounds = append(i.Rounds, round)

	suspects, err := randomSuspects()
	if err != nil {
		return i, err
	}
	i.Suspects = suspects

	cn := rand.Intn(len(suspects))
	i.CriminalUUID = i.Suspects[cn].UUID

	go GetAnswerFromAI(i.Rounds[0], i.CriminalUUID)
	log.Printf("NEW INVESTIGATION, criminal is: no. %d\n", cn+1)

	err = saveInvestigation(i)
	return i, err
}

func getCurrentInvestigation(gameUUID string) (Investigation, error) {
	var investigation = Investigation{GameUUID: gameUUID}
	var suspects_uuids = make([]string, 15)
	log.Printf("Getting investigation for game %s\n", gameUUID)
	row := database.QueryRow(`SELECT uuid, timestamp, criminal_uuid,
		sus1_uuid,
		sus2_uuid,
		sus3_uuid,
		sus4_uuid,
		sus5_uuid,
		sus6_uuid,
		sus7_uuid,
		sus8_uuid,
		sus9_uuid,
		sus10_uuid,
		sus11_uuid,
		sus12_uuid,
		sus13_uuid,
		sus14_uuid,
		sus15_uuid
		FROM investigations WHERE game_uuid = $1 ORDER BY timestamp DESC LIMIT 1`, gameUUID)
	err := row.Scan(&investigation.UUID, &investigation.Timestamp, &investigation.CriminalUUID,
		&suspects_uuids[0],
		&suspects_uuids[1],
		&suspects_uuids[2],
		&suspects_uuids[3],
		&suspects_uuids[4],
		&suspects_uuids[5],
		&suspects_uuids[6],
		&suspects_uuids[7],
		&suspects_uuids[8],
		&suspects_uuids[9],
		&suspects_uuids[10],
		&suspects_uuids[11],
		&suspects_uuids[12],
		&suspects_uuids[13],
		&suspects_uuids[14],
	)
	if err != nil {
		log.Printf("Could not get investigation: %v\n", err)
		return investigation, err
	}

	investigation.Rounds, err = getRounds(investigation.UUID)
	if err != nil {
		return investigation, err
	}

	investigation.Suspects, err = getSuspectsInInvestigation(suspects_uuids, investigation)
	if err != nil {
		return investigation, err
	}

	eliminated := 0
	for x := range investigation.Rounds {
		eliminated += len(investigation.Rounds[x].Eliminations)
	}
	if eliminated == (numSuspect - 1) {
		investigation.InvestigationOver = true
	}

	return investigation, nil
}

// MARK: ROUND

type Round struct {
	UUID              string        `json:"uuid"`
	InvestigationUUID string        `json:"InvestigationUUID"`
	Question          Question      `json:"Question"`
	AnswerUUID        string        `json:"AnswerUUID"`
	Answer            string        `json:"answer"` // TODO: Answer could be actually stored in table
	Eliminations      []Elimination `json:"Eliminations"`
	Timestamp         string        `json:"Timestamp"`
}

const createRoundsTable = `
	CREATE TABLE IF NOT EXISTS rounds (
		uuid TEXT PRIMARY KEY,
		investigation_uuid TEXT,
		question_uuid TEXT,
		answer TEXT,
		timestamp TEXT
	);`

func saveRound(r Round) error {
	query := `
		INSERT OR REPLACE INTO rounds (uuid, investigation_uuid, question_uuid, answer, timestamp)
		VALUES (?, ?, ?, ?, ?)
		`
	_, err := database.Exec(query, r.UUID, r.InvestigationUUID, r.Question.UUID, r.Answer, r.Timestamp)
	return err
}

func NewRound(investigationUUID string) (Round, error) {
	var r Round
	r.UUID = uuid.New().String()
	r.InvestigationUUID = investigationUUID
	r.Timestamp = TimestampNow()
	question, err := GetRandomQuestion()
	if err != nil {
		return r, err
	}
	r.Question = question

	err = saveRound(r)
	return r, err
}

func getRounds(investigationUUID string) ([]Round, error) {
	var rounds []Round
	log.Println("Getting rounds for investigation", investigationUUID)

	rows, err := database.Query("SELECT uuid, investigation_uuid, question_uuid, answer, timestamp FROM rounds WHERE investigation_uuid = $1 ORDER BY timestamp ASC", investigationUUID)
	if err != nil {
		log.Printf("Could not get rounds: %v\n", err)
		return rounds, err
	}
	defer rows.Close()

	for rows.Next() {
		var round Round
		err := rows.Scan(&round.UUID, &round.InvestigationUUID, &round.Question.UUID, &round.Answer, &round.Timestamp)
		if err != nil {
			log.Printf("Could not scan round: %v\n", err)
			return rounds, err
		}

		question, err := getQuestion(round.Question.UUID)
		if err != nil {
			log.Printf("Could not get question text for question_uuid=%s: %v", round.Question.UUID, err)
			return rounds, err
		}
		round.Question = question

		round.Eliminations, err = getEliminationsForRound(round.UUID)
		if err != nil {
			log.Printf("Could not get Eliminations for Round (%s): %v\n", round.UUID, err)
			return rounds, err
		}

		rounds = append(rounds, round)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v\n", err)
		return rounds, err
	}

	log.Println("Got rounds:", rounds)

	return rounds, nil
}

// MARK: ELIMINATION

type Elimination struct {
	UUID        string `json:"UUID"`
	RoundUUID   string `json:"RoundUUID"`
	SuspectUUID string `json:"SuspectUUID"`
	Timestamp   string `json:"Timestamp"`
}

const createEliminationsTable = `
	CREATE TABLE IF NOT EXISTS eliminations (
		UUID TEXT PRIMARY KEY,
		RoundUUID TEXT,
		SuspectUUID TEXT,
		Timestamp TEXT
	);`

// Save the Elimination, check if Criminal was not released
// and if not update the Game.Score accordingly.
func SaveElimination(suspectUUID, roundUUID, investigationUUID string) error {
	UUID := uuid.New().String()
	timestamp := TimestampNow()
	query := `INSERT OR REPLACE INTO eliminations (UUID, RoundUUID, SuspectUUID, Timestamp) VALUES (?, ?, ?, ?)`
	_, err := database.Exec(query, UUID, roundUUID, suspectUUID, timestamp)
	if err != nil {
		log.Printf("Could not save elimination of Suspect (%s) on Round (%s): %v\n", suspectUUID, roundUUID, err)
		return err
	}

	var criminalUUID string
	var gameUUID string
	row := database.QueryRow(`SELECT criminal_uuid, game_uuid FROM investigations WHERE uuid = $1`, investigationUUID)
	err = row.Scan(&criminalUUID, &gameUUID)
	if err != nil {
		log.Printf("Could not get criminal_uuid on Investigation (%s): %v\n", investigationUUID, err)
	}

	if criminalUUID != suspectUUID {

		increaseScore(gameUUID, roundUUID)
	} else {
		log.Println("Guilty criminal was released :(")
	}

	return nil
}

func getEliminationsForRound(roundUUID string) ([]Elimination, error) {
	var eliminations []Elimination
	log.Printf("Getting Eliminations for Round (%s)\n", roundUUID)

	rows, err := database.Query("SELECT UUID, RoundUUID, SuspectUUID, Timestamp FROM eliminations WHERE RoundUUID = $1 ORDER BY timestamp DESC", roundUUID)
	if err != nil {
		log.Printf("Could not get Eliminations: %v\n", err)
		return eliminations, err
	}
	defer rows.Close()

	for rows.Next() {
		var elimination Elimination
		err := rows.Scan(&elimination.UUID, &elimination.RoundUUID, &elimination.SuspectUUID, &elimination.Timestamp)
		if err != nil {
			log.Printf("Could not scan Elimination: %v\n", err)
			return eliminations, err
		}

		eliminations = append(eliminations, elimination)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during Eliminations rows iteration: %v\n", err)
		return eliminations, err
	}

	fmt.Println("GOT ELIMINATIONS:", eliminations)

	return eliminations, nil
}

// MARK: QUESTION

type Question struct {
	UUID    string `json:"UUID"`
	English string `json:"English"`
	Czech   string `json:"Czech"`
	Polish  string `json:"Polish"`
	Topic   string `json:"Topic"`
	Level   int    `json:"Level"`
}

var createQuestionsTable = `
	CREATE TABLE IF NOT EXISTS questions (
		UUID TEXT PRIMARY KEY,
		English TEXT,
		Czech TEXT,
		Polish TEXT,
		Topic TEXT,
		Level INT
	);`

func GetRandomQuestion() (Question, error) {
	var question Question
	row := database.QueryRow("SELECT UUID, English, Czech, Polish, Topic, Level FROM questions ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&question.UUID, &question.English, &question.Czech, &question.Polish, &question.Topic, &question.Level)
	return question, err
}

func InitQuestionsTable() error {
	_, err := database.Exec(createQuestionsTable)
	if err != nil {
		return err
	}
	for i := range defaultQuestions {
		err := SaveQuestion(defaultQuestions[i])
		if err != nil {
			log.Println("Cannot initialize Question:", err)
			return err
		}
	}

	return nil
}

// English is the cannonical text. If question with same English version exists, it will not overwrite.
func SaveQuestion(q Question) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM questions WHERE English = ?)"
	err := database.QueryRow(checkQuery, q.English).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	UUID := uuid.New().String()
	query := "INSERT into questions (UUID, English, Czech, Polish, Topic, Level) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = database.Exec(query, UUID, q.English, q.Czech, q.Polish, q.Topic, q.Level)
	if err != nil {
		log.Printf("Could not save Question %s (%s): %v", q.English, UUID, err)
		return err
	}

	return nil
}

func getQuestion(questionUUID string) (Question, error) {
	var question = Question{UUID: questionUUID}
	row := database.QueryRow("SELECT English, Czech, Polish, Topic, Level FROM questions WHERE UUID = $1 LIMIT 1", questionUUID)
	err := row.Scan(&question.English, &question.Czech, &question.Polish, &question.Topic, &question.Level)
	if err != nil {
		log.Printf("Could not scan question (%s): %v", questionUUID, err)
		return question, err
	}
	return question, nil
}

// MARK: ANSWER

// Save the Answer to the Round record in the database. There is then func WaitForAnswer()
// which is called from frontend once new Round is found (and so Question can be shown ASAP).
// But Answer takes time and when it is saved here the WaitForAnswer() retrieves it later.
func SaveAnswer(answer, roundUUID string) error {
	query := "UPDATE rounds SET answer = $1 WHERE uuid = $2"
	result, err := database.Exec(query, answer, roundUUID)
	if err != nil {
		log.Printf("Error updating answer for round %s: %v", roundUUID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for round %s: %v", roundUUID, err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("No rows were updated for round %s", roundUUID)
	}

	return nil
}

// MARK: LEVEL & SCORE

func getLevel(gameUUID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM investigations WHERE game_uuid = $1"

	err := database.QueryRow(query, gameUUID).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("error counting investigations records for game_uuid %s: %v", gameUUID, err)
	}

	return count, nil
}

// Increase the game.Score in the database after successful Elimination.
// Amount of increase is based on in which level we are and if it is 1st, 2nd or Nth
// Elimination in this round. Players are rewarded for risky behaviour - eliminating more than one suspect.
// But also they are rewarded for longevity - how much investigations they have solved.
func increaseScore(gameUUID string, roundUUID string) {
	level, err := getLevel(gameUUID)
	if err != nil {
		log.Println("Could not get level and increase score:", err)
		return
	}
	eliminations, err := getEliminationsForRound(roundUUID)
	if err != nil {
		log.Println("Could not get eliminations for this round and increase score:", err)
		return
	}

	amount := level * len(eliminations)

	query := "UPDATE games SET score = score + $1 WHERE uuid = $2"
	_, err = database.Exec(query, amount, gameUUID)
	if err != nil {
		log.Printf("error increasing score for gameUUID %s: %v", gameUUID, err)
		return
	}
	fmt.Printf("Score increased by %d\n", amount)
}

// This is used for High Scores list.
type FinalScore struct {
	Score        int    `json:"Score"`
	Position     int    `json:"Position"`
	Investigator string `json:"Investigator"`
	GameUUID     string `json:"GameUUID"`
	Timestamp    string `json:"Timestamp"`
}

func GetScores() ([]FinalScore, error) {
	var scores []FinalScore
	query := "SELECT uuid, score, investigator FROM games ORDER BY score DESC"
	rows, err := database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get scores: %w", err)
	}
	defer rows.Close()

	// Loop through the result set and scan into the games slice
	var position int
	for rows.Next() {
		position++
		var finalScore FinalScore
		var investigator sql.NullString
		var score sql.NullInt64
		err := rows.Scan(&finalScore.GameUUID, &score, &investigator)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		finalScore.Position = position
		finalScore.Investigator = investigator.String
		if score.Valid {
			finalScore.Score = int(score.Int64)
		}
		scores = append(scores, finalScore)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return scores, nil
}

func SaveScore(name, gameUUID string) error {
	query := "UPDATE games SET investigator = $1 WHERE uuid = $2"
	_, err := database.Exec(query, name, gameUUID)
	if err != nil {
		log.Printf("error saving investigator for gameUUID %s: %v", gameUUID, err)
	}
	return err
}

// MARK: AI SERVICES

type Service struct {
	Name  string `json:"Name"`
	Token string `json:"Token"`
}

const createServicesTable = `BEGIN;
CREATE TABLE IF NOT EXISTS services (
    Name TEXT PRIMARY KEY,
    Token TEXT
);
INSERT OR IGNORE INTO services (Name, Token)
VALUES
	('OpenAI', ''),
	('Anthropic', '');
COMMIT;` // list would be: ('OpenAI', ''), ('Google', ''), ('AWS', ''), ('Azure', '');

func GetService(name string) (Service, error) {
	var service Service
	query := "SELECT Name, Token FROM services WHERE name = $1"
	err := database.QueryRow(query, name).Scan(&service.Name, &service.Token)
	if err != nil {
		return service, fmt.Errorf("error geting Service for name %s: %v", name, err)
	}
	return service, nil
}

func GetServices() ([]Service, error) {
	var services []Service
	query := "SELECT Name, Token FROM services"
	rows, err := database.Query(query)
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var service Service
		err := rows.Scan(&service.Name, &service.Token)
		if err != nil {
			return services, err
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return services, err
	}
	return services, nil
}

func SaveToken(serviceName, token string) error {
	query := "UPDATE services SET Token = $1 WHERE Name = $2"
	result, err := database.Exec(query, token, serviceName)
	if err != nil {
		log.Printf("error saving Token for Service %s: %v", serviceName, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("SaveToken() rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("No record updated for service '%s'\n", serviceName)
		return err
	}
	fmt.Printf("SaveToken successful. Service=%s Token=%s\n", serviceName, token)
	return nil
}

// Wait until non-empty Answer appears on the Round record in Rounds table.
// Timeouts in 60 seconds, retries every 1 second. Answer is not modified anyhow,
// needs to be handled after return. On error returned answer is "". On error during
// generation of the answer, it is something like "failed OpenAI()". If everything was
// ok, the answer should be YES or NO, or lowercase variant - parse it later!
func WaitForAnswer(roundUUID string) string {
	pollInterval := 1 * time.Second
	timeout := 60 * time.Second
	start := time.Now()
	var answer string
	for {
		err := database.QueryRow("SELECT answer FROM rounds WHERE uuid = $1", roundUUID).Scan(&answer)
		if err == sql.ErrNoRows {
			log.Printf("Answer not available yet for Round (%s). Retrying...\n", roundUUID)
		} else if err != nil {
			log.Printf("Error querying answer for Round (%s), err: %v\n", roundUUID, err)
			return ""
		} else if answer == "" {
			log.Printf("Answer is still empty, lets sleep for a while...")
		} else {
			log.Printf("Answer found: %s", answer)
			return answer
		}
		if time.Since(start) > timeout {
			log.Printf("timed out waiting for answer to be available on Round (%s)\n", roundUUID)
			return ""
		}
		time.Sleep(pollInterval) // Wait for the polling interval before checking again
	}
}

// MARK: AI MODELS

type Model struct {
	Name    string `json:"Name"`
	Service string `json:"Service"`
	Active  bool   `json:"Active"`
}

const createModelsTable = `BEGIN;
CREATE TABLE IF NOT EXISTS models (
    Name TEXT PRIMARY KEY,
	Service TEXT,
    Active INT
);
INSERT OR IGNORE INTO models (Name, Service, Active)
VALUES
	('gpt-4o-2024-08-06', 'OpenAI', '1'),
	('chatgpt-4o-latest', 'OpenAI', '0'),
	('gpt-4o-mini-2024-07-18', 'OpenAI', '0'),
	('claude-3-haiku-20240307', 'Anthropic', '0'),
	('claude-3-5-sonnet-20240620', 'Anthropic', '0');
COMMIT;` // as defined in ai.go:supportedModels - hacky way, but OK for now

func GetActiveModel() (Model, error) {
	var model Model
	query := "SELECT Name, Service, Active FROM models WHERE Active = 1"
	err := database.QueryRow(query).Scan(&model.Name, &model.Service, &model.Active)
	if err != nil {
		return model, fmt.Errorf("error geting active Model: %v", err)
	}
	return model, nil
}

func GetAllModels() ([]Model, error) {
	var models []Model
	query := "SELECT Name, Service, Active FROM models"
	rows, err := database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving models: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var model Model
		err := rows.Scan(&model.Name, &model.Service, &model.Active)
		if err != nil {
			return nil, fmt.Errorf("error scanning model: %v", err)
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return models, nil
}

func SetActiveModel(modelName string) error {
	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deactivateQuery := "UPDATE models SET Active = 0"
	_, err = tx.Exec(deactivateQuery)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deactivating models: %v", err)
	}

	activateQuery := "UPDATE models SET Active = 1 WHERE Name = ?"
	_, err = tx.Exec(activateQuery, modelName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error activating model %s: %v", modelName, err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// MARK: DESCRIPTIONS

// Holds description of the Suspect image. There can be multiple descriptions for one Suspect.
// Descriptions can be made by different Services and different Models.
type Description struct {
	UUID        string `json:"UUID"`
	SuspectUUID string `json:"SuspectUUID"`
	Service     string `json:"Service"`
	Model       string `json:"Model"`
	Description string `json:"Description"`
	Prompt      string `json:"Prompt"`
	Timestamp   string `json:"Timestamp"`
}

const createDescriptionsTable = `
	CREATE TABLE IF NOT EXISTS descriptions (
		UUID TEXT PRIMARY KEY,
		SuspectUUID TEXT,
		Service TEXT,
		Model TEXT,
		Description TEXT,
		Prompt TEXT,
		Timestamp TEXT
	);`

func SaveDescription(d Description) error {
	query := `
		INSERT OR REPLACE INTO descriptions (UUID, SuspectUUID, Service, Model, Description, Prompt, Timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	timestamp := TimestampNow()
	if d.UUID == "" {
		d.UUID = uuid.New().String()
	}
	_, err := database.Exec(query, d.UUID, d.SuspectUUID, d.Service, d.Model, d.Description, d.Prompt, timestamp)
	return err
}

func GetDescriptionsForSuspect(suspectUUID, service, model string) ([]Description, error) {
	var descriptions []Description
	query := "SELECT UUID, Description, Prompt, Timestamp FROM descriptions WHERE SuspectUUID = $1 AND Service = $2 AND Model = $3"
	rows, err := database.Query(query, suspectUUID, service, model)
	if err != nil {
		return nil, fmt.Errorf("failed to get descriptions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d = Description{
			SuspectUUID: suspectUUID,
			Service:     service,
			Model:       model,
		}
		err := rows.Scan(&d.UUID, &d.Description, &d.Prompt, &d.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan description row: %w", err)
		}
		descriptions = append(descriptions, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("description rows iteration error: %w", err)
	}

	return descriptions, nil
}
