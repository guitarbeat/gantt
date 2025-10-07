// Package services provides service classes for the document generation system.
//
// This file demonstrates how to use the service classes together to generate documents.
// This is an example and should not be used in production - use the main CLI instead.
package services

import (
	"fmt"
	"log"

	"phd-dissertation-planner/src/core"
)

// ExampleUsage demonstrates how to use the service classes
func ExampleUsage() {
	// Create the main generator service
	gs, err := NewGeneratorService()
	if err != nil {
		log.Fatalf("Failed to create generator service: %v", err)
	}

	// Check service status
	status := gs.GetServiceStatus()
	fmt.Printf("Service Status: %+v\n", status)

	// Create a sample configuration
	cfg := core.Config{
		OutputDir: "example_output",
		StartYear: 2025,
		EndYear:   2026,
		Pages: []core.Page{
			{
				Name: "example_page",
				RenderBlocks: []core.RenderBlock{
					{
						FuncName: "monthly",
						Tpls:     []string{"monthly.tpl"},
					},
				},
			},
		},
	}

	// Validate configuration using ConfigLoader
	cl := NewConfigLoader()
	if err := cl.ValidateConfiguration(cfg); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create DocumentGenerator
	tm, err := NewTemplateManager()
	if err != nil {
		log.Fatalf("Failed to create template manager: %v", err)
	}

	dg := NewDocumentGenerator(tm)

	// Generate documents
	options := GenerationOptions{
		Preview:   false,
		OutputDir: cfg.OutputDir,
	}

	result, err := dg.GenerateDocument(cfg, options)
	if err != nil {
		log.Fatalf("Document generation failed: %v", err)
	}

	fmt.Printf("Generation Result: %+v\n", result)
	fmt.Printf("Generated Files: %d\n", len(result.GeneratedFiles))
	fmt.Printf("Pages: %d\n", result.PageCount)
	fmt.Printf("Modules: %d\n", result.ModuleCount)

	// Example of running test coverage
	ca := NewCoverageAnalyzer()
	coverageResult, err := ca.RunTestCoverage()
	if err != nil {
		log.Printf("Coverage analysis failed: %v", err)
	} else {
		fmt.Printf("Coverage Analysis: %+v\n", coverageResult)
	}
}
