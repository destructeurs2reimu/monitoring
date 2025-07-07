package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var s *discordgo.Session

var db *sql.DB

var commands = []*discordgo.ApplicationCommand{}
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	var err error

	token := os.Getenv("TOKEN")
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	var err error

	db, err = sql.Open("sqlite3", "monitoring.db")
	if err != nil {
		log.Fatalln(err)
	}
}

func createCommands() (err error) {
	for _, v := range commands {
		_, err = s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return
		}
	}
	return
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	h, ok := commandHandlers[name]
	if !ok {
		return
	}

	h(s, i)
}

func main() {
	var version string
	err := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Using sqlite %s", version)

	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS books (
			server_id INT PRIMARY KEY,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			member_count INT
		)
	`)
	if err != nil {
		log.Fatalln(err)
	}
	statement.Exec()

	err = s.Open()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Logged in as %s <@%s>\n", s.State.User.DisplayName(), s.State.User.ID)

	createCommands()
	s.AddHandler(commandHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
