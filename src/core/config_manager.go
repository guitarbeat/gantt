// Package core provides centralized configuration management.
//
// This package contains:
//   - ConfigManager: Centralized configuration loading, validation, and management
//   - Environment variable consolidation and validation
//   - Configuration hot-reloading for development
//   - Preset system integration
//   - Comprehensive error reporting and validation
//
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-yaml"
)

// ConfigManager handles centralized configuration management with validation and hot-reloading
type ConfigManager struct {
	config     Config
	configPath []string
	presetPath string
	logger     *Logger

	// Hot-reloading
	watcher       *fsnotify.Watcher
	reloadChan    chan struct{}
	stopChan      chan struct{}
	reloadMutex   sync.RWMutex
	isReloading   bool

	// Environment variables registry
	envVars map[string]EnvVarDefinition
}

// EnvVarDefinition defines an environment variable with validation rules
type EnvVarDefinition struct {
	Key          string
	Description  string
	DefaultValue string
	Required     bool
	Validator    func(string) error
}

// ConfigReloadEvent represents a configuration reload event
type ConfigReloadEvent struct {
	Timestamp time.Time
	Success   bool
	Error     error
	Config    *Config
	Reason    string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{
		config:    Config{},
		logger:    NewDefaultLogger(),
		reloadChan: make(chan struct{}, 1),
		stopChan:   make(chan struct{}),
		envVars:   make(map[string]EnvVarDefinition),
	}

	cm.registerEnvironmentVariables()
	return cm
}

// Load loads configuration from files and environment variables with validation
func (cm *ConfigManager) Load(paths []string, preset string) (*Config, error) {
	cm.configPath = paths
	cm.presetPath = preset

	// Load preset if specified
	if preset != "" {
		presetPaths, err := cm.loadPresetConfig(preset)
		if err != nil {
			return nil, fmt.Errorf("failed to load preset '%s': %w", preset, err)
		}
		paths = append(presetPaths, paths...)
	}

	// Load configuration from files
	config, err := cm.loadFromFiles(paths)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration files: %w", err)
	}

	// Apply environment variables
	if err := cm.applyEnvironmentVariables(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment variables: %w", err)
	}

	// Validate configuration
	if err := cm.validateConfiguration(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply defaults and post-processing
	cm.applyDefaults(config)

	cm.config = *config
	return config, nil
}

// ValidateAtStartup performs comprehensive validation at application startup
func (cm *ConfigManager) ValidateAtStartup(config *Config) error {
	cm.logger.Info("Performing startup configuration validation...")

	// Validate required fields
	if err := cm.validateRequiredFields(config); err != nil {
		return fmt.Errorf("required field validation failed: %w", err)
	}

	// Validate configuration structure and values
	if err := cm.validateConfigStructure(config); err != nil {
		return fmt.Errorf("structure validation failed: %w", err)
	}

	// Validate environment variables
	if err := cm.validateEnvironmentVariables(); err != nil {
		return fmt.Errorf("environment variable validation failed: %w", err)
	}

	// Validate file paths and permissions
	if err := cm.validateFilePaths(config); err != nil {
		return fmt.Errorf("file path validation failed: %w", err)
	}

	// Validate CSV data if specified
	if config.CSVFilePath != "" {
		if err := cm.validateCSVData(config); err != nil {
			return fmt.Errorf("CSV data validation failed: %w", err)
		}
	}

	cm.logger.Info("✅ Startup validation completed successfully")
	return nil
}

// StartHotReload starts configuration file watching for hot-reloading in development
func (cm *ConfigManager) StartHotReload(callback func(*ConfigReloadEvent)) error {
	if cm.watcher != nil {
		return fmt.Errorf("hot-reload already started")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	cm.watcher = watcher

	// Watch configuration files
	for _, path := range cm.configPath {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return fmt.Errorf("failed to watch config file '%s': %w", path, err)
		}
	}

	// Watch preset file if specified
	if cm.presetPath != "" {
		presetFullPath := filepath.Join("src", "core", "presets", cm.presetPath+".yaml")
		if err := watcher.Add(presetFullPath); err != nil {
			cm.logger.Debug("Preset file not found for watching: %s", presetFullPath)
		}
	}

	// Start watcher goroutine
	go cm.watchFiles(callback)

	cm.logger.Info("🔄 Hot-reload enabled for configuration files")
	return nil
}

