package app

import (
	"testing"
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
			result, err := dictFunc(tt.input...)

			if (err != nil) != tt.wantErr {
				t.Errorf("dictFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			if len(result) != tt.wantLen {
				t.Errorf("dictFunc() returned map of length %d, want %d", len(result), tt.wantLen)
			}

			// Check specific key-value pairs
			for key, expectedValue := range tt.checkKeys {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("dictFunc() missing key %q", key)
				} else if actualValue != expectedValue {
					t.Errorf("dictFunc() key %q = %v, want %v", key, actualValue, expectedValue)
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
			if got := incrFunc(tt.input); got != tt.want {
				t.Errorf("incrFunc(%d) = %d, want %d", tt.input, got, tt.want)
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
			if got := decFunc(tt.input); got != tt.want {
				t.Errorf("decFunc(%d) = %d, want %d", tt.input, got, tt.want)
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
			if got := isFunc(tt.input); got != tt.want {
				t.Errorf("isFunc(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTemplateFuncs(t *testing.T) {
	funcs := TemplateFuncs()

	// Check that all expected functions are present
	expectedFuncs := []string{"dict", "incr", "dec", "is", "hypertarget"}
	for _, name := range expectedFuncs {
		if _, exists := funcs[name]; !exists {
			t.Errorf("TemplateFuncs() missing function %q", name)
		}
	}

	// Check that we have the expected number of functions
	if len(funcs) != len(expectedFuncs) {
		t.Errorf("TemplateFuncs() returned %d functions, want %d", len(funcs), len(expectedFuncs))
	}
}
