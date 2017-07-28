package logger

import (
	"sync"

	telegram "gopkg.in/telegram-bot-api.v4"
)

type TelegSteam struct {
	bot    *telegram.BotAPI
	chatID int64
	mutex  sync.Mutex
	queue  chan string
	cancel chan struct{}
}

func NewTelegSteam(token string, chatID int64, chanSize int) (*TelegSteam, error) {
	bot, err := telegram.NewBotAPI(token)
	tele := &TelegSteam{bot: bot, queue: make(chan string, chanSize), chatID: chatID, mutex: sync.Mutex{}, cancel: make(chan struct{})}
	go tele.flusher()
	return tele, err
}

func (t *TelegSteam) Write(buf []byte) (n int, err error) {

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.queue <- string(buf)

	return len(buf), nil
}

func (t *TelegSteam) Cancel() {
	t.cancel <- struct{}{}
	close(t.queue)
}

func (t *TelegSteam) flusher() {
	for {
		select {
		case <-t.cancel:
			return
		case item := <-t.queue:
			msg := telegram.NewMessage(t.chatID, item)
			t.bot.Send(msg)
		}
	}
}
