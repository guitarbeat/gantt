package app_test

import (
	"testing"

	"phd-dissertation-planner/internal/app"
)

func TestDictFunc(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		wantErr   bool
		wantLen   int
		checkKeys map[string]interface{}
	}{
		{
			name:    "valid key-value pairs",
			input:   []interface{}{"key1", "value1", "key2", 42},
			wantErr: false,
			wantLen: 2,
			checkKeys: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
		},
		{
			name:    "empty input",
			input:   []interface{}{},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "odd number of arguments",
			input:   []interface{}{"key1", "value1", "key2"},
			wantErr: true,
		},
		{
			name:    "non-string key",
			input:   []interface{}{123, "value1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := app.dictFunc(tt.input...)

			if (err != nil) != tt.wantErr {
				t.Errorf("app.dictFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			if len(result) != tt.wantLen {
				t.Errorf("app.dictFunc() returned map of length %d, want %d", len(result), tt.wantLen)
			}

			// Check specific key-value pairs
			for key, expectedValue := range tt.checkKeys {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("app.dictFunc() missing key %q", key)
				} else if actualValue != expectedValue {
					t.Errorf("app.dictFunc() key %q = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestIncrFunc(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 5, 6},
		{"zero", 0, 1},
		{"negative", -1, 0},
		{"large", 999, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.incrFunc(tt.input); got != tt.want {
				t.Errorf("app.incrFunc(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestDecFunc(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 5, 4},
		{"zero", 0, -1},
		{"negative", -1, -2},
		{"large", 1000, 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.decFunc(tt.input); got != tt.want {
				t.Errorf("app.decFunc(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsFunc(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  bool
	}{
		{"bool true", true, true},
		{"bool false", false, false},
		{"nil", nil, false},
		{"string", "hello", true},
		{"empty string", "", true}, // Empty string is still non-nil
		{"int zero", 0, true},      // Zero is still non-nil
		{"int positive", 42, true},
		{"pointer nil", (*int)(nil), true}, // Typed nil is still non-nil interface value
		{"pointer non-nil", new(int), true},
		{"slice nil", []int(nil), true}, // Typed nil slice is still non-nil interface value
		{"slice empty", []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.isFunc(tt.input); got != tt.want {
				t.Errorf("app.isFunc(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTemplateFuncs(t *testing.T) {
	funcs := app.TemplateFuncs()

	// Check that all expected functions are present
	expectedFuncs := []string{"dict", "incr", "dec", "is", "hypertarget", "lower", "plus", "mod", "replace"}
	for _, name := range expectedFuncs {
		if _, exists := funcs[name]; !exists {
			t.Errorf("app.TemplateFuncs() missing function %q", name)
		}
	}

	// Check that we have the expected number of functions
	if len(funcs) != len(expectedFuncs) {
		t.Errorf("app.TemplateFuncs() returned %d functions, want %d", len(funcs), len(expectedFuncs))
	}
}

func TestRootFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple yaml file",
			input:    "config.yaml",
			expected: "config.tex",
		},
		{
			name:     "nested path",
			input:    "configs/base.yaml",
			expected: "base.tex",
		},
		{
			name:     "no extension",
			input:    "config",
			expected: "config.tex",
		},
		{
			name:     "multiple dots",
			input:    "config.special.yaml",
			expected: "config.special.tex",
		},
		{
			name:     "yml extension",
			input:    "config.yml",
			expected: "config.tex",
		},
		{
			name:     "json file",
			input:    "data.json",
			expected: "data.tex",
		},
		{
			name:     "empty string",
			input:    "",
			expected: ".tex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.RootFilename(tt.input)
			if result != tt.expected {
				t.Errorf("app.RootFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeLatex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special chars",
			input:    "Simple text",
			expected: "Simple text",
		},
		{
			name:     "ampersand",
			input:    "Tom & Jerry",
			expected: "Tom \\& Jerry",
		},
		{
			name:     "percent",
			input:    "100% complete",
			expected: "100\\% complete",
		},
		{
			name:     "dollar",
			input:    "$100 budget",
			expected: "\\$100 budget",
		},
		{
			name:     "hash",
			input:    "#1 priority",
			expected: "\\#1 priority",
		},
		{
			name:     "underscore",
			input:    "task_name",
			expected: "task\\_name",
		},
		{
			name:     "curly braces",
			input:    "{bold} text",
			expected: "\\{bold\\} text",
		},
		{
			name:     "multiple chars",
			input:    "Cost: $100 & 50% done_#1",
			expected: "Cost: \\$100 \\& 50\\% done\\_\\#1",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "backslash already present",
			input:    "already \\& escaped",
			expected: "already \\\\& escaped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.escapeLatex(tt.input)
			if result != tt.expected {
				t.Errorf("app.escapeLatex(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateCSVPriority(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected int
	}{
		{
			name:     "comprehensive file",
			filename: "research_timeline_v5_comprehensive.csv",
			expected: 10,
		},
		{
			name:     "v5.1 file",
			filename: "data_v5.1_tasks.csv",
			expected: 8,
		},
		{
			name:     "v5 file",
			filename: "timeline_v5.csv",
			expected: 6,
		},
		{
			name:     "regular file",
			filename: "tasks.csv",
			expected: 0,
		},
		{
			name:     "case insensitive comprehensive",
			filename: "COMPREHENSIVE_DATA.csv",
			expected: 10,
		},
		{
			name:     "case insensitive version",
			filename: "V5.1_EXPORT.csv",
			expected: 8,
		},
		{
			name:     "no extension",
			filename: "comprehensive_data",
			expected: 10,
		},
		{
			name:     "partial match comprehensive",
			filename: "mycomprehensivedata.csv",
			expected: 10,
		},
		{
			name:     "partial match v5.1",
			filename: "backup_v5.1_final.csv",
			expected: 8,
		},
		{
			name:     "partial match v5",
			filename: "archive_v5_old.csv",
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCSVPriority(tt.filename)
			if result != tt.expected {
				t.Errorf("calculateCSVPriority(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestLowerFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mixed case",
			input:    "Hello World",
			expected: "hello world",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "uppercase only",
			input:    "HELLO WORLD",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "special characters",
			input:    "Hello@World#123",
			expected: "hello@world#123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lowerFunc(tt.input)
			if result != tt.expected {
				t.Errorf("lowerFunc(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPlusFunc(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        5,
			b:        3,
			expected: 8,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: -5,
		},
		{
			name:     "mixed signs",
			a:        10,
			b:        -4,
			expected: 6,
		},
		{
			name:     "zero values",
			a:        0,
			b:        5,
			expected: 5,
		},
		{
			name:     "both zero",
			a:        0,
			b:        0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plusFunc(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("plusFunc(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestModFunc(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "exact division",
			a:        10,
			b:        5,
			expected: 0,
		},
		{
			name:     "remainder",
			a:        17,
			b:        5,
			expected: 2,
		},
		{
			name:     "negative dividend",
			a:        -17,
			b:        5,
			expected: -2,
		},
		{
			name:     "negative divisor",
			a:        17,
			b:        -5,
			expected: 2,
		},
		{
			name:     "both negative",
			a:        -17,
			b:        -5,
			expected: -2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modFunc(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("modFunc(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestReplaceFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		from     string
		to       string
		expected string
	}{
		{
			name:     "simple replacement",
			input:    "Hello world",
			from:     "world",
			to:       "Go",
			expected: "Hello Go",
		},
		{
			name:     "multiple occurrences",
			input:    "foo bar foo baz foo",
			from:     "foo",
			to:       "bar",
			expected: "bar bar bar baz bar",
		},
		{
			name:     "no occurrences",
			input:    "Hello world",
			from:     "xyz",
			to:       "abc",
			expected: "Hello world",
		},
		{
			name:     "empty from string",
			input:    "Hello world",
			from:     "",
			to:       "X",
			expected: "XHXeXlXlXoX XwXoXrXlXdX",
		},
		{
			name:     "replace with empty",
			input:    "Hello world",
			from:     " ",
			to:       "",
			expected: "Helloworld",
		},
		{
			name:     "empty input",
			input:    "",
			from:     "x",
			to:       "y",
			expected: "",
		},
		{
			name:     "special characters",
			input:    "a&b%c$d",
			from:     "&",
			to:       "AND",
			expected: "aANDb%c$d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFunc(tt.input, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("replaceFunc(%q, %q, %q) = %q, want %q", tt.input, tt.from, tt.to, result, tt.expected)
			}
		})
	}
}
