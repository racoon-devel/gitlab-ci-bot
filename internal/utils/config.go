package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
)

type BotConfig struct {
	Token string
	Chat  int64
}

type ServerConfig struct {
	Endpoint string
}

type Config struct {
	TelegramBot BotConfig
	Server      ServerConfig
}

func (c *Config) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("open file failed: %+w", err)
	}

	if err = toml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parse config file content failed: %+w", err)
	}

	// override some values from ENV variables
	if botToken := os.Getenv("TELEGRAM_BOT_TOKEN"); botToken != "" {
		c.TelegramBot.Token = botToken
	}
	if botChatID := os.Getenv("TELEGRAM_BOT_CHAT"); botChatID != "" {
		c.TelegramBot.Chat, err = strconv.ParseInt(botChatID, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot override bot chat ID via env variable: %w", err)
		}
	}

	return nil
}
