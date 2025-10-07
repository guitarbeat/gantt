// Package services provides service classes for the document generation system.
//
// ConfigLoader handles all configuration-related operations including auto-detection
// of CSV files, configuration files, and loading configuration with proper validation.
package services

import (
	"fmt"
	"os"
	"strings"

	"phd-dissertation-planner/src/app/helpers"
	"phd-dissertation-planner/src/core"

	"github.com/urfave/cli/v2"
)

// ConfigLoader handles configuration loading and auto-detection
type ConfigLoader struct {
	logger core.Logger
}

// ConfigLoadResult represents the result of configuration loading
type ConfigLoadResult struct {
	Config      core.Config
	PathConfigs []string
	CSVPath     string
	Reason      string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		logger: *core.NewDefaultLogger(),
	}
}

// LoadConfiguration loads and validates configuration from CLI context
func (cl *ConfigLoader) LoadConfiguration(c *cli.Context) (*ConfigLoadResult, error) {
	initialPathConfigs := strings.Split(c.Path("config"), ",")

	// Auto-detect CSV and adjust configuration accordingly
	csvPath := c.String("PLANNER_CSV_FILE")
	if csvPath == "" {
		autoPath, err := helpers.AutoDetectCSV()
		if err == nil {
			csvPath = autoPath
			// Set the CSV path for later use
			os.Setenv("PLANNER_CSV_FILE", csvPath)
			cl.logger.Info("Auto-detected CSV file: %s", csvPath)
		}
	}

	// Auto-detect configuration based on CSV
	pathConfigs := initialPathConfigs
	if csvPath != "" && len(initialPathConfigs) == 1 && initialPathConfigs[0] == "src/core/base.yaml" {
		configResult, err := helpers.AutoDetectConfig(csvPath)
		if err == nil && len(configResult.AdditionalConfigs) > 0 {
			pathConfigs = configResult.GetConfigPaths()
			cl.logger.Info("Auto-detected configuration files: %v", pathConfigs)
		}
	}

	cfg, err := core.NewConfig(pathConfigs...)
	if err != nil {
		return nil, core.NewConfigError(
			strings.Join(pathConfigs, ","),
			"",
			"failed to load configuration",
			err,
		)
	}

	// Override output directory from CLI flag if provided
	if od := strings.TrimSpace(c.Path("outdir")); od != "" {
		cfg.OutputDir = od
	}

	return &ConfigLoadResult{
		Config:      cfg,
		PathConfigs: pathConfigs,
		CSVPath:     csvPath,
		Reason:      "configuration loaded successfully",
	}, nil
}

// ValidateConfiguration validates the loaded configuration
func (cl *ConfigLoader) ValidateConfiguration(cfg core.Config) error {
	// Validate required fields
	if cfg.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	// Validate year range
	if cfg.StartYear > cfg.EndYear {
		return fmt.Errorf("start year (%d) cannot be greater than end year (%d)", cfg.StartYear, cfg.EndYear)
	}

	// Validate pages configuration
	if len(cfg.Pages) == 0 {
		return fmt.Errorf("at least one page must be configured")
	}

	// Validate each page
	for i, page := range cfg.Pages {
		if page.Name == "" {
			return fmt.Errorf("page %d must have a name", i)
		}
		if len(page.RenderBlocks) == 0 {
			return fmt.Errorf("page %d (%s) must have at least one render block", i, page.Name)
		}
	}

	return nil
}

// GetConfigurationSummary returns a summary of the loaded configuration
func (cl *ConfigLoader) GetConfigurationSummary(result *ConfigLoadResult) string {
	var summary strings.Builder
	
	summary.WriteString("Configuration Summary:\n")
	summary.WriteString(fmt.Sprintf("  Output Directory: %s\n", result.Config.OutputDir))
	summary.WriteString(fmt.Sprintf("  Year Range: %d-%d\n", result.Config.StartYear, result.Config.EndYear))
	summary.WriteString(fmt.Sprintf("  Pages: %d\n", len(result.Config.Pages)))
	
	if result.CSVPath != "" {
		summary.WriteString(fmt.Sprintf("  CSV File: %s\n", result.CSVPath))
	}
	
	summary.WriteString(fmt.Sprintf("  Config Files: %s\n", strings.Join(result.PathConfigs, ", ")))
	summary.WriteString(fmt.Sprintf("  Reason: %s\n", result.Reason))
	
	return summary.String()
}

// ReloadConfiguration reloads configuration from the same paths
func (cl *ConfigLoader) ReloadConfiguration(pathConfigs []string) (*ConfigLoadResult, error) {
	cfg, err := core.NewConfig(pathConfigs...)
	if err != nil {
		return nil, core.NewConfigError(
			strings.Join(pathConfigs, ","),
			"",
			"failed to reload configuration",
			err,
		)
	}

	return &ConfigLoadResult{
		Config:      cfg,
		PathConfigs: pathConfigs,
		Reason:      "configuration reloaded successfully",
	}, nil
}
