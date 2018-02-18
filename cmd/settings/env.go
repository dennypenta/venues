package settings

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func Load() {
	if err := godotenv.Load(os.ExpandEnv("$GOPATH/src/venues/.env")); err != nil {
		log.Warn("Error loading .env file")
	}
}

func MustGetSetting(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(fmt.Sprintf("Error retreiving \"%s\"", key))
	}

	return value
}
