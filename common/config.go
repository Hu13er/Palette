package common

import (
	"os"
	"strconv"
	"strings"
)

const prefix = "PALETTE_"

func ConfigString(key string) string {
	key = strings.ToUpper(key)
	return os.Getenv(prefix + key)
}

func ConfigBool(key string) bool {
	trues := []string{"YES", "Y", "TRUE", "T", "1"}
	value := ConfigString(key)
	value = strings.ToUpper(value)
	for _, v := range trues {
		if value == v {
			return true
		}
	}
	return false
}

func ConfigInt64(key string) int64 {
	outp, _ := strconv.ParseInt(ConfigString(key), 10, 32)
	return outp
}
