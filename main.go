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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

var s *discordgo.Session

func init() {
	var err error
	token := os.Getenv("TOKEN")
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}
}

var db *sql.DB

func init() {
	var err error
	var stmt *sql.Stmt

	db, err = sql.Open("sqlite3", "monitoring.db")
	if err != nil {
		log.Fatalln(err)
	}
	
	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS monitoring (
			date TIMESTAMP PRIMARY KEY DEFAULT CURRENT_TIMESTAMP NOT NULL,
			server_id INTEGER NOT NULL,
			member_count INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatalln(err)
	}
	
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}

var insertMemberCountStmt *sql.Stmt

func init() {
	var err error
	insertMemberCountStmt, err = db.Prepare(`
		INSERT INTO monitoring (server_id, member_count)
		VALUES (?, ?)
		ON CONFLICT (date)
		DO UPDATE SET member_count = excluded.member_count
	`)
	if err != nil {
		log.Fatalln(err)
	}
}

var deleteServerRowsStmt *sql.Stmt

func init() {
	var err error
	deleteServerRowsStmt, err = db.Prepare(`
		DELETE FROM monitoring WHERE server_id = (?)
	`)
	if err != nil {
		log.Fatalln(err)
	}
}

var commands = []*discordgo.ApplicationCommand{}
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {}

func createCommands() (err error) {
	for _, v := range commands {
		_, err = s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return
		}
	}
	return
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as %s <@%s>\n", s.State.User.DisplayName(), s.State.User.ID)
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	h, ok := commandHandlers[name]
	if !ok {
		return
	}

	h(s, i)
}

func onGuildCreate(_ *discordgo.Session, g *discordgo.GuildCreate) {
	_, err := insertMemberCountStmt.Exec(g.ID, g.MemberCount)
	if err != nil {
		log.Fatalln(err)
	}
}

func onGuildDelete(_ *discordgo.Session, g *discordgo.GuildDelete) {
	_, err := deleteServerRowsStmt.Exec(g.ID)
	if err != nil {
		log.Fatalln(err)
	}
}

func onGuildMemberEvent(s *discordgo.Session, m *discordgo.Member) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = insertMemberCountStmt.Exec(m.GuildID, g.MemberCount)
	if err != nil {
		log.Fatalln(err)
	}
}

func onGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	onGuildMemberEvent(s, m.Member)
}

func onGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	onGuildMemberEvent(s, m.Member)
}

func main() {
	var version string
	err := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Using sqlite %s", version)
	
	s.AddHandler(onReady)
	s.AddHandler(onInteractionCreate)
	s.AddHandler(onGuildCreate)
	s.AddHandler(onGuildDelete)
	s.AddHandler(onGuildMemberAdd)
	s.AddHandler(onGuildMemberRemove)

	err = s.Open()
	if err != nil {
		log.Fatalln(err)
	}

	createCommands()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
