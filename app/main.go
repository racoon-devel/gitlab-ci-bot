package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/racoon-devel/gitlab-ci-bot/internal/bot"
	"github.com/racoon-devel/gitlab-ci-bot/internal/utils"
)

const version = "0.1"

func main() {
	log.Printf("gitlab-ci-bot v%s", version)

	configPath := flag.String("config", "/etc/gitlab-ci-bot/config.toml", "path to config file")
	flag.Parse()

	var config utils.Config
	if err := config.Load(*configPath); err != nil {
		log.Fatalf("Load config failed: %s", err)
	}

	log.Printf("Config: %+v", config)

	botInstance, err := bot.New(config.TelegramBot.Token, config.TelegramBot.Chat)
	if err != nil {
		log.Fatalf("Start Telegram bot failed: %s", err)
	}

	handler := utils.WebHookHandler{Consumer: botInstance}

	mux := http.NewServeMux()
	mux.Handle("/webhook", handler)
	if err := http.ListenAndServe(config.Server.Endpoint, mux); err != nil {
		log.Fatalf("Start HTTP server failed: %s", err)
	}
}
