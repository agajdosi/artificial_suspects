package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"

	"golang.org/x/exp/rand"
)

// Holds the database connection. Is created via EnsureDBAvailable()
var database *sql.DB

const (
	appName           = "suspects"
	TimeFormat string = time.RFC3339Nano
	numSuspect        = 15 // How many suspects are in one investigation - there were 12 in original board game.
)

// MARK: APP HANDLERS

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

// Create a new game. Returns the suspects.
func (a *App) NewGame() Game {
	game, err := newGame()
	if err != nil {
		fmt.Println("NewGame() error:", err)
	}
	return game
}

// Loads the last game.
func (a *App) GetGame() Game {
	game, err := getCurrentGame()
	if err != nil {
		fmt.Println("GetGame()->getCurrentGame() error: ", err)
	}
	return game
}

func (a *App) NextInvestigation() Game {
	game, err := getCurrentGame()
	if err != nil {
		log.Printf("NextInvestigation() could not get current game: %v\n", err)
	}
	game.Investigation, err = newInvestigation(game.UUID)
	if err != nil {
		log.Printf("NextInvestigation() could not get new investigation: %v\n", err)
	}
	return game
}

// Next level is requested. Updates the question and level for the game object.
func (a *App) NextRound() Game {
	game, err := getCurrentGame()
	if err != nil {
		log.Printf("NextRound() could not get current game: %v\n", err)
	}

	round, err := newRound(game.Investigation.UUID)
	if err != nil {
		log.Printf("NextRound() could not get new Round: %v\n", err)
	}
	go GetAnswerFromAI(round, game.Investigation.CriminalUUID)

	game.Investigation.Rounds = append(game.Investigation.Rounds, round) // prepend

	fmt.Printf("New Round %d: %s\n", game.Level, game.Investigation.Rounds[len(game.Investigation.Rounds)-1].Question)
	return game
}

func (a *App) GetScores() []FinalScore {
	scores, err := getScores()
	if err != nil {
		log.Println("GetScores()", err)
	}
	return scores
}

// Wait until the Answer from AI is present in the database.
// TODO: implement the actuall retrieval from the DB.
func (a *App) WaitForAnswer(roundUUID string) string {
	pollInterval := 2 * time.Second
	timeout := 30 * time.Second
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			log.Printf("timed out waiting for answer to be available on Round (%s)\n", roundUUID)
			return ""
		}

		var answer string
		err := database.QueryRow("SELECT answer FROM rounds WHERE uuid = $1", roundUUID).Scan(&answer)
		if err == sql.ErrNoRows {
			log.Printf("Answer not available yet for Round (%s). Retrying...\n", roundUUID)
		} else if err != nil {
			log.Printf("Error querying answer for Round (%s), err: %v\n", roundUUID, err)
			return ""
		} else {
			return answer
		}

		// Wait for the polling interval before checking again
		time.Sleep(pollInterval)
	}
}

// User selected Suspect to be eliminated from the Investigation.
// This will save the Elimination to Eliminatios table. Also it
func (a *App) EliminateSuspect(suspectUUID, roundUUID, investigationUUID string) {
	fmt.Printf(">>> Eliminating suspect (%s) on Investigation (%s) in Round (%s)\n", suspectUUID, investigationUUID, roundUUID)
	err := saveElimination(suspectUUID, roundUUID, investigationUUID)
	if err != nil {
		log.Printf("EliminateSuspect() error: %v\n", err)
	}
}

func (a *App) SaveScore(name, gameUUID string) {
	fmt.Println("SaveScore:", name, gameUUID)
	query := "UPDATE games SET investigator = $1 WHERE uuid = $2"
	_, err := database.Exec(query, name, gameUUID)
	if err != nil {
		log.Printf("error saving investigator for gameUUID %s: %v", gameUUID, err)
	}
}

func (a *App) GetServices() []Service {
	services, err := getServices()
	if err != nil {
		log.Println("Could not get services:", err)
	}

	fmt.Println("Got Services:", services)

	return services
}

