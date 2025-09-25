package main

import (
	"fmt"
	"os"

	"phd-dissertation-planner/internal/application"
)

func main() {
	app := application.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Fatal error: %v\n", err)
		os.Exit(1)
	}
}
