package config

import (
	"bufio"
	"os"
)

var propertyMap map[string]string
var propertyFile = string("config.properties")

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
