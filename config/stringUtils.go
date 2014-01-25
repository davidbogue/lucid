package config

import (
	"strings"
)

func startsWith(main string, test string) bool {
	return strings.Index(main, test) == 0
}

func splitOnEquals(value string) (string, string) {
	index := strings.Index(value, "=")
	key := value[0:index]
	property := value[index+1:]
	return key, property
}
