package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var Config struct {
	Token string `json:"token"`
	Delay int    `json:"delay"`
}

func main() {
	loaded, err := ioutil.ReadFile("config.json")
	read := err == nil

	if err != nil {
		var token string
		var delay string

		fmt.Print("Enter your user token: ")
		_, _ = fmt.Scanln(&token)

		fmt.Print("Enter minutes to wait before deleting: ")
		_, _ = fmt.Scanln(&delay)

		delayInt, err := strconv.Atoi(delay)
		if err != nil {
			fmt.Println("Unable to parse input, setting to default 30 minutes.")
			delayInt = 30
		} else {
			fmt.Println("Deleting every message " + delay + " minute(s) after posting.")
		}

		Config.Token = token
		Config.Delay = delayInt
	} else {
		_ = json.Unmarshal(loaded, &Config)
	}

	session, err := discordgo.New(Config.Token)
	if err != nil {
		fmt.Println("Error authenticating (delete config.json if persists):", err)
		exit()
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			go func(s *discordgo.Session, m *discordgo.MessageCreate) {
				time.Sleep(time.Duration(Config.Delay) * time.Minute)
				if s.ChannelMessageDelete(m.ChannelID, m.ID) == nil {
					fmt.Println("Deleted the message", m.ID, "in channel", m.ChannelID+".")
				}
			}(s, m)
		}
	})

	if err = session.Open(); err != nil {
		fmt.Println("Error opening session:", err)
		exit()
	}

	if !read {
		bytes, _ := json.Marshal(Config)
		fmt.Println("The token entered is valid & will be saved to `config.json` for the future.")
		_ = ioutil.WriteFile("config.json", bytes, os.ModePerm)
	} else {
		fmt.Println("The token is has been loaded from `config.json` and authenticated successfully.")
	}

	_, _ = fmt.Scanln()
}

func exit() {
	_, _ = fmt.Scanln()
	os.Exit(0)
}
