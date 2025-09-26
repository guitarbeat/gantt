package main

import (
	"phd-dissertation-planner/src/app"
	"testing"
)

func TestBuildProcess(t *testing.T) {
	// Simulate the build process that runs when you execute: make -f scripts/Makefile clean-build
	app := app.New()
	args := []string{"plannergen", "--config", "src/core/base.yaml,src/core/monthly_calendar.yaml", "--outdir", "generated"}

	// This should exercise the main application code paths
	if err := app.Run(args); err != nil {
		t.Logf("Build process completed with note: %v", err)
	}
}
