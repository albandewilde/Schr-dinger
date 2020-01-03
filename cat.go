package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Read token in the `secrets.json` file
	secretFile, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		fmt.Println("Error while reading secrets:", err)
	}

	type Secrets struct {
		DISCORD string
	}

	var secrets Secrets

	// Parse json content
	err = json.Unmarshal(secretFile, &secrets)
	if err != nil {
		fmt.Println("Error while parsing secrets:", err)
	}

	// Create a new Discord session using the provided bot token.
	bot, err := discordgo.New("Bot " + secrets.DISCORD)
	if err != nil {
		fmt.Println("Error while creating the Discord session,", err)
		return
	}

	// Register the checkCat func as a callback for MessageCreate events.
	bot.AddHandler(checkCat)

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("Error while opening connection,", err)
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
			"\033[33m",
			time.Now().Format("2006-01-02 15:04:05"),
			" ",
			"\033[34m",
			m.Author.Username,
			"#",
			m.Author.Discriminator,
			"\033[0m",
			" want to know if the cat is alive.",
		)
		// Check if the cat is alive
		if isCatAlive() {
			s.ChannelMessageSend(m.ChannelID, "Meow !")
			fmt.Println(
				"\033[32m",
				"And he is !",
				"\033[0m",
			)
		} else {
			fmt.Println(
				"\033[31m",
				"Seems not...",
				"\033[0m",
			)
		}
	}
}
