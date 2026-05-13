package env

import (
	"fmt"

	"github.com/joho/godotenv"
)

var envMap map[string]string

func Load() error {
	var err error

	envMap, err = godotenv.Read()
	if err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}

	return nil
}

func GetEnv(key string, defaultVal string) string {
	if value, exists := envMap[key]; exists {
		return value
	}
	return defaultVal
}