// StopHotReload stops the hot-reloading watcher
func (cm *ConfigManager) StopHotReload() {
	if cm.watcher != nil {
		cm.stopChan <- struct{}{}
		cm.watcher.Close()
		cm.watcher = nil
		cm.logger.Info("🔄 Hot-reload stopped")
	}
}

// GetCurrentConfig returns a copy of the current configuration
func (cm *ConfigManager) GetCurrentConfig() Config {
	cm.reloadMutex.RLock()
	defer cm.reloadMutex.RUnlock()
	return cm.config
}

// registerEnvironmentVariables registers all supported environment variables
func (cm *ConfigManager) registerEnvironmentVariables() {
	cm.envVars = map[string]EnvVarDefinition{
		"PLANNER_YEAR": {
			Key:         "PLANNER_YEAR",
			Description: "Academic year for planner generation",
			Validator:   validateYear,
		},
		"PLANNER_START_YEAR": {
			Key:         "PLANNER_START_YEAR",
			Description: "Start year for multi-year planning",
			Validator:   validateYear,
		},
		"PLANNER_END_YEAR": {
			Key:         "PLANNER_END_YEAR",
			Description: "End year for multi-year planning",
			Validator:   validateYear,
		},
		"PLANNER_CSV_FILE": {
			Key:         "PLANNER_CSV_FILE",
			Description: "Path to CSV file containing task data",
			Validator:   validateFilePath,
		},
		"PLANNER_OUTPUT_DIR": {
			Key:         "PLANNER_OUTPUT_DIR",
			Description: "Output directory for generated files",
			DefaultValue: "build",
			Validator:   validateOutputDir,
		},
		"PLANNER_LAYOUT_PAPER_WIDTH": {
			Key:         "PLANNER_LAYOUT_PAPER_WIDTH",
			Description: "Paper width for PDF generation",
			Validator:   validateDimension,
		},
		"PLANNER_LAYOUT_PAPER_HEIGHT": {
			Key:         "PLANNER_LAYOUT_PAPER_HEIGHT",
			Description: "Paper height for PDF generation",
			Validator:   validateDimension,
		},
		"PLANNER_LAYOUT_PAPER_MARGIN_TOP": {
			Key:         "PLANNER_LAYOUT_PAPER_MARGIN_TOP",
			Description: "Top margin for PDF layout",
			Validator:   validateDimension,
		},
		"PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM": {
			Key:         "PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM",
			Description: "Bottom margin for PDF layout",
			Validator:   validateDimension,
		},
		"PLANNER_LAYOUT_PAPER_MARGIN_LEFT": {
			Key:         "PLANNER_LAYOUT_PAPER_MARGIN_LEFT",
			Description: "Left margin for PDF layout",
			Validator:   validateDimension,
		},
		"PLANNER_LAYOUT_PAPER_MARGIN_RIGHT": {
			Key:         "PLANNER_LAYOUT_PAPER_MARGIN_RIGHT",
			Description: "Right margin for PDF layout",
			Validator:   validateDimension,
		},
		"PLANNER_PRESET": {
			Key:         "PLANNER_PRESET",
			Description: "Configuration preset to use",
			Validator:   validatePreset,
		},
		"PLANNER_SILENT": {
			Key:         "PLANNER_SILENT",
			Description: "Suppress log output (true/false)",
			Validator:   validateBoolean,
		},
		"PLANNER_LOG_LEVEL": {
			Key:         "PLANNER_LOG_LEVEL",
			Description: "Logging level (silent/info/debug)",
			DefaultValue: "info",
			Validator:   validateLogLevel,
		},
		"DEV_TEMPLATES": {
			Key:         "DEV_TEMPLATES",
			Description: "Use filesystem templates instead of embedded (development)",
			Validator:   validateBoolean,
		},
	}
}

// loadFromFiles loads configuration from YAML files
func (cm *ConfigManager) loadFromFiles(paths []string) (*Config, error) {
	config := DefaultConfig()

	for _, path := range paths {
		// Skip missing files
		if _, err := os.Stat(path); os.IsNotExist(err) {
			cm.logger.Debug("Configuration file not found, skipping: %s", path)
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
		}

		// Skip empty files
		if len(strings.TrimSpace(string(content))) == 0 {
			continue
		}

		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML in '%s': %w", path, err)
		}

		cm.logger.Debug("Loaded configuration from: %s", path)
	}

	return &config, nil
}

