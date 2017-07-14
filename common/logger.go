package common

import (
	"github.com/Hu13er/telegrus"
	"github.com/sirupsen/logrus"
)

func init() {
	if level, err := logrus.ParseLevel(ConfigString("LOG_LEVEL")); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var (
		chatID = ConfigInt64("TELEGRAM_CHATID")
		token  = ConfigString("TELEGRAM_TOKEN")
	)

	if token == "" {
		logrus.Warnln("TELEGRAM_TOKEN not presented. Telegram logger DISABLED.")
		return
	}

	if chatID == int64(0) {
		logrus.Fatalln("TELEGRAM_CHATID not presented.")
	}

	telegramHooker := telegrus.NewHooker(token, chatID).SetLevel(logrus.GetLevel())

	if username := ConfigString("TELEGRAM_MENTION"); username != "" {
		telegramHooker.MentionOn(logrus.WarnLevel, username)
	}

	logrus.AddHook(telegramHooker)
}
