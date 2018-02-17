package main

import (
	"medesk/cmd/assembly"
	"medesk/cmd/settings"

	"github.com/joho/godotenv"

	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := assembly.NewApp()
	port := settings.MustGetSetting("PORT")
	app.Run(port)
}
