package config

import (
	"bufio"
	"os"
	"strings"
)

var propertyMap map[string]string
var propertyFile = string("./config/config.properties")

func init() {
	propertyMap = make(map[string]string)
	file, err := os.Open(propertyFile)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			propLine := scanner.Text()
			if !startsWith(propLine, "#") {
				key, value := splitOnEquals(propLine)
				propertyMap[key] = value
			}
		}
	}
}

func Value(key string) string {
	return propertyMap[key]
}

func startsWith(main string, test string) bool {
	return strings.Index(main, test) == 0
}

func splitOnEquals(value string) (string, string) {
	index := strings.Index(value, "=")
	key := value[0:index]
	property := value[index+1:]
	return key, property
}
