// Package main provides the entry point for the PhD Dissertation Planner application.
// This application generates LaTeX-based academic planners and calendars from CSV data.
//
// The application supports:
//   - CSV-based task management
//   - Monthly calendar generation
//   - LaTeX document output
//   - Configurable layouts and themes
//
// Usage:
//   plannergen [flags]
//
// Flags:
//   -config: Configuration file path (default: internal/common/base.yaml)
//   -preview: Render only one page per unique module
//   -outdir: Output directory for generated files
package main

import (
	"fmt"
	"os"

	"phd-dissertation-planner/internal/application"
)

// main is the entry point of the application.
// It initializes the CLI application and handles any fatal errors.
func main() {
	app := application.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Fatal error: %v\n", err)
		os.Exit(1)
	}
}
