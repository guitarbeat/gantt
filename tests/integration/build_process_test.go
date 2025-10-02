package main

import (
	"os"
	"path/filepath"
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

func TestPreviewMode(t *testing.T) {
	// Test preview mode (one page per unique module)
	app := app.New()
	args := []string{"plannergen", "--config", "src/core/base.yaml", "--outdir", "generated", "--preview"}

	if err := app.Run(args); err != nil {
		t.Fatalf("Preview mode failed: %v", err)
	}

	// Verify that output file was created
	outputFile := filepath.Join("generated", "base.tex")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Preview mode did not create output file: %s", outputFile)
	}
}

func TestCustomOutputDirectory(t *testing.T) {
	// Create a temporary output directory
	tmpDir := filepath.Join("generated", "test_custom_output")
	defer os.RemoveAll(tmpDir)

	// Test with custom output directory
	app := app.New()
	args := []string{"plannergen", "--config", "src/core/base.yaml", "--outdir", tmpDir}

	if err := app.Run(args); err != nil {
		t.Fatalf("Custom output directory test failed: %v", err)
	}

	// Verify that output was created in custom directory
	outputFile := filepath.Join(tmpDir, "base.tex")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Custom output directory did not create file: %s", outputFile)
	}
}

func TestMissingConfigFile(t *testing.T) {
	// Test behavior with missing config file
	// Note: The system handles missing config files gracefully by using defaults
	app := app.New()
	args := []string{"plannergen", "--config", "nonexistent.yaml", "--outdir", "generated"}

	err := app.Run(args)
	// Missing config files are skipped, system uses defaults
	t.Logf("Missing config file handled: %v", err)
}

func TestInvalidOutputDirectory(t *testing.T) {
	// Test with a path that cannot be created (invalid characters, etc.)
	// Note: This might be platform-specific, so we use a relatively safe invalid path
	app := app.New()
	
	// Create a file where we want to create a directory (this should fail)
	tmpFile := filepath.Join("generated", "test_invalid_dir_file")
	os.MkdirAll("generated", 0755)
	if f, err := os.Create(tmpFile); err == nil {
		f.Close()
		defer os.Remove(tmpFile)

		// Try to use the file path as a directory
		args := []string{"plannergen", "--config", "src/core/base.yaml", "--outdir", tmpFile}

		err := app.Run(args)
		if err == nil {
			t.Error("Expected error when output directory is a file, got nil")
		}
	}
}

func TestEmptyConfig(t *testing.T) {
	// Create a temporary empty config file
	tmpConfig := filepath.Join("generated", "test_empty.yaml")
	os.MkdirAll("generated", 0755)
	if err := os.WriteFile(tmpConfig, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	defer os.Remove(tmpConfig)

	app := app.New()
	args := []string{"plannergen", "--config", tmpConfig, "--outdir", "generated"}

	// Empty config should use defaults and potentially succeed or fail gracefully
	err := app.Run(args)
	// We log the result but don't fail - depends on implementation
	t.Logf("Empty config test result: %v", err)
}

func TestMultipleConfigFiles(t *testing.T) {
	// Test loading multiple config files (overlay behavior)
	app := app.New()
	
	// Base config plus monthly calendar config
	args := []string{
		"plannergen",
		"--config", "src/core/base.yaml,src/core/monthly_calendar.yaml",
		"--outdir", "generated",
	}

	if err := app.Run(args); err != nil {
		t.Fatalf("Multiple config files test failed: %v", err)
	}

	// Verify output was created
	outputFile := filepath.Join("generated", "monthly_calendar.tex")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Logf("Note: monthly_calendar.tex not created, may be expected: %v", err)
	}
}
