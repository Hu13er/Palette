package common

import (
	"os"
	"strings"
)

const prefix = "PALETTE_"

func ConfigString(key string) string {
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
