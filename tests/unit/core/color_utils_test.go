// Package core_test provides tests for color utilities
package core_test

import (
	"testing"

	"phd-dissertation-planner/internal/core"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "valid hex with hash",
			hex:      "#FF0000",
			expected: "255,0,0",
		},
		{
			name:     "valid hex without hash",
			hex:      "00FF00",
			expected: "0,255,0",
		},
		{
			name:     "blue color",
			hex:      "#0000FF",
			expected: "0,0,255",
		},
		{
			name:     "white",
			hex:      "#FFFFFF",
			expected: "255,255,255",
		},
		{
			name:     "black",
			hex:      "#000000",
			expected: "0,0,0",
		},
		{
			name:     "gray",
			hex:      "#808080",
			expected: "128,128,128",
		},
		{
			name:     "short hex",
			hex:      "#FFF",
			expected: "128,128,128",
		},
		{
			name:     "long hex",
			hex:      "#FF00000",
			expected: "128,128,128",
		},
		{
			name:     "invalid characters",
			hex:      "#GG0000",
			expected: "128,128,128",
		},
		{
			name:     "empty string",
			hex:      "",
			expected: "128,128,128",
		},
		{
			name:     "just hash",
			hex:      "#",
			expected: "128,128,128",
		},
		{
			name:     "lowercase",
			hex:      "#ff0000",
			expected: "255,0,0",
		},
		{
			name:     "mixed case",
			hex:      "#AbCdEf",
			expected: "171,205,239",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.HexToRGB(tt.hex)
			if result != tt.expected {
				t.Errorf("core.HexToRGB(%q) = %q, want %q", tt.hex, result, tt.expected)
			}
		})
	}
}
