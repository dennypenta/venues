package settings

import (
	"os"
	"github.com/labstack/gommon/log"
	"fmt"
)

func MustGetSetting(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(fmt.Sprintf("Error retreiving \"%s\"", key))
	}

	return value
}
