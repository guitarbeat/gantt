package config

import "testing"

func TestDefaultOutputDir(t *testing.T) {
	cfg, err := NewConfig() // no files, should fall back to defaults
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}
	if cfg.OutputDir != "build" {
		t.Fatalf("OutputDir should default to 'build', got %q", cfg.OutputDir)
	}
}
