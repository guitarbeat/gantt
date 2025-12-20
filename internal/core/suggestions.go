package core

import (
	"strings"
)

// SuggestCorrection finds the closest match from a list of valid options
func SuggestCorrection(input string, options []string) string {
	if input == "" || len(options) == 0 {
		return ""
	}

	input = strings.ToLower(input)
	bestMatch := ""
	bestDistance := 1000

	for _, option := range options {
		// Calculate Levenshtein distance
		dist := levenshtein(input, strings.ToLower(option))

		// If exact match (should generally be handled before calling this), return it
		if dist == 0 {
			return option
		}

		// Update best match if this is closer
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = option
		}
	}

	// Only return if it's a reasonable match:
	// 1. Distance is small (<= 4 edits to allow "in prog" -> "in progress")
	// 2. Distance is not too large relative to word length (avoid suggesting "a" for "bcd")
	if bestDistance <= 4 && float64(bestDistance) < float64(len(input))*0.8 {
		return bestMatch
	}

	return ""
}

// levenshtein calculates the Levenshtein distance between two strings
func levenshtein(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	// Optimization: use two rows instead of full matrix
	prev := make([]int, m+1)
	curr := make([]int, m+1)

	for j := 0; j <= m; j++ {
		prev[j] = j
	}

	for i := 1; i <= n; i++ {
		curr[0] = i
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}

			// Min of deletion, insertion, substitution
			curr[j] = min(
				prev[j]+1,     // deletion
				min(
					curr[j-1]+1, // insertion
					prev[j-1]+cost, // substitution
				),
			)
		}
		copy(prev, curr)
	}

	return prev[m]
}
