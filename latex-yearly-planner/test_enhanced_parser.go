package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"latex-yearly-planner/internal/data"
)

func main() {
	fmt.Println("=== Enhanced CSV Parser Test Suite ===")
	fmt.Println()

	// Test files
	testFiles := []string{
		"/Users/aaron/Downloads/gantt/input/data.cleaned.csv",
		"/Users/aaron/Downloads/gantt/input/data.cleaned.fixed.csv",
		"/Users/aaron/Downloads/gantt/input/test_single.csv",
		"/Users/aaron/Downloads/gantt/input/test_triple.csv",
	}

	for _, filePath := range testFiles {
		fmt.Printf("Testing file: %s\n", filepath.Base(filePath))
		fmt.Println(strings.Repeat("-", 50))

		// Create reader with enhanced options
		opts := &data.ReaderOptions{
			StrictMode:            false,
			SkipInvalid:           true,
			MaxMemoryMB:           100,
			Logger:                log.New(os.Stdout, "[TEST] ", log.LstdFlags),
			ValidateDependencies:  true,
			DetectCircularDeps:    true,
		}

		reader := data.NewReaderWithOptions(filePath, opts)

		// Test CSV format validation
		fmt.Println("1. Validating CSV format...")
		if err := reader.ValidateCSVFormat(); err != nil {
			fmt.Printf("   ❌ CSV format validation failed: %v\n", err)
		} else {
			fmt.Println("   ✅ CSV format validation passed")
		}

		// Test task parsing
		fmt.Println("2. Parsing tasks...")
		tasks, err := reader.ReadTasks()
		if err != nil {
			fmt.Printf("   ❌ Task parsing failed: %v\n", err)
		} else {
			fmt.Printf("   ✅ Successfully parsed %d tasks\n", len(tasks))
		}

		// Test error collection
		if reader.hasErrors() {
			fmt.Println("3. Error summary:")
			fmt.Printf("   %s\n", reader.getErrorSummary())
		} else {
			fmt.Println("3. ✅ No errors detected")
		}

		// Test task analysis
		if len(tasks) > 0 {
			fmt.Println("4. Task analysis:")
			
			// Count tasks by category
			categories := make(map[string]int)
			milestones := 0
			withDependencies := 0
			withParent := 0
			
			for _, task := range tasks {
				categories[task.Category]++
				if task.IsMilestone {
					milestones++
				}
				if len(task.Dependencies) > 0 {
					withDependencies++
				}
				if task.ParentID != "" {
					withParent++
				}
			}
			
			fmt.Printf("   - Categories: %v\n", categories)
			fmt.Printf("   - Milestones: %d\n", milestones)
			fmt.Printf("   - Tasks with dependencies: %d\n", withDependencies)
			fmt.Printf("   - Tasks with parent: %d\n", withParent)
			
			// Show sample tasks with enhanced features
			fmt.Println("5. Sample tasks with enhanced features:")
			count := 0
			for _, task := range tasks {
				if count >= 3 {
					break
				}
				if len(task.Dependencies) > 0 || task.ParentID != "" || task.IsMilestone {
					fmt.Printf("   Task %s: %s\n", task.ID, task.Name)
					if len(task.Dependencies) > 0 {
						fmt.Printf("     Dependencies: %v\n", task.Dependencies)
					}
					if task.ParentID != "" {
						fmt.Printf("     Parent: %s\n", task.ParentID)
					}
					if task.IsMilestone {
						fmt.Printf("     Milestone: Yes\n")
					}
					count++
				}
			}
		}

		fmt.Println()
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println()
	}

	fmt.Println("=== Test Suite Complete ===")
}
