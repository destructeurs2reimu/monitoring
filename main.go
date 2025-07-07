package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var s *discordgo.Session

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
	err := s.Open()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Logged in as %s (%s)\n", s.State.User.DisplayName(), s.State.User.ID)

	createCommands()
	s.AddHandler(commandHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
