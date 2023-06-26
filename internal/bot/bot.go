package bot

import (
	"fmt"
	"time"

	"vault/internal/vault"

	"go.uber.org/zap"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents a Telegram bot.
type Bot struct {
	token  string
	vault  *vault.Vault
	logger *zap.Logger
	*tg.BotAPI
	stopHiding   func()
	toHide       chan Message
	hideInterval int64
}

type messages struct {
	English    string
	Portuguese string
}

type keyboards struct {
	English    tg.InlineKeyboardMarkup
	Portuguese tg.InlineKeyboardMarkup
}

// New creates a new bot.
func New(token string, visibilityPeriod time.Duration, vault *vault.Vault, logger *zap.Logger) (*Bot, error) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &Bot{
		token:        token,
		vault:        vault,
		BotAPI:       bot,
		logger:       logger,
		hideInterval: 60,
	}, nil
}

// Start starts the bot.
func (bot *Bot) Start() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	bot.toHide, bot.stopHiding = bot.Watch()

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.CallbackQuery != nil {
			bot.handleCallbackQuery(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			bot.handleCommand(update.Message)
			continue
		}

		bot.handleMessage()
	}

}

// Stop stops the bot.
func (bot *Bot) Stop() {
	bot.StopReceivingUpdates()
	bot.stopHiding()
}
