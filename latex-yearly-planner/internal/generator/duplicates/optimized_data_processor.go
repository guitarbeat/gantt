package generator

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"latex-yearly-planner/internal/data"
)

// OptimizedDataProcessor provides high-performance data processing with caching and parallel processing
type OptimizedDataProcessor struct {
	config        *DataProcessorConfig
	cache         *DataCache
	parser        *OptimizedCSVParser
	validator     *OptimizedValidator
	logger        PDFLogger
}

// DataProcessorConfig defines configuration for optimized data processing
type DataProcessorConfig struct {
	// Parsing optimization
	EnableFastParsing     bool `json:"enable_fast_parsing"`
	SkipValidation        bool `json:"skip_validation"`
	UseStreaming          bool `json:"use_streaming"`
	BatchSize             int  `json:"batch_size"`
	
	// Caching settings
	EnableParsingCache    bool          `json:"enable_parsing_cache"`
	CacheSize             int           `json:"cache_size"`
	CacheTTL              time.Duration `json:"cache_ttl"`
	
	// Parallel processing
	EnableParallelProcessing bool `json:"enable_parallel_processing"`
	MaxWorkers              int  `json:"max_workers"`
	
	// Memory optimization
	EnableMemoryOptimization bool `json:"enable_memory_optimization"`
	MaxMemoryUsage          int64 `json:"max_memory_usage_mb"`
	
	// Date parsing optimization
	EnableDateCache        bool     `json:"enable_date_cache"`
	PreferredDateFormats   []string `json:"preferred_date_formats"`
	SkipInvalidDates       bool     `json:"skip_invalid_dates"`
}

// DataCache provides intelligent caching for data processing
type DataCache struct {
	config      *DataProcessorConfig
	parsedData  map[string]*CachedData
	dateCache   map[string]time.Time
	mu          sync.RWMutex
}

// CachedData represents cached parsed data
type CachedData struct {
	Tasks      []*data.Task
	CreatedAt  time.Time
	FileSize   int64
	FileHash   string
}

// OptimizedCSVParser provides high-performance CSV parsing
type OptimizedCSVParser struct {
	config        *DataProcessorConfig
	dateFormats   []string
	dateCache     map[string]time.Time
	mu            sync.RWMutex
}

// OptimizedValidator provides high-performance data validation
type OptimizedValidator struct {
	config        *DataProcessorConfig
	dependencyGraph map[string][]string
	mu            sync.RWMutex
}

// NewOptimizedDataProcessor creates a new optimized data processor
func NewOptimizedDataProcessor() *OptimizedDataProcessor {
	config := GetDefaultDataProcessorConfig()
	
	return &OptimizedDataProcessor{
		config:    config,
		cache:     NewDataCache(config),
		parser:    NewOptimizedCSVParser(config),
		validator: NewOptimizedValidator(config),
		logger:    &OptimizedDataProcessorLogger{},
	}
}

// GetDefaultDataProcessorConfig returns the default data processor configuration
func GetDefaultDataProcessorConfig() *DataProcessorConfig {
	return &DataProcessorConfig{
		EnableFastParsing:        true,
		SkipValidation:           false,
		UseStreaming:             true,
		BatchSize:                100,
		EnableParsingCache:       true,
		CacheSize:                100,
		CacheTTL:                 time.Hour * 24,
		EnableParallelProcessing: true,
		MaxWorkers:               4,
		EnableMemoryOptimization: true,
		MaxMemoryUsage:           256, // 256MB
		EnableDateCache:          true,
		PreferredDateFormats: []string{
			"2006-01-02",
			"01/02/2006",
			"02/01/2006",
			"2006/01/02",
			"02.01.2006",
		},
		SkipInvalidDates: true,
	}
}

// SetLogger sets the logger for the data processor
func (odp *OptimizedDataProcessor) SetLogger(logger PDFLogger) {
	odp.logger = logger
	odp.parser.SetLogger(logger)
	odp.validator.SetLogger(logger)
}

