package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the checkCat func as a callback for MessageCreate events.
	bot.AddHandler(checkCat)

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("I'm logged in ! (Press CTRL-C to exit.)\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func contains(arr []string, value string) bool {
	for _, elem := range arr {
		if elem == value {
			return true
		}
	}
	return false
}

func isCatAlive() bool {
	return rand.Float64() > 0.5
}

func checkCat(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Put the message in lowercase
	recv := strings.ToLower(m.Content)

	// Sentence we should check the cat
	msg := []string{
		"is the cat alive ?",
		"is the cat alive",
		"alive",
		"alive ?",
		"the cat is alive ?",
		"cat alive ?",
	}

	// Ignore our messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Answer (or not) the user
	if contains(msg, recv) {
		fmt.Print(
			m.Author.Username,
			"#",
			m.Author.Discriminator,
			" want to know if the cat is alive.",
		)
		// Check if the cat is alive
		if isCatAlive() {
			s.ChannelMessageSend(m.ChannelID, "Meow !")
			fmt.Println(" And he is !")
		} else {
			fmt.Println(" Seems not...")
		}
	}
}
