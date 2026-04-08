package env

import (
	"log"

	"github.com/joho/godotenv"
)

var envMap map[string]string

func Load() {
	var err error

	envMap, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key string, defaultVal string) string {
	if value, exists := envMap[key]; exists {
		return value
	}
	return defaultVal
}