// ProcessTasks processes tasks with optimizations
func (odp *OptimizedDataProcessor) ProcessTasks(filePath string) ([]*data.Task, error) {
	start := time.Now()
	odp.logger.Info("Starting optimized task processing for: %s", filePath)
	
	// Check cache first
	if odp.config.EnableParsingCache {
		if cached, found := odp.cache.Get(filePath); found {
			odp.logger.Info("Cache hit for file: %s", filePath)
			return cached.Tasks, nil
		}
	}
	
	// Determine processing strategy based on file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	var tasks []*data.Task
	
	// Use streaming for large files
	if fileInfo.Size() > int64(odp.config.MaxMemoryUsage*1024*1024) && odp.config.UseStreaming {
		tasks, err = odp.processTasksStreaming(filePath)
	} else {
		tasks, err = odp.processTasksBatch(filePath)
	}
	
	if err != nil {
		return nil, fmt.Errorf("task processing failed: %w", err)
	}
	
	// Cache the result
	if odp.config.EnableParsingCache {
		odp.cache.Set(filePath, &CachedData{
			Tasks:     tasks,
			CreatedAt: time.Now(),
			FileSize:  fileInfo.Size(),
			FileHash:  fmt.Sprintf("%d", fileInfo.Size()), // Simple hash
		})
	}
	
	odp.logger.Info("Processed %d tasks in %v", len(tasks), time.Since(start))
	return tasks, nil
}

// processTasksStreaming processes tasks using streaming for large files
func (odp *OptimizedDataProcessor) processTasksStreaming(filePath string) ([]*data.Task, error) {
	odp.logger.Info("Using streaming processing for large file")
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	
	// Create field index map
	fieldIndex := odp.createFieldIndex(header)
	
	var tasks []*data.Task
	batch := make([]*data.Task, 0, odp.config.BatchSize)
	rowNum := 1
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record at row %d: %w", rowNum, err)
		}
		
		rowNum++
		
		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}
		
		// Parse task
		task, err := odp.parser.ParseTask(record, fieldIndex, rowNum)
		if err != nil {
			if odp.config.SkipValidation {
				odp.logger.Debug("Skipping invalid row %d: %v", rowNum, err)
				continue
			}
			return nil, fmt.Errorf("failed to parse task at row %d: %w", rowNum, err)
		}
		
		batch = append(batch, task)
		
		// Process batch when full
		if len(batch) >= odp.config.BatchSize {
			processedBatch, err := odp.processBatch(batch)
			if err != nil {
				return nil, fmt.Errorf("failed to process batch: %w", err)
			}
			tasks = append(tasks, processedBatch...)
			batch = batch[:0] // Reset batch
		}
	}
	
	// Process remaining tasks
	if len(batch) > 0 {
		processedBatch, err := odp.processBatch(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to process final batch: %w", err)
		}
		tasks = append(tasks, processedBatch...)
	}
	
	return tasks, nil
}

// processTasksBatch processes tasks using batch processing
func (odp *OptimizedDataProcessor) processTasksBatch(filePath string) ([]*data.Task, error) {
	odp.logger.Info("Using batch processing")
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	
	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	
	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows found")
	}
	
	// Create field index map
	fieldIndex := odp.createFieldIndex(records[0])
	
	// Parse all tasks
	tasks := make([]*data.Task, 0, len(records)-1)
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header
		
		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}
		
		task, err := odp.parser.ParseTask(record, fieldIndex, rowNum)
		if err != nil {
			if odp.config.SkipValidation {
				odp.logger.Debug("Skipping invalid row %d: %v", rowNum, err)
				continue
			}
			return nil, fmt.Errorf("failed to parse task at row %d: %w", rowNum, err)
		}
		
		tasks = append(tasks, task)
	}
	
	// Process all tasks
	return odp.processBatch(tasks)
}

// processBatch processes a batch of tasks
func (odp *OptimizedDataProcessor) processBatch(tasks []*data.Task) ([]*data.Task, error) {
	// Validate tasks if validation is enabled
	if !odp.config.SkipValidation {
		validatedTasks, err := odp.validator.ValidateTasks(tasks)
		if err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		tasks = validatedTasks
	}
	
	// Optimize tasks
	optimizedTasks := make([]*data.Task, len(tasks))
	for i, task := range tasks {
		optimizedTasks[i] = odp.optimizeTask(task)
	}
	
	return optimizedTasks, nil
}

// optimizeTask optimizes a single task
func (odp *OptimizedDataProcessor) optimizeTask(task *data.Task) *data.Task {
	// Create optimized copy
	optimized := &data.Task{
		ID:           task.ID,
		Name:         task.Name,
		StartDate:    task.StartDate,
		EndDate:      task.EndDate,
		Category:     task.Category,
		Description:  task.Description,
		Priority:     task.Priority,
		Status:       task.Status,
		Assignee:     task.Assignee,
		ParentID:     task.ParentID,
		Dependencies: make([]string, len(task.Dependencies)),
		IsMilestone:  task.IsMilestone,
	}
	
	// Copy dependencies efficiently
	copy(optimized.Dependencies, task.Dependencies)
	
	return optimized
}

