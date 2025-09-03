package main

import (
	"fmt"
	"time"
)

func main() {
	// Test with UTC
	t := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("UTC: August 1, 2025 is %s\n", t.Weekday())

	// Test with Local (CDT)
	t2 := time.Date(2025, 8, 1, 0, 0, 0, 0, time.Local)
	fmt.Printf("Local: August 1, 2025 is %s\n", t2.Weekday())

	// Test the shift calculation with Local
	wd := time.Monday // weekstart = 1
	weekday := t2.Weekday()
	shift := (7 + int(weekday) - int(wd)) % 7
	fmt.Printf("wd=%d, weekday=%d, shift=%d\n", wd, weekday, shift)
}
