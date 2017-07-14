package common

import (
	"github.com/Hu13er/telegrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	if level, err := log.ParseLevel(ConfigString("LOG_LEVEL")); err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	var (
		chatID = ConfigInt64("TELEGRAM_CHATID")
		token  = ConfigString("TELEGRAM_TOKEN")
	)

	if token == "" {
		log.Warnln("TELEGRAM_TOKEN not presented. Telegram logger DISABLED.")
		return
	}

	if chatID == int64(0) {
		log.Fatalln("TELEGRAM_CHATID not presented.")
	}

	telegramHooker := telegrus.NewHooker(token, chatID).SetLevel(log.GetLevel())

	if username := ConfigString("TELEGRAM_MENTION"); username != "" {
		telegramHooker.MentionOn(log.WarnLevel, username)
	}

	log.AddHook(telegramHooker)
}