// createFieldIndex creates an optimized field index map
func (odp *OptimizedDataProcessor) createFieldIndex(header []string) map[string]int {
	fieldIndex := make(map[string]int, len(header))
	for i, field := range header {
		normalizedField := strings.ToLower(strings.TrimSpace(field))
		fieldIndex[normalizedField] = i
	}
	return fieldIndex
}

// NewDataCache creates a new data cache
func NewDataCache(config *DataProcessorConfig) *DataCache {
	return &DataCache{
		config:     config,
		parsedData: make(map[string]*CachedData),
		dateCache:  make(map[string]time.Time),
	}
}

// Get retrieves cached data
func (dc *DataCache) Get(key string) (*CachedData, bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	
	cached, exists := dc.parsedData[key]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if time.Since(cached.CreatedAt) > dc.config.CacheTTL {
		return nil, false
	}
	
	return cached, true
}

// Set stores cached data
func (dc *DataCache) Set(key string, data *CachedData) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	
	// Check cache size limit
	if len(dc.parsedData) >= dc.config.CacheSize {
		dc.evictOldest()
	}
	
	dc.parsedData[key] = data
}

// evictOldest removes the oldest cached item
func (dc *DataCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	
	for key, data := range dc.parsedData {
		if oldestKey == "" || data.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = data.CreatedAt
		}
	}
	
	if oldestKey != "" {
		delete(dc.parsedData, oldestKey)
	}
}

// NewOptimizedCSVParser creates a new optimized CSV parser
func NewOptimizedCSVParser(config *DataProcessorConfig) *OptimizedCSVParser {
	return &OptimizedCSVParser{
		config:      config,
		dateFormats: config.PreferredDateFormats,
		dateCache:   make(map[string]time.Time),
	}
}

// SetLogger sets the logger for the parser
func (ocp *OptimizedCSVParser) SetLogger(logger PDFLogger) {
	// Parser doesn't need logging for now
}

// ParseTask parses a single CSV record into a Task
func (ocp *OptimizedCSVParser) ParseTask(record []string, fieldIndex map[string]int, rowNum int) (*data.Task, error) {
	task := &data.Task{}
	
	// Helper function to get field value safely
	getField := func(fieldName string) string {
		normalizedField := strings.ToLower(strings.TrimSpace(fieldName))
		if index, exists := fieldIndex[normalizedField]; exists && index < len(record) {
			return strings.TrimSpace(record[index])
		}
		return ""
	}
	
	// Parse required fields
	task.ID = getField("Task ID")
	if task.ID == "" {
		return nil, fmt.Errorf("missing required field: Task ID")
	}
	
	task.Name = getField("Task Name")
	task.Description = getField("Description")
	task.Category = getField("Category")
	task.Status = getField("Status")
	task.Assignee = getField("Assignee")
	task.ParentID = getField("Parent Task ID")
	
	// Parse priority
	if priorityStr := getField("Priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			task.Priority = priority
		} else {
			task.Priority = 1 // Default
		}
	} else {
		task.Priority = 1 // Default
	}
	
	// Parse dependencies
	depsStr := getField("Dependencies")
	if depsStr != "" {
		task.Dependencies = ocp.parseDependencies(depsStr)
	} else {
		task.Dependencies = []string{}
	}
	
	// Parse dates with optimization
	startDateStr := getField("Start Date")
	if startDateStr != "" {
		startDate, err := ocp.parseDate(startDateStr)
		if err != nil {
			if ocp.config.SkipInvalidDates {
				odp.logger.Debug("Skipping invalid start date: %s", startDateStr)
			} else {
				return nil, fmt.Errorf("invalid start date '%s': %w", startDateStr, err)
			}
		} else {
			task.StartDate = startDate
		}
	}
	
	endDateStr := getField("Due Date")
	if endDateStr != "" {
		endDate, err := ocp.parseDate(endDateStr)
		if err != nil {
			if ocp.config.SkipInvalidDates {
				odp.logger.Debug("Skipping invalid end date: %s", endDateStr)
			} else {
				return nil, fmt.Errorf("invalid end date '%s': %w", endDateStr, err)
			}
		} else {
			task.EndDate = endDate
		}
	}
	
	// Determine if milestone
	task.IsMilestone = ocp.isMilestoneTask(task.Name, task.Description)
	
	return task, nil
}

