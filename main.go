package main

import (
	"./cmd/assembly"
)

func main() {
	app := assembly.NewApp()

	app.Run("8000")
}
