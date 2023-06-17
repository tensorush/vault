package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"vault-bot/config"
	"vault-bot/internal/bot"
	"vault-bot/internal/database"
	"vault-bot/internal/database/queries"
	"vault-bot/internal/vault"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	db, err := database.New(cfg.Type, cfg.Data)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap error: %s", err)
	}

	vault, err := vault.New(db, cfg.EncryptionKey, logger)
	if err != nil {
		log.Fatalf("vault error: %s", err)
	}

	bot, err := bot.New(cfg.Token, cfg.ExpirationPeriod, vault, logger)
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