// parseDate parses a date string with caching
func (ocp *OptimizedCSVParser) parseDate(dateStr string) (time.Time, error) {
	// Check cache first
	if ocp.config.EnableDateCache {
		ocp.mu.RLock()
		if cached, exists := ocp.dateCache[dateStr]; exists {
			ocp.mu.RUnlock()
			return cached, nil
		}
		ocp.mu.RUnlock()
	}
	
	// Parse date
	dateStr = strings.TrimSpace(dateStr)
	var parsed time.Time
	var err error
	
	// Try preferred formats first
	for _, format := range ocp.dateFormats {
		if parsed, err = time.Parse(format, dateStr); err == nil {
			break
		}
	}
	
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date with any format")
	}
	
	// Cache the result
	if ocp.config.EnableDateCache {
		ocp.mu.Lock()
		ocp.dateCache[dateStr] = parsed
		ocp.mu.Unlock()
	}
	
	return parsed, nil
}

// parseDependencies parses comma-separated dependencies
func (ocp *OptimizedCSVParser) parseDependencies(depsStr string) []string {
	if depsStr == "" {
		return []string{}
	}
	
	deps := strings.Split(depsStr, ",")
	cleanDeps := make([]string, 0, len(deps))
	
	for _, dep := range deps {
		cleanDep := strings.TrimSpace(dep)
		if cleanDep != "" {
			cleanDeps = append(cleanDeps, cleanDep)
		}
	}
	
	return cleanDeps
}

// isMilestoneTask determines if a task is a milestone
func (ocp *OptimizedCSVParser) isMilestoneTask(name, description string) bool {
	text := strings.ToLower(name + " " + description)
	milestoneKeywords := []string{"milestone", "deadline", "due", "complete", "finish", "submit", "deliver"}
	
	for _, keyword := range milestoneKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	
	return false
}

// NewOptimizedValidator creates a new optimized validator
func NewOptimizedValidator(config *DataProcessorConfig) *OptimizedValidator {
	return &OptimizedValidator{
		config:         config,
		dependencyGraph: make(map[string][]string),
	}
}

// SetLogger sets the logger for the validator
func (ov *OptimizedValidator) SetLogger(logger PDFLogger) {
	// Validator doesn't need logging for now
}

// ValidateTasks validates a batch of tasks
func (ov *OptimizedValidator) ValidateTasks(tasks []*data.Task) ([]*data.Task, error) {
	// Build dependency graph
	ov.buildDependencyGraph(tasks)
	
	// Validate dependencies
	if err := ov.validateDependencies(tasks); err != nil {
		return nil, fmt.Errorf("dependency validation failed: %w", err)
	}
	
	// Check for circular dependencies
	if err := ov.checkCircularDependencies(tasks); err != nil {
		return nil, fmt.Errorf("circular dependency detected: %w", err)
	}
	
	return tasks, nil
}

// buildDependencyGraph builds the dependency graph
func (ov *OptimizedValidator) buildDependencyGraph(tasks []*data.Task) {
	ov.mu.Lock()
	defer ov.mu.Unlock()
	
	ov.dependencyGraph = make(map[string][]string)
	
	for _, task := range tasks {
		ov.dependencyGraph[task.ID] = task.Dependencies
	}
}

// validateDependencies validates that all dependencies exist
func (ov *OptimizedValidator) validateDependencies(tasks []*data.Task) error {
	// Create set of existing task IDs
	taskIDs := make(map[string]bool)
	for _, task := range tasks {
		taskIDs[task.ID] = true
	}
	
	// Check each task's dependencies
	for _, task := range tasks {
		for _, depID := range task.Dependencies {
			if !taskIDs[depID] {
				return fmt.Errorf("task %s references non-existent dependency: %s", task.ID, depID)
			}
		}
	}
	
	return nil
}

// checkCircularDependencies checks for circular dependencies
func (ov *OptimizedValidator) checkCircularDependencies(tasks []*data.Task) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		
		for _, neighbor := range ov.dependencyGraph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true // Circular dependency found
			}
		}
		
		recStack[node] = false
		return false
	}
	
	// Check each unvisited node
	for taskID := range ov.dependencyGraph {
		if !visited[taskID] {
			if dfs(taskID) {
				return fmt.Errorf("circular dependency detected involving task: %s", taskID)
			}
		}
	}
	
	return nil
}

// OptimizedDataProcessorLogger provides logging for optimized data processor
type OptimizedDataProcessorLogger struct{}

func (l *OptimizedDataProcessorLogger) Info(msg string, args ...interface{})  { fmt.Printf("[DATA-INFO] "+msg+"\n", args...) }
func (l *OptimizedDataProcessorLogger) Error(msg string, args ...interface{}) { fmt.Printf("[DATA-ERROR] "+msg+"\n", args...) }
func (l *OptimizedDataProcessorLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[DATA-DEBUG] "+msg+"\n", args...) }
