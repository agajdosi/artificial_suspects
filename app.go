package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/exp/rand"
)

var game Game

// Holds the database connection. Is created via EnsureDBAvailable()
var database *sql.DB

const (
	appName = "suspects"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Create a new game. Returns the suspects.
func (a *App) NewGame() Game {
	fmt.Println("Loading new game")
	question, err := GetRandomQuestion()
	if err != nil {
		fmt.Println("NewGame()->GetRandomQuestion() error:", err)
	}

	// DUMMY
	game = Game{
		UUID:  uuid.New().String(),
		Level: 1,
		Investigation: Investigation{
			Suspects: []Suspect{
				{UUID: "1", ImageSource: "2.jpg"},
				{UUID: "2", ImageSource: "2.jpg"},
				{UUID: "3", ImageSource: "1.jpg"},
				{UUID: "4", ImageSource: "2.jpg"},
				{UUID: "5", ImageSource: "2.jpg"},
				{UUID: "6", ImageSource: "2.jpg"},
				{UUID: "7", ImageSource: "1.jpg"},
				{UUID: "8", ImageSource: "2.jpg"},
				{UUID: "9", ImageSource: "1.jpg"},
				{UUID: "10", ImageSource: "2.jpg"},
				{UUID: "11", ImageSource: "1.jpg"},
				{UUID: "12", ImageSource: "2.jpg"},
				{UUID: "13", ImageSource: "1.jpg"},
				{UUID: "14", ImageSource: "2.jpg"},
				{UUID: "15", ImageSource: "1.jpg"},
			},
			Rounds: []Round{
				{Question: question.Text},
			},
		},
	}
	return game
}

// Loads the last game.
func (a *App) GetGame() Game {
	game := a.NewGame()
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

func GetRandomQuestion() (*Question, error) {
	var question Question
	row := database.QueryRow("SELECT uuid, question, topic, level FROM questions ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&question.UUID, &question.Text, &question.Topic, &question.Level)
	return &question, err
}

var QS = []Question{
	{Text: "Does the suspect love pizza?", Topic: "food", Level: 1},
	{Text: "Does the suspect love spicy food?", Topic: "food", Level: 1},
	{Text: "Does the suspect hate immigrants?", Topic: "political", Level: 1},
	{Text: "Is the suspect a leftist?", Topic: "political", Level: 1},
}

func InitQuestionsTable(db *sql.DB) error {
	_, err := db.Exec(createQuestionsTable)
	if err != nil {
		return err
	}

	for i := range QS {
		err := SaveQuestion(db, QS[i])
		if err != nil {
			log.Println("Cannot initialize question:", err)
			return err
		}
	}

	return nil
}

func SaveQuestion(db *sql.DB, q Question) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM questions WHERE question = ?)"
	err := db.QueryRow(checkQuery, q.Text).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	UUID := uuid.New().String()
	query := "INSERT into questions (uuid, question, topic, level) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, UUID, q.Text, q.Topic, q.Level)
	return err
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

	err = initDB(db)
	if err != nil {
		return err
	}

	database = db // setting to global variable
	fmt.Println("Database prepared and available!")

	return nil
}

func initDB(db *sql.DB) error {
	var tables = []string{
		createSuspectsTable,
		createGamesTable,
		createInvestigationsTable,
		createRoundsTable,
		createEliminationsTable,
	}
	for i := range tables {
		_, err := db.Exec(tables[i])
		if err != nil {
			fmt.Printf("Error initializing table: '%s', error: %v", tables[i], err)
			return err
		}
	}
	err := InitQuestionsTable(db)

	return err
}

// MARK: SUSPECTS

type Suspect struct {
	UUID        string `json:"uuid"`
	ImageSource string `json:"imageSource"`
	Free        bool   `json:"free"`
}

const createSuspectsTable = `
	CREATE TABLE IF NOT EXISTS suspects (
		uuid TEXT PRIMARY KEY,
		image TEXT,
		timestamp TEXT
	);`

// MARK: GAMES

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

// Get all users from the database.
func GetGame() ([]*Game, error) {
	var users []*Game
	rows, err := database.Query("SELECT uuid FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(Game)
		err := rows.Scan(
			&user.UUID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// MARK: CASES

// Investigation is a set of X Suspects, User needs to find a Criminal among them.
type Investigation struct {
	UUID      string    `json:"uuid"`
	Suspects  []Suspect `json:"suspects"`
	Level     int       `json:"level"` // but can be taken from len of Rounds
	Rounds    []Round   `json:"rounds"`
	Timestamp string
}

const createInvestigationsTable = `
	CREATE TABLE IF NOT EXISTS investigations (
		uuid TEXT PRIMARY KEY,
		game_uuid TEXT,
		timestamp TEXT
	);`

// MARK: ROUNDS

type Round struct {
	UUID      string `json:"uuid"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Timestamp string
}

const createRoundsTable = `
	CREATE TABLE IF NOT EXISTS rounds (
		uuid TEXT PRIMARY KEY,
		investigation_uuid TEXT,
		question_uuid TEXT,
		answer TEXT,
		timestamp TEXT
	);`

// MARK: ELIMINATIONS

type Elimination struct {
	UUID        string `json:"uuid"`
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
