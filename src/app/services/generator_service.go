// Package services provides service classes for the document generation system.
//
// GeneratorService orchestrates the complete document generation process using
// the various service classes to provide a clean, high-level interface.
package services

import (
	"fmt"
	"strings"

	"phd-dissertation-planner/src/core"

	"github.com/urfave/cli/v2"
)

// GeneratorService orchestrates the complete document generation process
type GeneratorService struct {
	configLoader      *ConfigLoader
	templateManager   *TemplateManager
	documentGenerator *DocumentGenerator
	coverageAnalyzer  *CoverageAnalyzer
	logger            core.Logger
}

// GenerationRequest represents a request for document generation
type GenerationRequest struct {
	ConfigPaths []string
	OutputDir   string
	Preview     bool
	RunTests    bool
}

// GenerationResponse represents the response from document generation
type GenerationResponse struct {
	Success        bool
	GeneratedFiles []string
	PageCount      int
	ModuleCount    int
	CoverageResult *CoverageAnalysisResult
	Errors         []error
	Summary        string
}

// NewGeneratorService creates a new generator service with all dependencies
func NewGeneratorService() (*GeneratorService, error) {
	// Create template manager
	tm, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create template manager: %w", err)
	}

	// Create other services
	cl := NewConfigLoader()
	dg := NewDocumentGenerator(tm)
	ca := NewCoverageAnalyzer()

	return &GeneratorService{
		configLoader:      cl,
		templateManager:   tm,
		documentGenerator: dg,
		coverageAnalyzer:  ca,
		logger:            *core.NewDefaultLogger(),
	}, nil
}

// ProcessCLIRequest processes a CLI request and generates documents
func (gs *GeneratorService) ProcessCLIRequest(c *cli.Context) (*GenerationResponse, error) {
	// Check if test coverage is requested
	if c.Bool("test-coverage") {
		return gs.runTestCoverage()
	}

	// Load configuration
	configResult, err := gs.configLoader.LoadConfiguration(c)
	if err != nil {
		return &GenerationResponse{
			Success: false,
			Errors:  []error{err},
		}, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := gs.configLoader.ValidateConfiguration(configResult.Config); err != nil {
		return &GenerationResponse{
			Success: false,
			Errors:  []error{err},
		}, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log configuration summary
	gs.logger.Info(gs.configLoader.GetConfigurationSummary(configResult))

	// Generate documents
	options := GenerationOptions{
		Preview:   c.Bool("preview"),
		OutputDir: configResult.Config.OutputDir,
	}

	// Generate root document
	rootFile, err := gs.documentGenerator.GenerateRootDocument(configResult.Config, configResult.PathConfigs)
	if err != nil {
		return &GenerationResponse{
			Success: false,
			Errors:  []error{err},
		}, fmt.Errorf("failed to generate root document: %w", err)
	}

	// Generate pages
	pageResult, err := gs.documentGenerator.GenerateDocument(configResult.Config, options)
	if err != nil {
		return &GenerationResponse{
			Success: false,
			Errors:  []error{err},
		}, fmt.Errorf("failed to generate pages: %w", err)
	}

	// Combine results
	allFiles := append([]string{rootFile}, pageResult.GeneratedFiles...)

	response := &GenerationResponse{
		Success:        true,
		GeneratedFiles: allFiles,
		PageCount:      pageResult.PageCount + 1, // +1 for root document
		ModuleCount:    pageResult.ModuleCount,
		Errors:         pageResult.Errors,
		Summary:        gs.generateSummary(allFiles, pageResult.PageCount+1, pageResult.ModuleCount),
	}

	return response, nil
}

// runTestCoverage runs test coverage analysis
func (gs *GeneratorService) runTestCoverage() (*GenerationResponse, error) {
	coverageResult, err := gs.coverageAnalyzer.RunTestCoverage()
	if err != nil {
		return &GenerationResponse{
			Success: false,
			Errors:  []error{err},
		}, err
	}

	summary := gs.coverageAnalyzer.GetCoverageSummary(coverageResult)

	return &GenerationResponse{
		Success:        coverageResult.Success,
		CoverageResult: coverageResult,
		Summary:        summary,
		Errors:         coverageResult.Errors,
	}, nil
}

// generateSummary creates a summary of the generation process
func (gs *GeneratorService) generateSummary(files []string, pageCount, moduleCount int) string {
	summary := fmt.Sprintf("Document Generation Complete:\n")
	summary += fmt.Sprintf("  Generated Files: %d\n", len(files))
	summary += fmt.Sprintf("  Pages: %d\n", pageCount)
	summary += fmt.Sprintf("  Modules: %d\n", moduleCount)
	
	if len(files) > 0 {
		summary += fmt.Sprintf("  Output Directory: %s\n", files[0][:len(files[0])-len(files[0][strings.LastIndex(files[0], "/"):])])
	}

	return summary
}

// GetServiceStatus returns the status of all services
func (gs *GeneratorService) GetServiceStatus() map[string]bool {
	status := make(map[string]bool)
	
	// Check template manager
	status["template_manager"] = gs.templateManager != nil
	
	// Check config loader
	status["config_loader"] = gs.configLoader != nil
	
	// Check document generator
	status["document_generator"] = gs.documentGenerator != nil
	
	// Check coverage analyzer
	status["coverage_analyzer"] = gs.coverageAnalyzer != nil
	
	return status
}

// ReloadTemplates reloads templates (useful for development)
func (gs *GeneratorService) ReloadTemplates() error {
	tm, err := NewTemplateManager()
	if err != nil {
		return fmt.Errorf("failed to reload templates: %w", err)
	}
	gs.templateManager = tm
	gs.documentGenerator.templateManager = tm
	return nil
}
