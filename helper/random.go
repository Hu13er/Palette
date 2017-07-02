package helper

import (
	"math/rand"
	"time"
)

type Charset string

const (
	DefaultCharset = Charset("abcdefghijklmnopqrstuvwxyz1234567890")
	NumricCharset  = Charset("1234567890")
)

func init() {
	rand.Seed(time.Now().Unix())
}

func (cs Charset) RandomStr(size int) string {
	var (
		charset = string(cs)
		outp    = ""
	)

	for i := 0; i < size; i++ {
		outp += string(charset[rand.Intn(len(charset))])
	}
	return outp
}
