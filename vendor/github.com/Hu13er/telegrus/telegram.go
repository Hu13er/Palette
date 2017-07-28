package telegrus

import (
	"fmt"
	"net/http"
)

const sendMessageRequest = "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

type telegramBot struct {
	botToken string
	chatID   int64
	queue    chan string
	cancel   chan struct{}
}

func newTelegramBot(botToken string, chatID int64, chanSize int) *telegramBot {
	bot := &telegramBot{
		botToken: botToken,
		chatID:   chatID,
		queue:    make(chan string, chanSize),
		cancel:   make(chan struct{}),
	}

	go bot.flush()
	return bot
}

func (tb *telegramBot) SendMsg(text string) {
	tb.queue <- text
}

func (tb *telegramBot) Cancel() {
	tb.cancel <- struct{}{}
}

func (tb *telegramBot) flush() {
	for {
		select {
		case txt := <-tb.queue:
			query := fmt.Sprintf(sendMessageRequest, tb.botToken, tb.chatID, txt)
			http.Get(query)
		case <-tb.cancel:
			break
		}
	}
}
