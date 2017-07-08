package common

import (
	"io"
	"log"
	"os"
)

var Stdlog io.Writer

func init() {
	var (
		chatID   = ConfigInt64("TELEGRAM_CHATID")
		token    = ConfigString("TELEGRAM_TOKEN")
		chanSize = 32
	)

	if token == "" {
		log.Fatalln("TELEGRAM_TOKEN not presented.")
	}

	if chatID == int64(0) {
		log.Fatalln("TELEGRAM_CHATID not presented.")
	}

	telegramSteam, err := newTelegSteam(token, chatID, chanSize)
	if err != nil {
		log.Fatalln("Can not connect to Telegram bot")
	}

	stdErr := os.Stderr

	Stdlog = intigrate(telegramSteam, stdErr)

	log.SetOutput(Stdlog)
	log.SetFlags(log.Llongfile)
}

type intigrater []io.Writer

func intigrate(args ...io.Writer) intigrater {
	return intigrater(args)
}

func (i intigrater) Write(p []byte) (n int, err error) {
	for _, v := range i {
		n, err = v.Write(p)
		if err != nil {
			return
		}
	}
	return
}
