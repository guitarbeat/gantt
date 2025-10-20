package core_test

import (
	"errors"
	"strings"
	"testing"

	"phd-dissertation-planner/internal/core"
)

func TestConfigError(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		field       string
		message     string
		err         error
		wantContain []string
	}{
		{
			name:        "with field",
			file:        "config.yaml",
			field:       "weekstart",
			message:     "invalid value",
			err:         errors.New("parse error"),
			wantContain: []string{"config.yaml", "weekstart", "invalid value"},
		},
		{
			name:        "without field",
			file:        "config.yaml",
			field:       "",
			message:     "file not found",
			err:         nil,
			wantContain: []string{"config.yaml", "file not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Newcore.ConfigError(tt.file, tt.field, tt.message, tt.err)

			errStr := err.Error()
			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("ConfigError.Error() = %q, should contain %q", errStr, want)
				}
			}

			// Test Unwrap
			if tt.err != nil {
				if unwrapped := err.Unwrap(); unwrapped != tt.err {
					t.Errorf("ConfigError.Unwrap() = %v, want %v", unwrapped, tt.err)
				}
			}
		})
	}
}

func TestFileError(t *testing.T) {
	baseErr := errors.New("permission denied")
	err := NewFileError("/path/to/file.txt", "write", baseErr)

	errStr := err.Error()
	if !strings.Contains(errStr, "/path/to/file.txt") {
		t.Errorf("FileError should contain file path")
	}
	if !strings.Contains(errStr, "write") {
		t.Errorf("FileError should contain operation")
	}

	// Test Unwrap
	if unwrapped := err.Unwrap(); unwrapped != baseErr {
		t.Errorf("FileError.Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestTemplateError(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		line        int
		message     string
		wantContain []string
	}{
		{
			name:        "with line number",
			template:    "document.tpl",
			line:        42,
			message:     "syntax error",
			wantContain: []string{"document.tpl", "42", "syntax error"},
		},
		{
			name:        "without line number",
			template:    "page.tpl",
			line:        0,
			message:     "execution failed",
			wantContain: []string{"page.tpl", "execution failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewTemplateError(tt.template, tt.line, tt.message, nil)

			errStr := err.Error()
			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("TemplateError.Error() = %q, should contain %q", errStr, want)
				}
			}
		})
	}
}

func TestDataError(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		row         int
		column      string
		message     string
		wantContain []string
	}{
		{
			name:        "with row and column",
			source:      "data.csv",
			row:         10,
			column:      "StartDate",
			message:     "invalid date",
			wantContain: []string{"data.csv", "10", "StartDate", "invalid date"},
		},
		{
			name:        "with row only",
			source:      "data.csv",
			row:         5,
			column:      "",
			message:     "parse error",
			wantContain: []string{"data.csv", "5", "parse error"},
		},
		{
			name:        "without row",
			source:      "data.csv",
			row:         0,
			column:      "",
			message:     "file corrupt",
			wantContain: []string{"data.csv", "file corrupt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDataError(tt.source, tt.row, tt.column, tt.message, nil)

			errStr := err.Error()
			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("DataError.Error() = %q, should contain %q", errStr, want)
				}
			}
		})
	}
}

func TestErrorAggregator(t *testing.T) {
	t.Run("empty aggregator", func(t *testing.T) {
		agg := NewErrorAggregator()

		if agg.HasErrors() {
			t.Error("Empty aggregator should not have errors")
		}
		if agg.HasWarnings() {
			t.Error("Empty aggregator should not have warnings")
		}
		if agg.ErrorCount() != 0 {
			t.Errorf("ErrorCount() = %d, want 0", agg.ErrorCount())
		}
		if agg.WarningCount() != 0 {
			t.Errorf("WarningCount() = %d, want 0", agg.WarningCount())
		}
	})

	t.Run("add errors", func(t *testing.T) {
		agg := NewErrorAggregator()

		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		agg.AddError(err1)
		agg.AddError(err2)

		if !agg.HasErrors() {
			t.Error("Aggregator should have errors")
		}
		if agg.ErrorCount() != 2 {
			t.Errorf("ErrorCount() = %d, want 2", agg.ErrorCount())
		}
	})

	t.Run("add warnings", func(t *testing.T) {
		agg := NewErrorAggregator()

		warn1 := errors.New("warning 1")
		warn2 := errors.New("warning 2")

		agg.AddWarning(warn1)
		agg.AddWarning(warn2)

		if !agg.HasWarnings() {
			t.Error("Aggregator should have warnings")
		}
		if agg.WarningCount() != 2 {
			t.Errorf("WarningCount() = %d, want 2", agg.WarningCount())
		}
	})

	t.Run("summary", func(t *testing.T) {
		agg := NewErrorAggregator()

		agg.AddError(errors.New("error 1"))
		agg.AddWarning(errors.New("warning 1"))

		summary := agg.Summary()
		if !strings.Contains(summary, "Errors") {
			t.Error("Summary should contain 'Errors'")
		}
		if !strings.Contains(summary, "Warnings") {
			t.Error("Summary should contain 'Warnings'")
		}
	})

	t.Run("clear", func(t *testing.T) {
		agg := NewErrorAggregator()

		agg.AddError(errors.New("error"))
		agg.AddWarning(errors.New("warning"))

		agg.Clear()

		if agg.HasErrors() || agg.HasWarnings() {
			t.Error("Clear() should remove all errors and warnings")
		}
	})

	t.Run("error interface", func(t *testing.T) {
		agg := NewErrorAggregator()

		// Empty aggregator
		if agg.Error() != "no errors" {
			t.Errorf("Empty aggregator Error() = %s, want 'no errors'", agg.Error())
		}

		// Single error
		agg.AddError(errors.New("test error"))
		errStr := agg.Error()
		if !strings.Contains(errStr, "test error") {
			t.Errorf("Error() should contain the error message")
		}

		// Multiple errors
		agg.AddError(errors.New("another error"))
		errStr = agg.Error()
		if !strings.Contains(errStr, "2 errors") {
			t.Errorf("Error() should indicate multiple errors")
		}
	})
}
