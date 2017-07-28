# logger

## Installing:
```bash
go get -u github.com/Hu13er/logger
```

## Example:
```go
package main

import (
	"fmt"
	"os"

	"github.com/Hu13er/logger"
)

var log *logger.LogFmt

func init() {
	prefix := map[int]string{
		logger.DebugLevel: "(d) ",
		logger.InfoLevel:  "(!) ",
		logger.WarnLevel:  "[*] ",
		logger.ErrorLevel: "[X] ",
		logger.PanicLevel: "[#] ",
	}
	stdlog := logger.New(os.Stderr)
	stdlog.SetPrefix(prefix)

	var (
		token    string // Your bot token
		chatID   int64  // Your ChatID
		chanSize int    // Size of message queue
	)
	telegramStream, err := logger.NewTelegSteam(token, chatID, chanSize)
	if err != nil {
		panic(err)
	}

	telegramLog := logger.New(telegramStream)
	telegramLog.SetPrefix(prefix)

	inti := logger.NewLoggerIntegrate(stdlog, telegramLog)
	log = &logger.LogFmt{Logger: inti}
}

func foo() {
	log.WithHeaderln("In foo:").Infoln("I am a info in foo")
}

func main() {
	log.WithHeaderln("In main:").Warn("I am a warn in main")
	foo()
	fmt.Scanln()
}
```

Coded with <3 for Mansur :D