func (a *App) SaveToken(serviceName, token string) {
	query := "UPDATE services SET Token = $1 WHERE Name = $2"
	result, err := database.Exec(query, token, serviceName)
	if err != nil {
		log.Printf("error saving Token for Service %s: %v", serviceName, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("SaveToken() rows affected: %v", err)
		return
	}
	if rowsAffected == 0 {
		log.Printf("No record updated for service '%s'\n", serviceName)
		return
	}
	fmt.Printf("SaveToken successful. Service=%s Token=%s\n", serviceName, token)
}

// MARK: DATABASE

func GetDataDirPath() string {
	return filepath.Join(xdg.ConfigHome, appName)
}

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

	fmt.Println("Stating the dbPATH")
	_, err := os.Stat(gameDBPath)
	if os.IsNotExist(err) {
		file, err := os.Create(gameDBPath)
		if err != nil {
			return err
		}
		file.Close()
	}

	fmt.Println("gameDBPath created")

	db, err := sql.Open("sqlite3", gameDBPath)
	if err != nil {
		log.Fatal(err)
	}
	database = db // setting to global variable

	err = initDB()
	if err != nil {
		return err
	}

	fmt.Println("Database prepared and available!")

	return nil
}

func initDB() error {
	var tables = []string{
		createGamesTable,
		createInvestigationsTable,
		createRoundsTable,
		createEliminationsTable,
		createServicesTable,
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

	err = InitSuspectsTable()
	if err != nil {
		return err
	}

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

func InitSuspectsTable() error {
	_, err := database.Exec(createSuspectsTable)
	if err != nil {
		return err
	}

	for i := range defaultSuspects {
		err := SaveSuspect(defaultSuspects[i])
		if err != nil {
			log.Println("Cannot initialize Suspect:", err)
			return err
		}
	}

	return nil
}

func SaveSuspect(suspect Suspect) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM suspects WHERE image = ?)"
	err := database.QueryRow(checkQuery, suspect.Image).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	UUID := uuid.New().String()
	timestamp := time.Now().String()
	query := "INSERT into suspects (uuid, image, timestamp) VALUES (?, ?, ?)"
	_, err = database.Exec(query, UUID, suspect.Image, timestamp)
	if err != nil {
		log.Printf("Could not save Suspect %s (%s): %v", suspect.Image, UUID, err)
		return err
	}

	return nil
}

