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

	"golang.org/x/exp/rand"
)

// Holds the database connection. Is created via EnsureDBAvailable()
var database *sql.DB

const (
	appName           = "suspects"
	TimeFormat string = time.RFC3339Nano
	numSuspect        = 15 // How many suspects are in one investigation - there were 12 in original board game.
)

// MARK: APP

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

// Next level is requested. Updates the question and level for the game object.
func (a *App) NextLevel() Game {
	var game Game

	q, err := GetRandomQuestion()
	if err != nil {
		fmt.Println("GetRandomQuestion() error:", err)
		return game
	}

	game.Investigation.Rounds[0].Question = q.Text
	game.Level++
	fmt.Printf("New level %d: %s\n", game.Level, game.Investigation.Rounds[0].Question)
	return game
}

// Asks the AI whether it thinks the
func (a *App) GetAnswerFromAI() bool {
	return true
}

// User selected suspect to be freed.
func (a *App) FreeSuspect(uuid string) bool {
	fmt.Printf("Freeing suspect: %s\n", uuid)
	return rand.Intn(2) == 1
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

var QS = []Question{
	{Text: "Does the suspect love pizza?", Topic: "food", Level: 1},
	{Text: "Does the suspect love spicy food?", Topic: "food", Level: 1},
	{Text: "Does the suspect hate immigrants?", Topic: "political", Level: 1},
	{Text: "Is the suspect a leftist?", Topic: "political", Level: 1},
}

func InitQuestionsTable() error {
	_, err := database.Exec(createQuestionsTable)
	if err != nil {
		return err
	}

	for i := range QS {
		err := SaveQuestion(QS[i])
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

// MARK: SUSPECTS

type Suspect struct {
	UUID      string `json:"UUID"`
	Image     string `json:"Image"`
	Free      bool   `json:"Free"`
	Timestamp string `json:"Timestamp"`
}

const createSuspectsTable = `
	CREATE TABLE IF NOT EXISTS suspects (
		uuid TEXT PRIMARY KEY,
		image TEXT,
		timestamp TEXT
	);`

var defaultSuspects = []Suspect{
	{UUID: "001", Image: "1.jpg"},
	{UUID: "002", Image: "2.jpg"},
	{UUID: "003", Image: "3.jpg"},
	{UUID: "004", Image: "4.jpg"},
	{UUID: "005", Image: "5.jpg"},
	{UUID: "006", Image: "6.jpg"},
	{UUID: "007", Image: "7.jpg"},
	{UUID: "008", Image: "8.jpg"},
	{UUID: "009", Image: "9.jpg"},
	{UUID: "010", Image: "10.jpg"},
	{UUID: "011", Image: "11.jpg"},
	{UUID: "012", Image: "12.jpg"},
	{UUID: "013", Image: "13.jpg"},
	{UUID: "014", Image: "14.jpg"},
	{UUID: "015", Image: "15.jpg"},
}

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

// DUMMY for now
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
type Game struct {
	UUID          string        `json:"uuid"`
	Level         int           `json:"level"`
	Investigation Investigation `json:"investigation"`
	Timestamp     string
}

const createGamesTable = `
	CREATE TABLE IF NOT EXISTS games (
		uuid TEXT PRIMARY KEY,
		timestamp TEXT
	);`

func newGame() (Game, error) {
	var game Game
	game.UUID = uuid.New().String()
	game.Timestamp = time.Now().String()
	err := saveGame(game)
	if err != nil {
		return game, err
	}

	game.Investigation, err = newInvestigation(game.UUID)
	if err != nil {
		return game, err
	}

	return game, err
}

func getCurrentGame() (Game, error) {
	var game Game
	row := database.QueryRow("SELECT uuid, timestamp FROM games ORDER BY timestamp DESC LIMIT 1")
	err := row.Scan(&game.UUID, &game.Timestamp)

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

	return game, nil
}

func saveGame(game Game) error {
	query := `INSERT INTO games (uuid, timestamp)VALUES (?, ?)`
	_, err := database.Exec(query, game.UUID, game.Timestamp)
	return err
}

// MARK: INVESTIGATION

// Investigation is a set of X Suspects, User needs to find a Criminal among them.
type Investigation struct {
	UUID      string    `json:"uuid"`
	GameUUID  string    `json:"game_uuid"`
	Suspects  []Suspect `json:"suspects"`
	Level     int       `json:"level"` // but can be taken from len of Rounds
	Rounds    []Round   `json:"rounds"`
	Criminal  string
	Timestamp string
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
		investigation.Criminal,
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
	i.Level = 1

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
	i.Criminal = i.Suspects[cn].UUID

	err = saveInvestigation(i)
	return i, err
}

func getCurrentInvestigation(gameUUID string) (Investigation, error) {
	var investigation = Investigation{GameUUID: gameUUID}
	log.Printf("Getting investigation for game %s\n", gameUUID)
	row := database.QueryRow("SELECT uuid, timestamp FROM investigations WHERE game_uuid = $1 ORDER BY timestamp DESC LIMIT 1", gameUUID)
	err := row.Scan(&investigation.UUID, &investigation.Timestamp)
	if err != nil {
		log.Printf("Could not get investigation: %v\n", err)
		return investigation, err
	}

	investigation.Rounds, err = getRounds(investigation.UUID)
	if err != nil {
		return investigation, err
	}

	return investigation, nil
}

// MARK: ROUNDS

type Round struct {
	UUID              string `json:"uuid"`
	InvestigationUUID string
	QuestionUUID      string
	Question          string `json:"question"` // TODO: Question could be actually the whole object
	Answer            string `json:"answer"`   // TODO: Answer could be actually stored in table
	Timestamp         string
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

	rows, err := database.Query("SELECT uuid, investigation_uuid, question_uuid, answer, timestamp FROM rounds WHERE investigation_uuid = $1 ORDER BY timestamp DESC", investigationUUID)
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
		}
		round.Question = question.Text

		rounds = append(rounds, round)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v\n", err)
		return rounds, err
	}

	log.Println("Got rounds:", rounds)

	return rounds, nil
}

// MARK: ELIMINATIONS

type Elimination struct {
	UUID        string `json:"uuid"`
	RoundUUID   string
	SuspectUUID string `json:"suspectUUID"`
	Timestamp   string
}

const createEliminationsTable = `
	CREATE TABLE IF NOT EXISTS eliminations (
		uuid TEXT PRIMARY KEY,
		round_uuid TEXT,
		suspect_uuid TEXT,
		timestamp TEXT
	);`

func saveElimination(e Elimination) error {
	query := `
		INSERT OR REPLACE INTO eliminations (uuid, round_uuid, suspect_uuid, timestamp)
		VALUES (?, ?, ?, ?, ?)
		`
	_, err := database.Exec(query, e.UUID, e.RoundUUID, e.SuspectUUID, e.Timestamp)
	return err
}
