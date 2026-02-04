package app

import (
	"testing"
)

func TestEscapeLatex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple text", "Simple text"},
		{"Text with & and %", "Text with \\& and \\%"},
		{"$100 #tag", "\\$100 \\#tag"},
		{"{braces}", "\\{braces\\}"},
		{"_underscore", "\\_underscore"},
		{"Nested { braces { inside } }", "Nested \\{ braces \\{ inside \\} \\}"},
		{"$$%%&&##__", "\\$\\$\\%\\%\\&\\&\\#\\#\\_\\_"},
	}

	for _, tt := range tests {
		actual := EscapeLatex(tt.input)
		if actual != tt.expected {
			t.Errorf("EscapeLatex(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func BenchmarkEscapeLatex(b *testing.B) {
	input := "Project Setup & Proposal: 50% completed_task #123 {Draft} $100"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EscapeLatex(input)
	}
}
