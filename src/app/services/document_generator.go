// Package services provides service classes for the document generation system.
//
// DocumentGenerator handles the complete document generation pipeline including
// page composition, module validation, and file output operations.
package services

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"phd-dissertation-planner/src/core"
)

// DocumentGenerator handles the complete document generation pipeline
type DocumentGenerator struct {
	templateManager *TemplateManager
	logger          core.Logger
}

// GenerationOptions holds options for document generation
type GenerationOptions struct {
	Preview bool
	OutputDir string
}

// GenerationResult represents the result of document generation
type GenerationResult struct {
	GeneratedFiles []string
	PageCount      int
	ModuleCount    int
	Errors         []error
}

// NewDocumentGenerator creates a new document generator
func NewDocumentGenerator(tm *TemplateManager) *DocumentGenerator {
	return &DocumentGenerator{
		templateManager: tm,
		logger:          *core.NewDefaultLogger(),
	}
}

// GenerateDocument generates the complete document from configuration
func (dg *DocumentGenerator) GenerateDocument(cfg core.Config, options GenerationOptions) (*GenerationResult, error) {
	result := &GenerationResult{
		GeneratedFiles: make([]string, 0),
		Errors:         make([]error, 0),
	}

	// Setup output directory
	if err := dg.setupOutputDirectory(cfg.OutputDir); err != nil {
		return result, fmt.Errorf("failed to setup output directory: %w", err)
	}

	// Generate pages
	for _, file := range cfg.Pages {
		pageResult, err := dg.generateSinglePage(cfg, file, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to generate page %s: %w", file.Name, err))
			continue
		}
		
		result.GeneratedFiles = append(result.GeneratedFiles, pageResult.FilePath)
		result.ModuleCount += pageResult.ModuleCount
	}

	result.PageCount = len(cfg.Pages)
	return result, nil
}

// generateSinglePage generates a single page file
func (dg *DocumentGenerator) generateSinglePage(cfg core.Config, file core.Page, options GenerationOptions) (*PageGenerationResult, error) {
	wr := &bytes.Buffer{}

	// Compose all modules for this page
	modules, err := dg.composePageModules(cfg, file, options.Preview)
	if err != nil {
		return nil, fmt.Errorf("failed to compose modules: %w", err)
	}

	// Validate module alignment
	if err := dg.validateModuleAlignment(modules, file.Name); err != nil {
		return nil, fmt.Errorf("module validation failed: %w", err)
	}

	// Render modules to buffer
	if err := dg.templateManager.RenderModules(wr, modules, file); err != nil {
		return nil, fmt.Errorf("failed to render modules: %w", err)
	}

	// Write page file
	filePath, err := dg.writePageFile(cfg.OutputDir, file.Name, wr.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write page file: %w", err)
	}

	return &PageGenerationResult{
		FilePath:    filePath,
		ModuleCount: dg.countModules(modules),
	}, nil
}

// composePageModules composes all modules for a page by calling composer functions
func (dg *DocumentGenerator) composePageModules(cfg core.Config, file core.Page, preview bool) ([]core.Modules, error) {
	var modules []core.Modules

	for _, block := range file.RenderBlocks {
		fn, ok := core.ComposerMap[block.FuncName]
		if !ok {
			return nil, fmt.Errorf("unknown composer function %q - check configuration", block.FuncName)
		}

		blockModules, err := fn(cfg, block.Tpls)
		if err != nil {
			return nil, fmt.Errorf("failed to compose modules for %q: %w", block.FuncName, err)
		}

		// Only one page per unique module if preview flag is enabled
		if preview {
			blockModules = core.FilterUniqueModules(blockModules)
		}

		modules = append(modules, blockModules)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("no modules generated for page %q", file.Name)
	}

	return modules, nil
}

// validateModuleAlignment ensures all module arrays have the same length
func (dg *DocumentGenerator) validateModuleAlignment(modules []core.Modules, pageName string) error {
	if len(modules) == 0 {
		return nil
	}

	expectedLen := len(modules[0])
	for _, mods := range modules {
		if len(mods) != expectedLen {
			return fmt.Errorf("module alignment error for page %q: expected %d modules, got %d", pageName, expectedLen, len(mods))
		}
	}

	return nil
}

// countModules counts the total number of modules across all module arrays
func (dg *DocumentGenerator) countModules(modules []core.Modules) int {
	total := 0
	for _, mods := range modules {
		total += len(mods)
	}
	return total
}

// setupOutputDirectory ensures the output directory exists
func (dg *DocumentGenerator) setupOutputDirectory(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return core.NewFileError(outputDir, "create directory", err)
	}
	dg.logger.Info("Output directory: %s", outputDir)
	return nil
}

// writePageFile writes the page content to a file
func (dg *DocumentGenerator) writePageFile(outputDir, pageName string, content []byte) (string, error) {
	pageFile := filepath.Join(outputDir, pageName+".tex")
	if err := os.WriteFile(pageFile, content, 0o600); err != nil {
		return "", core.NewFileError(pageFile, "write", err)
	}
	dg.logger.Info("Generated page: %s", pageFile)
	return pageFile, nil
}

// GenerateRootDocument creates the main LaTeX document file
func (dg *DocumentGenerator) GenerateRootDocument(cfg core.Config, pathConfigs []string) (string, error) {
	wr := &bytes.Buffer{}

	if err := dg.templateManager.ExecuteDocument(wr, cfg); err != nil {
		return "", core.NewTemplateError("document.tpl", 0, "failed to generate LaTeX document", err)
	}

	outputFile := filepath.Join(cfg.OutputDir, dg.rootFilename(pathConfigs[len(pathConfigs)-1]))
	if err := os.WriteFile(outputFile, wr.Bytes(), 0o600); err != nil {
		return "", core.NewFileError(outputFile, "write", err)
	}
	dg.logger.Info("Generated LaTeX file: %s", outputFile)
	return outputFile, nil
}

// rootFilename generates the root filename from the config path
func (dg *DocumentGenerator) rootFilename(pathConfig string) string {
	filename := filepath.Base(pathConfig)
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + ".tex"
}

// PageGenerationResult represents the result of generating a single page
type PageGenerationResult struct {
	FilePath    string
	ModuleCount int
}
