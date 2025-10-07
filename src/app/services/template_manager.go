// Package services provides service classes for the document generation system.
//
// This package contains service classes that encapsulate specific responsibilities
// and provide clean interfaces for the main application logic.
//
// TemplateManager handles all template-related operations including loading,
// parsing, and executing templates with proper error handling.
package services

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"phd-dissertation-planner/src/app"
	"phd-dissertation-planner/src/core"
	tmplfs "phd-dissertation-planner/src/shared/templates"
)

// TemplateManager handles template loading, parsing, and execution
type TemplateManager struct {
	templates *template.Template
	logger    core.Logger
}

// TemplateManagerConfig holds configuration for the template manager
type TemplateManagerConfig struct {
	DevMode        bool
	TemplateSubDir string
	TemplatePath   string
	TemplatePattern string
}

// NewTemplateManager creates a new template manager with default configuration
func NewTemplateManager() (*TemplateManager, error) {
	config := TemplateManagerConfig{
		DevMode:         os.Getenv("DEV_TEMPLATES") != "",
		TemplateSubDir:  "monthly",
		TemplatePath:    "src/shared/templates/monthly",
		TemplatePattern: "*.tpl",
	}
	return NewTemplateManagerWithConfig(config)
}

// NewTemplateManagerWithConfig creates a new template manager with custom configuration
func NewTemplateManagerWithConfig(config TemplateManagerConfig) (*TemplateManager, error) {
	tm := &TemplateManager{
		logger: *core.NewDefaultLogger(),
	}

	if err := tm.loadTemplates(config); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return tm, nil
}

// loadTemplates loads templates from either embedded files or filesystem
func (tm *TemplateManager) loadTemplates(config TemplateManagerConfig) error {
	// Create template with custom functions
	t := template.New("").Funcs(tm.getTemplateFuncs())

	// Choose source of templates: embedded by default, filesystem when DEV_TEMPLATES is set
	var (
		err   error
		useFS fs.FS
	)

	if config.DevMode {
		// Use on-disk templates for development override
		tm.logger.Debug("Loading templates from filesystem: %s", config.TemplatePath)
		useFS = os.DirFS(filepath.Join("src", "shared", "templates", config.TemplateSubDir))
	} else {
		// Use embedded templates from templates.FS
		tm.logger.Debug("Loading embedded templates from: %s", config.TemplateSubDir)
		// Narrow to monthly/ subdir
		var sub fs.FS
		sub, err = fs.Sub(tmplfs.FS, config.TemplateSubDir)
		if err != nil {
			return fmt.Errorf("failed to access embedded templates directory '%s': %v (check that templates are properly embedded)", config.TemplateSubDir, err)
		}
		useFS = sub
	}

	// Parse all *.tpl templates from the selected FS
	t, err = t.ParseFS(useFS, config.TemplatePattern)
	if err != nil {
		// Provide detailed error message with troubleshooting hints
		if config.DevMode {
			return fmt.Errorf("failed to parse templates from filesystem '%s' with pattern '%s': %v\n"+
				"Check that template files exist and have valid syntax", config.TemplatePath, config.TemplatePattern, err)
		} else {
			return fmt.Errorf("failed to parse embedded templates with pattern '%s': %v\n"+
				"This may indicate a build issue - ensure templates are embedded correctly", config.TemplatePattern, err)
		}
	}

	tm.templates = t
	tm.logger.Debug("Successfully loaded templates with pattern: %s", config.TemplatePattern)
	return nil
}

// getTemplateFuncs returns the template function map
func (tm *TemplateManager) getTemplateFuncs() template.FuncMap {
	// Import template functions from the app package
	return app.TemplateFuncs()
}

// Execute executes a template with the given data
func (tm *TemplateManager) Execute(wr io.Writer, name string, data interface{}) error {
	// Check if template exists before trying to execute
	if tm.templates.Lookup(name) == nil {
		availableTemplates := make([]string, 0)
		for _, tmpl := range tm.templates.Templates() {
			availableTemplates = append(availableTemplates, tmpl.Name())
		}
		return core.NewTemplateError(
			name,
			0,
			fmt.Sprintf("template not found (available: %v)", availableTemplates),
			nil,
		)
	}

	if err := tm.templates.ExecuteTemplate(wr, name, data); err != nil {
		return core.NewTemplateError(name, 0, "failed to execute template", err)
	}

	return nil
}

// ExecuteDocument executes the main document template
func (tm *TemplateManager) ExecuteDocument(wr io.Writer, cfg core.Config) error {
	type pack struct {
		Cfg   core.Config
		Pages []core.Page
	}

	data := pack{Cfg: cfg, Pages: cfg.Pages}
	return tm.Execute(wr, "document.tpl", data)
}

// RenderModules renders all modules to the writer using the template
func (tm *TemplateManager) RenderModules(wr io.Writer, modules []core.Modules, file core.Page) error {
	if len(modules) == 0 {
		return nil
	}

	moduleCount := len(modules[0])
	for i := 0; i < moduleCount; i++ {
		for j, mod := range modules {
			if err := tm.Execute(wr, mod[i].Tpl, mod[i]); err != nil {
				return core.NewTemplateError(
					mod[i].Tpl,
					0,
					fmt.Sprintf("failed to execute template for function %s", file.RenderBlocks[j].FuncName),
					err,
				)
			}
		}
	}

	return nil
}

// GetAvailableTemplates returns a list of available template names
func (tm *TemplateManager) GetAvailableTemplates() []string {
	var templates []string
	for _, tmpl := range tm.templates.Templates() {
		templates = append(templates, tmpl.Name())
	}
	return templates
}

// ValidateTemplate checks if a template exists
func (tm *TemplateManager) ValidateTemplate(name string) error {
	if tm.templates.Lookup(name) == nil {
		return core.NewTemplateError(
			name,
			0,
			fmt.Sprintf("template not found (available: %v)", tm.GetAvailableTemplates()),
			nil,
		)
	}
	return nil
}
