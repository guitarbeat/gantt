package main

import (
	"os"

	"phd-dissertation-planner/internal/application"
)

func main() {
	app := application.New()
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
