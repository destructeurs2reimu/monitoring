package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	var err error

	token := os.Getenv("TOKEN")
	s, err = discordgo.New(token)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	err := s.Open()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Logged in as %s (%s)\n", s.State.User.DisplayName(), s.State.User.ID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}