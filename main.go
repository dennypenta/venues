package main

import (
	"venues/cmd/assembly"
	"venues/cmd/settings"
)

func main() {
	settings.Load()

	app := assembly.NewApp()
	port := settings.MustGetSetting("PORT")
	app.Run(port)
}