// applyEnvironmentVariables applies environment variable overrides
func (cm *ConfigManager) applyEnvironmentVariables(config *Config) error {
	if err := env.Parse(config); err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	cm.logger.Debug("Applied environment variable overrides")
	return nil
}

// validateConfiguration performs comprehensive configuration validation
func (cm *ConfigManager) validateConfiguration(config *Config) error {
	validator := NewConfigValidator()
	result, err := validator.ValidateConfigFileContent(config)
	if err != nil {
		return err
	}

	if !result.IsValid {
		var errors []string
		for _, issue := range result.Errors {
			errors = append(errors, issue.Error())
		}
		return fmt.Errorf("configuration validation failed:\n%s", strings.Join(errors, "\n"))
	}

	if len(result.Warnings) > 0 {
		cm.logger.Info("Configuration warnings:")
		for _, warning := range result.Warnings {
			cm.logger.Info("  ⚠️  %s", warning.Message)
		}
	}

	return nil
}

// applyDefaults applies default values and post-processing
func (cm *ConfigManager) applyDefaults(config *Config) {
	// Apply defaults for unset values
	if config.Year == 0 {
		config.Year = time.Now().Year()
	}

	if strings.TrimSpace(config.OutputDir) == "" {
		config.OutputDir = Defaults.DefaultOutputDir
	}

	// Apply layout engine defaults
	config.setLayoutEngineDefaults()

	// Set algorithmic colors
	config.setAlgorithmicColors()

	// Set date range from CSV if available
	if config.CSVFilePath != "" {
		if err := config.setDateRangeFromCSV(); err != nil {
			cm.logger.Warn("Failed to set date range from CSV: %v", err)
		}
	}
}

// loadPresetConfig loads preset configuration
func (cm *ConfigManager) loadPresetConfig(preset string) ([]string, error) {
	presetFile := filepath.Join("src", "core", "presets", preset+".yaml")

	// Check if preset exists
	if _, err := os.Stat(presetFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("preset '%s' not found", preset)
	}

	cm.logger.Info("Loading preset: %s", preset)
	return []string{presetFile}, nil
}

// watchFiles monitors configuration files for changes
func (cm *ConfigManager) watchFiles(callback func(*ConfigReloadEvent)) {
	for {
		select {
		case event, ok := <-cm.watcher.Events:
			if !ok {
				return
			}

			// Only reload on write events
			if event.Has(fsnotify.Write) {
				cm.logger.Info("Configuration file changed: %s", event.Name)
				go cm.handleReload(callback, "file_changed")
			}

		case err, ok := <-cm.watcher.Errors:
			if !ok {
				return
			}
			cm.logger.Error("File watcher error: %v", err)

		case <-cm.stopChan:
			return
		}
	}
}

// handleReload performs the configuration reload
func (cm *ConfigManager) handleReload(callback func(*ConfigReloadEvent), reason string) {
	cm.reloadMutex.Lock()
	if cm.isReloading {
		cm.reloadMutex.Unlock()
		return // Already reloading
	}
	cm.isReloading = true
	cm.reloadMutex.Unlock()

	defer func() {
		cm.reloadMutex.Lock()
		cm.isReloading = false
		cm.reloadMutex.Unlock()
	}()

	event := &ConfigReloadEvent{
		Timestamp: time.Now(),
		Reason:    reason,
	}

	// Reload configuration
	newConfig, err := cm.Load(cm.configPath, cm.presetPath)
	if err != nil {
		event.Success = false
		event.Error = err
		cm.logger.Error("Configuration reload failed: %v", err)
	} else {
		event.Success = true
		event.Config = newConfig

		// Update current config
		cm.reloadMutex.Lock()
		cm.config = *newConfig
		cm.reloadMutex.Unlock()

		cm.logger.Info("✅ Configuration reloaded successfully")
	}

	if callback != nil {
		callback(event)
	}
}

