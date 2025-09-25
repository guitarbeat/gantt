package main

import (
	"fmt"
	"os"

	"phd-dissertation-planner/src/app"
)

func main() {
	app := app.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Fatal error: %v\n", err)
		os.Exit(1)
	}
}
