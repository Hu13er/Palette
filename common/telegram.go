package common

import (
	"sync"

	telegram "gopkg.in/telegram-bot-api.v4"
)

type telegSteam struct {
	bot    *telegram.BotAPI
	chatID int64
	mutex  sync.Mutex
	queue  chan string
	cancel chan struct{}
}

func newTelegSteam(token string, chatID int64, chanSize int) (*telegSteam, error) {
	bot, err := telegram.NewBotAPI(token)
	tele := &telegSteam{bot: bot, queue: make(chan string, chanSize), chatID: chatID, mutex: sync.Mutex{}, cancel: make(chan struct{})}
	go tele.flusher()
	return tele, err
}

func (t *telegSteam) Write(buf []byte) (n int, err error) {

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.queue <- string(buf)

	return len(buf), nil
}

func (t *telegSteam) Cancel() {
	t.cancel <- struct{}{}
	close(t.queue)
}

func (t *telegSteam) flusher() {
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