// validateRequiredFields validates that all required configuration fields are present
// validateConfigStructure validates the overall structure and values of configuration
func (cm *ConfigManager) validateConfigStructure(config *Config) error {
	validator := NewConfigValidator()

	// Convert config to ValidationResult by validating it
	issues := validator.validateConfigStructure(*config)
	if len(issues) > 0 {
		var errors []string
		for _, issue := range issues {
			errors = append(errors, issue.Error())
		}
		return fmt.Errorf("configuration structure validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func (cm *ConfigManager) validateRequiredFields(config *Config) error {
	var errors []string

	// Check required layout fields
	if config.Layout.Paper.Width == "" {
		errors = append(errors, "layout.paper.width is required")
	}
	if config.Layout.Paper.Height == "" {
		errors = append(errors, "layout.paper.height is required")
	}

	// Check pages configuration
	if len(config.Pages) == 0 {
		errors = append(errors, "at least one page must be defined")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

// validateFilePaths validates file paths and permissions
func (cm *ConfigManager) validateFilePaths(config *Config) error {
	var errors []string

	// Validate output directory
	if config.OutputDir != "" {
		if strings.Contains(config.OutputDir, "..") {
			errors = append(errors, "output directory path cannot contain '..' for security")
		}

		// Try to create output directory to test permissions
		if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
			errors = append(errors, fmt.Sprintf("cannot create output directory '%s': %v", config.OutputDir, err))
		}
	}

	// Validate CSV file path if specified
	if config.CSVFilePath != "" {
		if _, err := os.Stat(config.CSVFilePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("CSV file does not exist: %s", config.CSVFilePath))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

// validateCSVData validates CSV data if specified
func (cm *ConfigManager) validateCSVData(config *Config) error {
	if config.CSVFilePath == "" {
		return nil
	}

	validator := NewCSVValidator()
	result, err := validator.ValidateCSVFile(config.CSVFilePath)
	if err != nil {
		return err
	}

	if !result.IsValid {
		var errors []string
		for _, issue := range result.Errors {
			errors = append(errors, issue.Error())
		}
		return fmt.Errorf("CSV validation failed:\n%s", strings.Join(errors, "\n"))
	}

	cm.logger.Info("✅ CSV data validation passed (%d rows)", result.RowCount)
	return nil
}

// validateEnvironmentVariables validates environment variable values
func (cm *ConfigManager) validateEnvironmentVariables() error {
	var errors []string

	for _, envVar := range cm.envVars {
		value := os.Getenv(envVar.Key)
		if value == "" && envVar.Required {
			errors = append(errors, fmt.Sprintf("required environment variable '%s' is not set", envVar.Key))
			continue
		}

		if value != "" && envVar.Validator != nil {
			if err := envVar.Validator(value); err != nil {
				errors = append(errors, fmt.Sprintf("environment variable '%s': %v", envVar.Key, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

// Environment variable validators

func validateYear(value string) error {
	// Basic year validation (4 digits, reasonable range)
	if len(value) != 4 {
		return fmt.Errorf("year must be 4 digits")
	}
	year := 0
	if _, err := fmt.Sscanf(value, "%d", &year); err != nil {
		return fmt.Errorf("invalid year format")
	}
	if year < 2000 || year > 2100 {
		return fmt.Errorf("year must be between 2000 and 2100")
	}
	return nil
}

func validateFilePath(value string) error {
	if strings.Contains(value, "..") {
		return fmt.Errorf("file path cannot contain '..' for security")
	}
	return nil
}

func validateOutputDir(value string) error {
	return validateFilePath(value)
}

func validateDimension(value string) error {
	// Basic dimension validation (should end with common units)
	validUnits := []string{"pt", "mm", "cm", "in", "em", "ex"}
	valid := false
	for _, unit := range validUnits {
		if strings.HasSuffix(value, unit) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("dimension must end with one of: %s", strings.Join(validUnits, ", "))
	}
	return nil
}

func validatePreset(value string) error {
	validPresets := []string{"academic", "compact", "presentation"}
	for _, preset := range validPresets {
		if value == preset {
			return nil
		}
	}
	return fmt.Errorf("preset must be one of: %s", strings.Join(validPresets, ", "))
}

func validateBoolean(value string) error {
	if value != "true" && value != "false" {
		return fmt.Errorf("must be 'true' or 'false'")
	}
	return nil
}

func validateLogLevel(value string) error {
	validLevels := []string{"silent", "info", "debug"}
	for _, level := range validLevels {
		if value == level {
			return nil
		}
	}
	return fmt.Errorf("log level must be one of: %s", strings.Join(validLevels, ", "))
}
