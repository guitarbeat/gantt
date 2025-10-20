package core_test

import (
	"bytes"
	"os"
	"testing"

	"phd-dissertation-planner/internal/core"
)

func TestNewLogger(t *testing.T) {
	logger := core.NewLogger("[test] ")

	if logger == nil {
		t.Fatal("core.NewLogger() should not return nil")
	}

	if logger.Writer == nil {
		t.Error("core.NewLogger() should initialize writer")
	}

	if logger.Level < 0 {
		t.Error("core.NewLogger() should set valid log level")
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := core.NewDefaultLogger()

	if logger == nil {
		t.Fatal("NewDefaultLogger() should not return nil")
	}
}

func TestIsSilent(t *testing.T) {
	// Save original env
	originalSilent := os.Getenv("PLANNER_SILENT")
	originalLevel := os.Getenv("PLANNER_LOG_LEVEL")
	defer func() {
		os.Setenv("PLANNER_SILENT", originalSilent)
		os.Setenv("PLANNER_LOG_LEVEL", originalLevel)
	}()

	tests := []struct {
		name       string
		silentEnv  string
		levelEnv   string
		wantSilent bool
	}{
		{"no env vars", "", "", false},
		{"PLANNER_SILENT=1", "1", "", true},
		{"PLANNER_LOG_LEVEL=silent", "", "silent", true},
		{"PLANNER_LOG_LEVEL=info", "", "info", false},
		{"PLANNER_LOG_LEVEL=debug", "", "debug", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PLANNER_SILENT", tt.silentEnv)
			os.Setenv("PLANNER_LOG_LEVEL", tt.levelEnv)

			if got := core.IsSilent(); got != tt.wantSilent {
				t.Errorf("IsSilent() = %v, want %v", got, tt.wantSilent)
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	// Create a logger with a custom output buffer for testing
	var buf bytes.Buffer

	// We can't easily test the actual logging since it goes through log.Logger
	// But we can test that the methods don't panic
	logger := core.NewLogger("[test] ")

	// These should not panic
	t.Run("Info", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Info() panicked: %v", r)
			}
		}()
		logger.Info("test message")
	})

	t.Run("Debug", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Debug() panicked: %v", r)
			}
		}()
		logger.Debug("test message")
	})

	t.Run("Error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Error() panicked: %v", r)
			}
		}()
		logger.Error("test message")
	})

	t.Run("Warn", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Warn() panicked: %v", r)
			}
		}()
		logger.Warn("test message")
	})

	t.Run("Printf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Printf() panicked: %v", r)
			}
		}()
		logger.Printf("test message")
	})

	// Test with formatting
	t.Run("formatted messages", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Formatted logging panicked: %v", r)
			}
		}()
		logger.Info("test %s %d", "message", 42)
		logger.Debug("test %s %d", "message", 42)
		logger.Error("test %s %d", "message", 42)
		logger.Warn("test %s %d", "message", 42)
	})

	_ = buf // Keep compiler happy
}

func TestLogLevelDetection(t *testing.T) {
	// Save and restore environment
	originalSilent := os.Getenv("PLANNER_SILENT")
	originalLevel := os.Getenv("PLANNER_LOG_LEVEL")
	defer func() {
		os.Setenv("PLANNER_SILENT", originalSilent)
		os.Setenv("PLANNER_LOG_LEVEL", originalLevel)
	}()

	tests := []struct {
		name      string
		silentEnv string
		levelEnv  string
		wantLevel int
	}{
		{"default", "", "", core.LogLevelInfo},
		{"silent flag", "1", "", core.LogLevelSilent},
		{"explicit silent", "", "silent", core.LogLevelSilent},
		{"explicit info", "", "info", core.LogLevelInfo},
		{"explicit debug", "", "debug", core.LogLevelDebug},
		{"invalid level", "", "invalid", core.LogLevelInfo}, // Should default to info
		{"case insensitive", "", "SILENT", core.LogLevelSilent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PLANNER_SILENT", tt.silentEnv)
			os.Setenv("PLANNER_LOG_LEVEL", tt.levelEnv)

			// Create a new logger to pick up environment
			logger := core.NewLogger("[test] ")

			if logger.Level != tt.wantLevel {
				t.Errorf("logger.Level = %d, want %d", logger.Level, tt.wantLevel)
			}
		})
	}
}
