package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"vault/configs"
	"vault/internal/bot"
	"vault/internal/db"
	"vault/internal/db/queries"
	"vault/internal/vault"
)

func main() {
	config, err := configs.LoadConfig("./configs/")
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	db, err := db.New(config.PostgresDSN)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap error: %s", err)
	}

	vault, err := vault.New(db, config.BotEncryptionKey, logger)
	if err != nil {
		log.Fatalf("vault error: %s", err)
	}

	bot, err := bot.New(config.BotToken, config.BotVisibilityPeriod, vault, logger)
	if err != nil {
		log.Fatalf("bot error: %s", err)
	}

	log.Println("Starting up vault bot...")

	go bot.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down vault bot...")

	bot.Stop()
	if err := queries.Close(); err != nil {
		log.Fatalf("queries close error: %s", err)
	}
}
