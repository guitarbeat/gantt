package main

import (
	"log"
	"os"

	"latex-yearly-planner/internal/app"
)

// * Main entry point for the plannergen CLI application
// * Creates and runs the CLI application using the internal app package
func main() {
	// * Get the CLI application instance from the internal app package
	cliApp := app.New()
	
	// * Run the CLI application with command line arguments
	// * If an error occurs, log it and exit with status code 1
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}