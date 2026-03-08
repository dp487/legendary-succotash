package main

import (
	"log"

	"github.com/dp487/legendary-succotash/app"
)

func main() {
	err := app.SetupAndRunApp()
	if err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}