// Get the basic suspect data from the Database without Suspect.Free field!
// Because Suspect.Free and Suspect.Fled needs information from table Investigation->Rounds->Eliminations.
func getSuspect(suspectUUID string) (Suspect, error) {
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
func getSuspects(suspectUUIDs []string, investigation Investigation) ([]Suspect, error) {
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
		suspect, err = getSuspect(suspectUUIDs[x])
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

func newGame() (Game, error) {
	var game Game
	game.UUID = uuid.New().String()
	game.Timestamp = time.Now().String()
	game.Score = 0
	err := saveGame(game)
	if err != nil {
		return game, err
	}

	game.Investigation, err = newInvestigation(game.UUID)
	if err != nil {
		return game, err
	}
	game.Level, err = getLevel(game.UUID)
	if err != nil {
		return game, err
	}

	GetAnswerFromAI(game.Investigation.Rounds[0], game.Investigation.CriminalUUID)

	return game, err
}

func getCurrentGame() (Game, error) {
	var game Game
	row := database.QueryRow("SELECT uuid, timestamp, score FROM games ORDER BY timestamp DESC LIMIT 1")
	err := row.Scan(&game.UUID, &game.Timestamp, &game.Score)

	// No game found - first play
	if err == sql.ErrNoRows {
		log.Println("Warning: No games in DB, creating new game")
		return newGame()
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
	query := `INSERT INTO games (uuid, timestamp, score) VALUES (?, ?, ?)`
	_, err := database.Exec(query, game.UUID, game.Timestamp, game.Score)
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
func newInvestigation(gameUUID string) (Investigation, error) {
	var i Investigation
	i.UUID = uuid.New().String()
	i.GameUUID = gameUUID
	i.Timestamp = time.Now().String()

	round, err := newRound(i.UUID)
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

	investigation.Suspects, err = getSuspects(suspects_uuids, investigation)
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
	QuestionUUID      string        `json:"QuestionUUID"`
	Question          string        `json:"question"` // TODO: Question could be actually the whole object
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
	_, err := database.Exec(query, r.UUID, r.InvestigationUUID, r.QuestionUUID, r.Answer, r.Timestamp)
	return err
}

func newRound(investigationUUID string) (Round, error) {
	var r Round
	r.UUID = uuid.New().String()
	r.InvestigationUUID = investigationUUID
	r.Timestamp = time.Now().Format(TimeFormat)
	question, err := GetRandomQuestion()
	if err != nil {
		return r, err
	}
	r.Question = question.Text
	r.QuestionUUID = question.UUID

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
		err := rows.Scan(&round.UUID, &round.InvestigationUUID, &round.QuestionUUID, &round.Answer, &round.Timestamp)
		if err != nil {
			log.Printf("Could not scan round: %v\n", err)
			return rounds, err
		}

		question, err := getQuestion(round.QuestionUUID)
		if err != nil {
			log.Printf("Could not get question text for question_uuid=%s: %v", round.QuestionUUID, err)
			return rounds, err
		}
		round.Question = question.Text

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
func saveElimination(suspectUUID, roundUUID, investigationUUID string) error {
	UUID := uuid.New().String()
	timestamp := time.Now().Format(TimeFormat)
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
	UUID  string `json:"uuid"`
	Text  string `json:"text"`
	Topic string `json:"topic"`
	Level int    `json:"level"`
}

var createQuestionsTable = `
	CREATE TABLE IF NOT EXISTS questions (
		uuid TEXT PRIMARY KEY,
		question TEXT,
		topic TEXT,
		level INT
	);`

func GetRandomQuestion() (Question, error) {
	var question Question
	row := database.QueryRow("SELECT uuid, question, topic, level FROM questions ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&question.UUID, &question.Text, &question.Topic, &question.Level)
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

func SaveQuestion(q Question) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM questions WHERE question = ?)"
	err := database.QueryRow(checkQuery, q.Text).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	UUID := uuid.New().String()
	query := "INSERT into questions (uuid, question, topic, level) VALUES (?, ?, ?, ?)"
	_, err = database.Exec(query, UUID, q.Text, q.Topic, q.Level)
	if err != nil {
		log.Printf("Could not save Question %s (%s): %v", q.Text, UUID, err)
		return err
	}

	return nil
}

func getQuestion(questionUUID string) (Question, error) {
	var question = Question{UUID: questionUUID}
	row := database.QueryRow("SELECT question, topic, level FROM questions WHERE uuid = $1 LIMIT 1", questionUUID)
	err := row.Scan(&question.Text, &question.Topic, &question.Level)
	if err != nil {
		log.Printf("Could not scan question (%s): %v", questionUUID, err)
		return question, err
	}
	return question, nil
}

// MARK: ANSWER

// Get the Answer to Question from the AI model and save it into the database.
// Does not return anything, for retrieval App later uses WaitForAnswer().
// TODO: actually implement this.
func GetAnswerFromAI(round Round, criminalUUID string) {
	fmt.Println(">>> GetAnswerFromAI called!")
	var answer string
	n := rand.Intn(2)
	if n%2 == 0 {
		answer = "yes"
	} else {
		answer = "no"
	}

	fmt.Println("Answer is:", answer)
	SaveAnswer(answer, round.UUID)
}

func GetAnswerFromOpenAI(question Question, criminalUUID string) (string, error) {
	prompt := `ROLE: You are a player of Unusual Suspects. Even though you cannot normally help with identifying or analyzing personal traits,
you can for purposes of simple and funny game.
TASK: You take a look on the image of a person and try to describe them. Not just visually, but also sociologically, economically, politically.
Are they old or young, female, male or nonbinary? Do they conservative, wealthy, or poor, liberal? Are they educated, do they work with hands?
Who are they? How they behave?
You generate a 200 words contemplating about what are the aspects of the person you see.
Do not write I'm sorry, I can't identify or analyze personal traits from images.
`
	fmt.Println("Sending prompt to AI:", prompt)

	service, err := getService("openai")
	if err != nil {
		return "", err
	}

	client := openai.NewClient(service.Token)
	// Call the OpenAI API with the image and question
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL:    "https://io.google/2023/speakers/200d922d-3cba-4ff9-8682-5a1e9b57a1d9_500.webp",
								Detail: openai.ImageURLDetailHigh,
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		return "", err
	}
	fmt.Println("RESPONSE:", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

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

func getScores() ([]FinalScore, error) {
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

// MARK: AI MODELS

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
	('OpenAI', '');
COMMIT;` // list would be: ('OpenAI', ''), ('Google', ''), ('AWS', ''), ('Azure', '');

func getService(name string) (Service, error) {
	var service Service
	query := "SELECT Name, Token FROM services WHERE name = $1"
	err := database.QueryRow(query, name).Scan(&service.Name, &service.Token)
	if err != nil {
		return service, fmt.Errorf("error geting Service for name %s: %v", name, err)
	}
	return service, nil
}

func getServices() ([]Service, error) {
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
