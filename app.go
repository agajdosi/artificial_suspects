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
var db *sql.DB

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
	question := GetQuestion()

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
				{Question: question},
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
	game.Investigation.Rounds[0].Question = GetQuestion() // HERE APPEND A NEW ROUND
	game.Level++
	fmt.Printf("New level %d: %s\n", game.Level, game.Investigation.Rounds[0].Question)
	return game
}

// Asks the AI whether it thinks the
func (a *App) GetAnswerFromAI() bool {
	return true
}

// User selected suspect to be freed.
func (a *App) FreeSuspect(suspectUUID string) bool {
	fmt.Printf("Freeing suspect: %s\n", suspectUUID)
	return rand.Intn(2) == 1
}

// MARK: QUESTION

func GetQuestion() string {
	i := rand.Intn(len(Questions)) // DUMMY
	return Questions[i]
}

var Questions = []string{
	"Does the suspect love pizza?",
	"Does the suspect hate immigrants?",
	"Is the suspect a leftist?",
	"Does the suspect love spicy food?",
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

	database, err := sql.Open("sqlite3", gameDBPath)
	if err != nil {
		log.Fatal(err)
	}

	err = initDB(database)
	if err != nil {
		return err
	}

	db = database // setting to global variable
	return nil
}

func initDB(db *sql.DB) error {
	var tables = []string{createSuspectsTable, createGamesTable, createInvestigationsTable, createRoundsTable, createEliminationsTable}
	for i := range tables {
		_, err := db.Exec(tables[i])
		if err != nil {
			return err
		}
	}
	return nil
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
	rows, err := db.Query("SELECT uuid FROM games")
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
		questionUUID TEXT,
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
