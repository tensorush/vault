package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"vault-bot/internal/database"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCommand handles commands.
func (b *Bot) handleCommand(msg *tg.Message) {
	switch msg.Command() {
	case start:
		b.handleStart(msg)
	case set:
		b.handleSet(msg)
	case get:
		b.handleGet(msg)
	case del:
		b.handleDel(msg)
	}
}

// handleMessage handles messages.
func (b *Bot) handleMessage() {

}

// handleMessageLang handles language messages.
func (b *Bot) handleMessageLang(msg string, chatID int64) string {
	lang := b.vault.GetLang(chatID)
	switch lang {
	case "en":
		return allMessages[msg].English
	default:
		return allMessages[msg].Portuguese
	}
}

// handleKeyboardLang handles keyboards languages.
func (b *Bot) handleKeyboardLang(keyboard string, chatID int64) tg.InlineKeyboardMarkup {
	lang := b.vault.GetLang(chatID)
	switch lang {
	case "en":
		return allKeyboards[keyboard].English
	default:
		return allKeyboards[keyboard].Portuguese
	}
}

// handleStart handles start command.
func (b *Bot) handleStart(msg *tg.Message) {
	msgConfig := tg.NewMessage(msg.Chat.ID, fmt.Sprintf(b.handleMessageLang(start, msg.Chat.ID), b.hideInterval))
	msgConfig.ReplyMarkup = b.handleKeyboardLang(startKeyboard, msg.Chat.ID)

	_, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}

// handleSet handles set command.
func (b *Bot) handleSet(msg *tg.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tg.NewMessage(msg.Chat.ID, b.handleMessageLang(set, msg.Chat.ID))
	if len(split) != 4 {
		msgConfig = tg.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
		m, err := b.Send(msgConfig)
		if err != nil {
			log.Println("send error: ", err)
		} else {
			b.toHide <- Message{
				chatID:    msg.Chat.ID,
				id:        msg.MessageID,
				createdAt: time.Now(),
			}

			b.toHide <- Message{
				chatID:    m.Chat.ID,
				id:        m.MessageID,
				createdAt: time.Now(),
			}
		}
		return
	}

	err := b.vault.Save(msg.Chat.ID, split[1], split[2], split[3])
	if err != nil {
		msgConfig.Text = b.handleMessageLang(setErr, msg.Chat.ID)
		log.Printf("save error: %v\n", err)
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- Message{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- Message{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

// handleGet handles get command.
func (b *Bot) handleGet(msg *tg.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tg.NewMessage(msg.Chat.ID, b.handleMessageLang(get, msg.Chat.ID))
	if len(split) != 2 {
		msgConfig = tg.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
		m, err := b.Send(msgConfig)
		if err != nil {
			log.Println("send error: ", err)
		} else {
			b.toHide <- Message{
				chatID:    msg.Chat.ID,
				id:        msg.MessageID,
				createdAt: time.Now(),
			}

			b.toHide <- Message{
				chatID:    m.Chat.ID,
				id:        m.MessageID,
				createdAt: time.Now(),
			}
		}
		return
	}
	service := split[1]

	cred, err := b.vault.Get(msg.Chat.ID, service)
	if err != nil {
		if errors.Is(err, database.ErrServiceNotFound) {
			msgConfig.Text = b.handleMessageLang(serviceNotFoundErr, msg.Chat.ID)
		} else {
			msgConfig.Text = b.handleMessageLang(getErr, msg.Chat.ID)
		}
		log.Printf("get error: %v\n", err)
	} else {
		msgConfig.ReplyMarkup = b.handleKeyboardLang(hideKeyboard, msg.Chat.ID)
		msgConfig.Text = fmt.Sprintf(b.handleMessageLang(get, msg.Chat.ID), service, cred.Login, cred.Password)
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- Message{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- Message{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

// handleDel handles delete command.
func (b *Bot) handleDel(msg *tg.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tg.NewMessage(msg.Chat.ID, b.handleMessageLang(del, msg.Chat.ID))
	if len(split) != 2 {
		msgConfig = tg.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
		m, err := b.Send(msgConfig)
		if err != nil {
			log.Println("send error: ", err)
		} else {
			b.toHide <- Message{
				chatID:    msg.Chat.ID,
				id:        msg.MessageID,
				createdAt: time.Now(),
			}

			b.toHide <- Message{
				chatID:    m.Chat.ID,
				id:        m.MessageID,
				createdAt: time.Now(),
			}
		}
		return
	}
	service := split[1]

	err := b.vault.Delete(msg.Chat.ID, service)
	if err != nil {
		if errors.Is(err, database.ErrServiceNotFound) {
			msgConfig.Text = b.handleMessageLang(serviceNotFoundErr, msg.Chat.ID)
		} else {
			msgConfig.Text = b.handleMessageLang(delErr, msg.Chat.ID)
			log.Printf("del error: %v\n", err)
		}
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- Message{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- Message{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

// handleCallbackQuery handles callback queries from user.
func (b *Bot) handleCallbackQuery(query *tg.CallbackQuery) {
	split := strings.Split(query.Data, "::")
	if len(split) == 0 {
		return
	}

	defer b.logger.Sync()

	text := split[0]

	switch text {
	case hide:
		msg := tg.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
		}

		msg = tg.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID-1)
		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
		}
	case changeLang:
		msg := tg.NewEditMessageTextAndMarkup(
			query.Message.Chat.ID,
			query.Message.MessageID,
			"Choose a new language 🌎",
			b.handleKeyboardLang(setLangKeyboard, query.Message.Chat.ID),
		)

		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("send error: %v", err.Error()))
		}
	case change:
		if len(split) == 1 {
			return
		}

		b.vault.SetLang(query.Message.Chat.ID, split[1])

		msg := tg.NewEditMessageTextAndMarkup(
			query.Message.Chat.ID, query.Message.MessageID,
			fmt.Sprintf(b.handleMessageLang(start, query.Message.Chat.ID), b.hideInterval),
			b.handleKeyboardLang(startKeyboard, query.Message.Chat.ID),
		)

		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("send error: %v", err.Error()))
		}

	}
}
