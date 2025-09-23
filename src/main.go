package main

import (
	"os"

	"phd-dissertation-planner/internal/core"
)

func main() {
	app := core.New()
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
