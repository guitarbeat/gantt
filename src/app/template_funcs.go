package app

import (
	"errors"
	"text/template"
)

// TemplateFuncs returns a FuncMap with all custom template functions
// These functions are available to all templates during rendering
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"dict": dictFunc,
		"incr": incrFunc,
		"dec":  decFunc,
		"is":   isFunc,
	}
}

// dictFunc creates a map from key-value pairs for use in templates
// Usage: {{ template "name" (dict "key1" value1 "key2" value2) }}
func dictFunc(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("dict requires an even number of arguments (key-value pairs)")
	}

	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

// incrFunc increments an integer by 1
// Usage: {{ incr .Index }}
func incrFunc(i int) int {
	return i + 1
}

// decFunc decrements an integer by 1
// Usage: {{ dec .Index }}
func decFunc(i int) int {
	return i - 1
}

// isFunc checks if a value is truthy
// Returns true for:
// - bool true
// - non-nil values
// Returns false for:
// - bool false
// - nil values
// Usage: {{ if is .Value }}...{{ end }}
func isFunc(i interface{}) bool {
	// Handle explicit boolean values
	if value, ok := i.(bool); ok {
		return value
	}

	// All other non-nil values are truthy
	return i != nil
}

// Additional template helper functions can be added here
// Examples:
// - formatDate: Format time.Time values
// - join: Join string slices
// - default: Provide default values
// - escape: Escape special characters
