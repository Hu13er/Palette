# Telegram Hooks for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Install

```shell
$ go get github.com/Hu13er/telegrus
```

## Usage

```go
package main

import (
    "github.com/sirupsen/logrus"
    "github.com/Hu13er/telegrus"
)

func main() {
    log := logrus.New()
    
    var (
        botToken = ""      // Your Bot token
        chatID = int64(0) // Your chatID
    )

    logrus.AddHook(
        telegrus.NewHooker(botToken, chatID).
            MentionOn(logrus.WarnLevel,
                "Hu13er", "foobar").
            MentionOn(logrus.ErrorLevel,
                "Huberrr").
            SetLevel(logrus.InfoLevel),
    )

    logrus.Debugln("This is a DEBUG")
    logrus.Infoln("This is an INFO")
    logrus.Warnln("This is a WARN")
    logrus.Errorln("This is an ERROR")
    fmt.Scanln()
}
```

## License
*MIT*
