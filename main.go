package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/go-fluxer/gofluxer"
)

type Config struct {
	Token  string `json:"FLUXERBOTTOKEN"`
	Prefix string `json:"FLUXERBOTPREFIX"`
	WebhookId string `json:"FLUXERWEBHOOKID"`
	WebhookToken string `json:"FLUXERWEBHOOKTOKEN"`
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

type GuildSettings struct {
	LogChannelID string `json:"log_channel_id"`
}
type Database struct {
	mu       sync.RWMutex
	Filename string
	Guilds   map[string]GuildSettings
}
var db *Database
func InitDB(filename string) *Database {
	d := &Database{
		Filename: filename,
		Guilds:   make(map[string]GuildSettings),
	}
	file, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(file, &d.Guilds)
	}
	return d
}
func (d *Database) Save() error {
	data, err := json.MarshalIndent(d.Guilds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(d.Filename, data, 0644)
}

func main() {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Something went wrong while loading config.json file: %v\n", err)
		return
	}
	db = InitDB("guilds.json")

	bot := gofluxer.NewBot(cfg.Token, cfg.Prefix)
	// bot.NewBotInstance("https://example.com/api/v1", "wss://example.com/gateway/?v=1")
	// The bot.NewBotInstance() function allows you to use a 3rd party Fluxer instance instead of using the main one.
	bot.NewBotConfig(true, 500)



	bot.OnReady(func() {
		fmt.Println("Demo [Fluxer.app] Bot is Ready")
	})

	bot.OnUserJoin(func(u *gofluxer.UserJoinPayload) {
		db.mu.RLock()
		settings, exists := db.Guilds[u.GuildID]
		db.mu.RUnlock()
		if exists && settings.LogChannelID != "" {
			welcomeMsg := fmt.Sprintf("**MEMBER JOINED:** Welcome %s to %s!", u.User.Username, u.GuildID)
			bot.SendMessage(settings.LogChannelID, welcomeMsg)
		}
	})

	bot.OnUserLeave(func(l *gofluxer.UserLeavePayload) {
		db.mu.RLock()
		settings, exists := db.Guilds[l.GuildID]
		db.mu.RUnlock()
		if exists && settings.LogChannelID != "" {
			welcomeMsg := fmt.Sprintf("**MEMBER LEAVE:** %s (ID: %s) has left %s.\n", l.User.Username, l.UserID, l.GuildName)
			bot.SendMessage(settings.LogChannelID, welcomeMsg)
		}
	})

	bot.OnMessageDelete(func(d *gofluxer.MessageDeletePayload) {
		db.mu.RLock()
		settings, exists := db.Guilds[d.GuildID]
		db.mu.RUnlock()
		if exists && settings.LogChannelID != "" {
			deleteMsg := fmt.Sprintf("**MESSAGE DELETED:** The following message ID %s was deleted from <#%s>.", d.MessageID, d.ChannelID)
			bot.SendMessage(settings.LogChannelID, deleteMsg)
		}
	})



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
			"description": "help\nping\nsay\nesay\nsetlogchannel\ndisablelogchannel\nannounce\ntest-reply\ntest-sowner\ntest-nsfw",
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
	
	// CONFIG COMMANDS

	bot.AddCommand("setlogchannel", func(m *gofluxer.Message, args []string) {
		if !bot.IsOwner(m) {
			bot.SendMessage(m.ChannelID, "You need to have the **SERVER OWNER** permission to use this command.")
			return
		}
		if len(args) == 0 {
			bot.SendMessage(m.ChannelID, "Arguments required. Please provide a valid target Channel ID.")
			return
		}

		targetChannelID := args[0]
		db.mu.Lock()
		db.Guilds[m.GuildID] = GuildSettings{LogChannelID: targetChannelID}
		err := db.Save()
		db.mu.Unlock()
		if err != nil {
			bot.SendMessage(m.ChannelID, "An error occured while saving the configuration file.")
			return
		}
		bot.SendMessage(m.ChannelID, fmt.Sprintf("Logging channel successfully set to <#%s>!", targetChannelID))
	})

	bot.AddCommand("disablelogchannel", func(m *gofluxer.Message, args []string) {
		if !bot.IsOwner(m) {
			bot.SendMessage(m.ChannelID, "You need to have the **SERVER OWNER** permission to use this command.")
			return
		}

		db.mu.Lock()
		_, exists := db.Guilds[m.GuildID]
		if exists {
			delete(db.Guilds, m.GuildID)
			err = db.Save()
		}
		db.mu.Unlock()
		if !exists {
			bot.SendMessage(m.ChannelID, "Logging is already disabled or wasn't set up for this server.")
			return
		}
		if err != nil {
			bot.SendMessage(m.ChannelID, "An error occured while saving the configuration file.")
			return
		}
		bot.SendMessage(m.ChannelID, "Logging has been disabled.")
	})

	bot.AddCommand("announce", func(m *gofluxer.Message, args []string) {
		const devID = "1473453696008249383" // Your Fluxer user ID here....
		if m.Author.ID == devID {
			if len(args) == 0 {
				bot.SendMessage(m.ChannelID, "What do you want me to say?")
				return
			}
			wh := gofluxer.NewWebhookClient(cfg.WebhookId, cfg.WebhookToken)
			wh.Execute(strings.Join(args, " "))
			bot.SendMessage(m.ChannelID, "Message sent to webhook")
		} else {
			bot.SendMessage(m.ChannelID, "Only Bot Staff members can use this command")
		}
	})

	// TEST COMMANDS
	
	bot.AddCommand("test-reply", func(m *gofluxer.Message, args []string) {
		bot.ReplyMessage(m, "Hello World!")
	})

	bot.AddCommand("test-sowner", func(m *gofluxer.Message, args []string) {
		if !bot.IsOwner(m) {
			bot.SendMessage(m.ChannelID, "You need to have the **SERVER OWNER** permission to use this command.")
			return
		}
		bot.SendMessage(m.ChannelID, "Hello server owner")
	})

	bot.AddCommand("test-nsfw", func(m *gofluxer.Message, args []string) {
		if !bot.IsNSFW(m.ChannelID) {
			bot.SendMessage(m.ChannelID, "This command can only be used in NSFW marked channels for safety reasons")
			return
		}

		bot.SendMessage(m.ChannelID, "This is a NSFW marked channel.")
	})





	bot.Run()
}