package bot

import (
	"fmt"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Message contains information about message.
type Message struct {
	chatID    int64
	id        int
	createdAt time.Time
}

// Watch watches messages and deletes them after hideInterval.
func (b *Bot) Watch() (chan Message, func()) {
	messagesCh := make(chan Message, 10000)
	cancelCh := make(chan bool)

	go func() {
		for {
			select {
			case <-cancelCh:
				close(messagesCh)
				return
			case msg := <-messagesCh:
				for time.Now().Unix()-b.hideInterval <= msg.createdAt.Unix() {
				}

				msgDelConfig := tg.NewDeleteMessage(msg.chatID, msg.id)
				if _, err := b.Send(msgDelConfig); err != nil {
					b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
				}
			}
		}
	}()

	return messagesCh, func() {
		close(cancelCh)
	}
}
