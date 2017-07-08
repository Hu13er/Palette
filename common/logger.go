package common

import (
	"github.com/Hu13er/telegrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	var (
		chatID = ConfigInt64("TELEGRAM_CHATID")
		token  = ConfigString("TELEGRAM_TOKEN")
	)

	if token == "" {
		log.Fatalln("TELEGRAM_TOKEN not presented.")
	}

	if chatID == int64(0) {
		log.Fatalln("TELEGRAM_CHATID not presented.")
	}

	if level, err := log.ParseLevel(ConfigString("LOG_LEVEL")); err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	telegramHooker := telegrus.NewHooker(token, chatID)

	if username := ConfigString("TELEGRAM_MENTION"); username != "" {
		telegramHooker.SetMention(map[log.Level][]string{
			log.WarnLevel:  []string{"Huberrr"},
			log.ErrorLevel: []string{"Huberrr"},
			log.PanicLevel: []string{"Huberrr"},
		})
	}

	log.AddHook(telegramHooker)
}
