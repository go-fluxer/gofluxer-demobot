package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"github.com/go-fluxer/gofluxer"
)

type Config struct {
	Token  string `json:"FLUXERBOTTOKEN"`
	Prefix string `json:"FLUXERBOTPREFIX"`
}
func LoadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func main() {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Something went wrong while loading config.json file: %v\n", err)
		return
	}

	bot := gofluxer.NewBot(cfg.Token, cfg.Prefix)

	bot.OnMessage(func(m *gofluxer.Message) {
		fmt.Printf("[%s]: %s\n", m.Author.Username, m.Content)
		// ON MESSAGE EVENT
		if strings.ToLower(m.Content) == "ping" {
			bot.SendMessage(m.ChannelID, "Pong!")
		}
	})

	// MAIN COMMANDS

	bot.AddCommand("help", func(m *gofluxer.Message, args []string) {
		bot.SendEmbed(m.ChannelID, map[string]interface{}{
			"title":       "Bot",
			"description": "help\nping\nsay\nesay\ntest-sowner\ntest-nsfw\ntest-dev",
			"color":       0x00FFFF,
		})
	})

	bot.AddCommand("ping", func(m *gofluxer.Message, args []string) {
		bot.SendMessage(m.ChannelID, "Pong! 🏓")
	})

	bot.AddCommand("say", func(m *gofluxer.Message, args []string) {
		if len(args) == 0 {
			bot.SendMessage(m.ChannelID, "What do you want me to say?")
			return
		}
		bot.SendMessage(m.ChannelID, strings.Join(args, " "))
	})

	bot.AddCommand("esay", func(m *gofluxer.Message, args []string) {
		fullInput := strings.Join(args, " ")
		parts := strings.Split(fullInput, "|")
		if len(parts) < 2 {
			bot.SendMessage(m.ChannelID, "Usage: !essay Title | Description")
			return
		}

		bot.SendEmbed(m.ChannelID, map[string]interface{}{
			"title":       strings.TrimSpace(parts[0]),
			"description": strings.TrimSpace(parts[1]),
			"color":       0x00FFFF,
		})
	})

	// TEST COMMANDS

	bot.AddCommand("test-sowner", func(m *gofluxer.Message, args []string) {
		if !bot.IsOwner(m) {
			bot.SendMessage(m.ChannelID, "You need to have **SERVER OWNER** permission to use this command.")
			return
		}
		bot.SendMessage(m.ChannelID, "Hello server owner")
	})

	bot.AddCommand("test-nsfw", func(m *gofluxer.Message, args []string) {
		if !bot.IsNSFW(m.ChannelID) {
			bot.SendMessage(m.ChannelID, "This command can only be used in NSFW marked channels for safety reason")
			return
		}

		bot.SendMessage(m.ChannelID, "This is a NSFW marked channel.")
	})

	bot.AddCommand("test-dev", func(m *gofluxer.Message, args []string) {
		const devID = "1473453696008249383" // Your Fluxer user ID here....
		if m.Author.ID == devID {
			bot.SendMessage(m.ChannelID, "Hello, my bot dev")
		} else {
			bot.SendMessage(m.ChannelID, "Only Bot Staff members can use this command")
		}
	})





	fmt.Println("[Fluxer.app] Bot is Ready")
	if err := bot.Run(); err != nil {
		fmt.Printf("[Fluxer.app] Bot has stopped: %v\n", err)
	}
}