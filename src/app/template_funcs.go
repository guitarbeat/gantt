// Package app - Template functions provide custom helpers for template rendering.
//
// This module defines all custom template functions available in the template
// rendering engine. Each function is designed to simplify common template operations.
//
// Available functions:
//
// dict: Create maps from key-value pairs
//
//	Usage: {{ template "name" (dict "key1" value1 "key2" value2) }}
//	Useful for passing multiple parameters to templates
//
// incr: Increment an integer by 1
//
//	Usage: {{ incr .Index }}
//	Common for loop counters and pagination
//
// dec: Decrement an integer by 1
//
//	Usage: {{ dec .Index }}
//	Useful for zero-based indexing
//
// is: Check if a value is truthy
//
//	Usage: {{ if is .Value }}...{{ end }}
//	Returns true for non-nil values and explicit true booleans
//	Returns false for nil and explicit false booleans
//
// All functions are thoroughly tested with 100% code coverage.
// See template_funcs_test.go for comprehensive test cases.
//
// Example template usage:
//
//	{{/* Pass multiple values to a sub-template */}}
//	{{ template "header" (dict "title" "My Page" "subtitle" "Details") }}
//
//	{{/* Loop with 1-based indexing */}}
//	{{ range $i, $item := .Items }}
//	    {{ incr $i }}. {{ $item.Name }}
//	{{ end }}
//
//	{{/* Conditional rendering */}}
//	{{ if is .OptionalField }}
//	    Optional content: {{ .OptionalField }}
//	{{ end }}
package app

import (
	"errors"
	"text/template"

	"phd-dissertation-planner/src/shared/templates"
)

// TemplateFuncs returns a FuncMap with all custom template functions
// These functions are available to all templates during rendering
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"dict":       dictFunc,
		"incr":       incrFunc,
		"dec":        decFunc,
		"is":         isFunc,
		"hypertarget": templates.Hypertarget,
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
