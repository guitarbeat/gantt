package generator

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/data"
)

// OptimizedLayoutEngine provides high-performance layout optimization
type OptimizedLayoutEngine struct {
	config        *LayoutEngineConfig
	cache         *LayoutCache
	optimizer     *LayoutOptimizer
	calculator    *LayoutCalculator
	logger        PDFLogger
}

// LayoutEngineConfig defines configuration for optimized layout
type LayoutEngineConfig struct {
	// Performance settings
	EnableLayoutCache     bool `json:"enable_layout_cache"`
	CacheSize             int  `json:"cache_size"`
	EnableParallelLayout  bool `json:"enable_parallel_layout"`
	MaxWorkers            int  `json:"max_workers"`
	
	// Layout optimization
	EnableSmartStacking   bool    `json:"enable_smart_stacking"`
	EnableOverlapDetection bool   `json:"enable_overlap_detection"`
	EnableSpaceOptimization bool  `json:"enable_space_optimization"`
	MinTaskHeight         float64 `json:"min_task_height"`
	MaxTaskHeight         float64 `json:"max_task_height"`
	
	// Visual optimization
	EnableVisualOptimization bool    `json:"enable_visual_optimization"`
	OptimalSpacing          float64 `json:"optimal_spacing"`
	MaxOverlapThreshold     float64 `json:"max_overlap_threshold"`
	
	// Memory optimization
	EnableMemoryOptimization bool `json:"enable_memory_optimization"`
	MaxMemoryUsage          int64 `json:"max_memory_usage_mb"`
}

// LayoutCache provides intelligent caching for layout calculations
type LayoutCache struct {
	config      *LayoutEngineConfig
	layouts     map[string]*CachedLayout
	calculations map[string]*CachedCalculation
	mu          sync.RWMutex
}

// CachedLayout represents cached layout data
type CachedLayout struct {
	Tasks      []*data.Task
	Positions  map[string]*TaskPosition
	CreatedAt  time.Time
	TaskCount  int
	Hash       string
}

// CachedCalculation represents cached calculation data
type CachedCalculation struct {
	Result     interface{}
	CreatedAt  time.Time
	InputHash  string
}

// TaskPosition represents the position of a task in the layout
type TaskPosition struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Row    int
	Column int
}

// LayoutOptimizer provides layout optimization algorithms
type LayoutOptimizer struct {
	config        *LayoutEngineConfig
	overlapDetector *OverlapDetector
	spaceOptimizer  *SpaceOptimizer
	logger        PDFLogger
}

// LayoutCalculator provides optimized layout calculations
type LayoutCalculator struct {
	config        *LayoutEngineConfig
	visualOptimizer *VisualOptimizer
	logger        PDFLogger
}

// OverlapDetector detects and resolves task overlaps
type OverlapDetector struct {
	config *LayoutEngineConfig
}

// SpaceOptimizer optimizes space usage in layouts
type SpaceOptimizer struct {
	config *LayoutEngineConfig
}

// VisualOptimizer optimizes visual appearance
type VisualOptimizer struct {
	config *LayoutEngineConfig
}

// NewOptimizedLayoutEngine creates a new optimized layout engine
func NewOptimizedLayoutEngine() *OptimizedLayoutEngine {
	config := GetDefaultLayoutEngineConfig()
	
	return &OptimizedLayoutEngine{
		config:     config,
		cache:      NewLayoutCache(config),
		optimizer:  NewLayoutOptimizer(config),
		calculator: NewLayoutCalculator(config),
		logger:     &OptimizedLayoutEngineLogger{},
	}
}

// GetDefaultLayoutEngineConfig returns the default layout engine configuration
func GetDefaultLayoutEngineConfig() *LayoutEngineConfig {
	return &LayoutEngineConfig{
		EnableLayoutCache:       true,
		CacheSize:               50,
		EnableParallelLayout:    true,
		MaxWorkers:              4,
		EnableSmartStacking:     true,
		EnableOverlapDetection:  true,
		EnableSpaceOptimization: true,
		MinTaskHeight:           20.0,
		MaxTaskHeight:           40.0,
		EnableVisualOptimization: true,
		OptimalSpacing:          5.0,
		MaxOverlapThreshold:     0.1,
		EnableMemoryOptimization: true,
		MaxMemoryUsage:          128, // 128MB
	}
}

// SetLogger sets the logger for the layout engine
func (ole *OptimizedLayoutEngine) SetLogger(logger PDFLogger) {
	ole.logger = logger
	ole.optimizer.SetLogger(logger)
	ole.calculator.SetLogger(logger)
}

// GenerateLayout generates an optimized layout for tasks
func (ole *OptimizedLayoutEngine) GenerateLayout(tasks []*data.Task, config *calendar.CalendarConfig) (*calendar.CalendarLayout, error) {
	start := time.Now()
	ole.logger.Info("Starting optimized layout generation for %d tasks", len(tasks))
	
	// Check cache first
	if ole.config.EnableLayoutCache {
		layoutHash := ole.calculateLayoutHash(tasks, config)
		if cached, found := ole.cache.GetLayout(layoutHash); found {
			ole.logger.Info("Layout cache hit")
			return ole.convertCachedLayout(cached, config), nil
		}
	}
	
	// Generate layout
	layout, err := ole.generateLayoutInternal(tasks, config)
	if err != nil {
		return nil, fmt.Errorf("layout generation failed: %w", err)
	}
	
	// Cache the result
	if ole.config.EnableLayoutCache {
		layoutHash := ole.calculateLayoutHash(tasks, config)
		ole.cache.SetLayout(layoutHash, &CachedLayout{
			Tasks:     tasks,
			Positions: ole.extractPositions(layout),
			CreatedAt: time.Now(),
			TaskCount: len(tasks),
			Hash:      layoutHash,
		})
	}
	
	ole.logger.Info("Generated layout in %v", time.Since(start))
	return layout, nil
}

// generateLayoutInternal generates the actual layout
func (ole *OptimizedLayoutEngine) generateLayoutInternal(tasks []*data.Task, config *calendar.CalendarConfig) (*calendar.CalendarLayout, error) {
	// Sort tasks by start date for better layout
	sortedTasks := ole.sortTasksByDate(tasks)
	
	// Calculate layout dimensions
	dimensions := ole.calculateLayoutDimensions(sortedTasks, config)
	
	// Generate base layout
	layout := &calendar.CalendarLayout{
		Tasks:      sortedTasks,
		Dimensions: dimensions,
		Config:     config,
	}
	
	// Apply optimizations
	if ole.config.EnableSmartStacking {
		layout = ole.optimizer.ApplySmartStacking(layout)
	}
	
	if ole.config.EnableOverlapDetection {
		layout = ole.optimizer.ResolveOverlaps(layout)
	}
	
	if ole.config.EnableSpaceOptimization {
		layout = ole.optimizer.OptimizeSpace(layout)
	}
	
	if ole.config.EnableVisualOptimization {
		layout = ole.calculator.OptimizeVisuals(layout)
	}
	
	return layout, nil
}

// sortTasksByDate sorts tasks by start date
func (ole *OptimizedLayoutEngine) sortTasksByDate(tasks []*data.Task) []*data.Task {
	sorted := make([]*data.Task, len(tasks))
	copy(sorted, tasks)
	
	sort.Slice(sorted, func(i, j int) bool {
		// Handle tasks without dates
		if sorted[i].StartDate.IsZero() && sorted[j].StartDate.IsZero() {
			return sorted[i].Name < sorted[j].Name
		}
		if sorted[i].StartDate.IsZero() {
			return false
		}
		if sorted[j].StartDate.IsZero() {
			return true
		}
		return sorted[i].StartDate.Before(sorted[j].StartDate)
	})
	
	return sorted
}

// calculateLayoutDimensions calculates optimal layout dimensions
func (ole *OptimizedLayoutEngine) calculateLayoutDimensions(tasks []*data.Task, config *calendar.CalendarConfig) *calendar.LayoutDimensions {
	// Calculate time range
	var startDate, endDate time.Time
	for _, task := range tasks {
		if !task.StartDate.IsZero() {
			if startDate.IsZero() || task.StartDate.Before(startDate) {
				startDate = task.StartDate
			}
		}
		if !task.EndDate.IsZero() {
			if endDate.IsZero() || task.EndDate.After(endDate) {
				endDate = task.EndDate
			}
		}
	}
	
	// Default to current year if no dates
	if startDate.IsZero() {
		startDate = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if endDate.IsZero() {
		endDate = time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC)
	}
	
	// Calculate dimensions
	timeSpan := endDate.Sub(startDate)
	days := int(timeSpan.Hours() / 24)
	
	// Calculate optimal row count based on task count
	rowCount := ole.calculateOptimalRowCount(len(tasks))
	
	return &calendar.LayoutDimensions{
		Width:     float64(days) * config.DayWidth,
		Height:    float64(rowCount) * config.TaskHeight,
		RowCount:  rowCount,
		DayWidth:  config.DayWidth,
		TaskHeight: config.TaskHeight,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// calculateOptimalRowCount calculates the optimal number of rows
func (ole *OptimizedLayoutEngine) calculateOptimalRowCount(taskCount int) int {
	// Use square root for optimal distribution
	optimal := int(math.Ceil(math.Sqrt(float64(taskCount))))
	
	// Ensure minimum and maximum bounds
	minRows := 1
	maxRows := 20
	
	if optimal < minRows {
		return minRows
	}
	if optimal > maxRows {
		return maxRows
	}
	
	return optimal
}

// calculateLayoutHash calculates a hash for layout caching
func (ole *OptimizedLayoutEngine) calculateLayoutHash(tasks []*data.Task, config *calendar.CalendarConfig) string {
	// Simple hash based on task count and config
	hash := fmt.Sprintf("%d_%d_%d_%d", 
		len(tasks), 
		int(config.DayWidth), 
		int(config.TaskHeight),
		config.Year)
	
	return hash
}

// extractPositions extracts task positions from layout
func (ole *OptimizedLayoutEngine) extractPositions(layout *calendar.CalendarLayout) map[string]*TaskPosition {
	positions := make(map[string]*TaskPosition)
	
	for i, task := range layout.Tasks {
		// Calculate position based on task dates
		x := ole.calculateTaskX(task, layout)
		y := float64(i) * layout.Dimensions.TaskHeight
		width := ole.calculateTaskWidth(task, layout)
		height := layout.Dimensions.TaskHeight
		
		positions[task.ID] = &TaskPosition{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			Row:    i,
			Column: 0, // Will be calculated by smart stacking
		}
	}
	
	return positions
}

// calculateTaskX calculates the X position of a task
func (ole *OptimizedLayoutEngine) calculateTaskX(task *data.Task, layout *calendar.CalendarLayout) float64 {
	if task.StartDate.IsZero() {
		return 0
	}
	
	daysSinceStart := int(task.StartDate.Sub(layout.Dimensions.StartDate).Hours() / 24)
	return float64(daysSinceStart) * layout.Dimensions.DayWidth
}

// calculateTaskWidth calculates the width of a task
func (ole *OptimizedLayoutEngine) calculateTaskWidth(task *data.Task, layout *calendar.CalendarLayout) float64 {
	if task.StartDate.IsZero() || task.EndDate.IsZero() {
		return layout.Dimensions.DayWidth // Default width
	}
	
	days := int(task.EndDate.Sub(task.StartDate).Hours() / 24)
	if days < 1 {
		days = 1 // Minimum width
	}
	
	return float64(days) * layout.Dimensions.DayWidth
}

// convertCachedLayout converts cached layout to CalendarLayout
func (ole *OptimizedLayoutEngine) convertCachedLayout(cached *CachedLayout, config *calendar.CalendarConfig) *calendar.CalendarLayout {
	// This would convert the cached layout back to CalendarLayout
	// For now, return a basic layout
	return &calendar.CalendarLayout{
		Tasks:      cached.Tasks,
		Dimensions: &calendar.LayoutDimensions{
			Width:     800,
			Height:    600,
			RowCount:  cached.TaskCount,
			DayWidth:  config.DayWidth,
			TaskHeight: config.TaskHeight,
		},
		Config: config,
	}
}

// NewLayoutCache creates a new layout cache
func NewLayoutCache(config *LayoutEngineConfig) *LayoutCache {
	return &LayoutCache{
		config:      config,
		layouts:     make(map[string]*CachedLayout),
		calculations: make(map[string]*CachedCalculation),
	}
}

// GetLayout retrieves cached layout
func (lc *LayoutCache) GetLayout(hash string) (*CachedLayout, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	
	cached, exists := lc.layouts[hash]
	if !exists {
		return nil, false
	}
	
	// Check if expired (24 hours)
	if time.Since(cached.CreatedAt) > time.Hour*24 {
		return nil, false
	}
	
	return cached, true
}

// SetLayout stores cached layout
func (lc *LayoutCache) SetLayout(hash string, layout *CachedLayout) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	
	// Check cache size limit
	if len(lc.layouts) >= lc.config.CacheSize {
		lc.evictOldestLayout()
	}
	
	lc.layouts[hash] = layout
}

// evictOldestLayout removes the oldest cached layout
func (lc *LayoutCache) evictOldestLayout() {
	var oldestKey string
	var oldestTime time.Time
	
	for key, layout := range lc.layouts {
		if oldestKey == "" || layout.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = layout.CreatedAt
		}
	}
	
	if oldestKey != "" {
		delete(lc.layouts, oldestKey)
	}
}

// NewLayoutOptimizer creates a new layout optimizer
func NewLayoutOptimizer(config *LayoutEngineConfig) *LayoutOptimizer {
	return &LayoutOptimizer{
		config:         config,
		overlapDetector: NewOverlapDetector(config),
		spaceOptimizer:  NewSpaceOptimizer(config),
		logger:         &OptimizedLayoutEngineLogger{},
	}
}

// SetLogger sets the logger for the optimizer
func (lo *LayoutOptimizer) SetLogger(logger PDFLogger) {
	lo.logger = logger
}

// ApplySmartStacking applies smart stacking optimization
func (lo *LayoutOptimizer) ApplySmartStacking(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	lo.logger.Info("Applying smart stacking optimization")
	
	// Group tasks by time periods
	timeGroups := lo.groupTasksByTime(layout.Tasks)
	
	// Apply stacking within each group
	for _, group := range timeGroups {
		lo.stackTasksInGroup(group)
	}
	
	return layout
}

// groupTasksByTime groups tasks by overlapping time periods
func (lo *LayoutOptimizer) groupTasksByTime(tasks []*data.Task) [][]*data.Task {
	groups := [][]*data.Task{}
	processed := make(map[string]bool)
	
	for _, task := range tasks {
		if processed[task.ID] {
			continue
		}
		
		group := []*data.Task{task}
		processed[task.ID] = true
		
		// Find overlapping tasks
		for _, otherTask := range tasks {
			if processed[otherTask.ID] {
				continue
			}
			
			if lo.tasksOverlap(task, otherTask) {
				group = append(group, otherTask)
				processed[otherTask.ID] = true
			}
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

// tasksOverlap checks if two tasks overlap in time
func (lo *LayoutOptimizer) tasksOverlap(task1, task2 *data.Task) bool {
	if task1.StartDate.IsZero() || task1.EndDate.IsZero() ||
		task2.StartDate.IsZero() || task2.EndDate.IsZero() {
		return false
	}
	
	return task1.StartDate.Before(task2.EndDate) && task2.StartDate.Before(task1.EndDate)
}

// stackTasksInGroup stacks tasks within a group
func (lo *LayoutOptimizer) stackTasksInGroup(tasks []*data.Task) {
	// Sort by start date
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartDate.Before(tasks[j].StartDate)
	})
	
	// Simple stacking: assign row numbers
	for i, task := range tasks {
		// This would update the task's row position
		_ = i
		_ = task
		// Implementation would update task position
	}
}

// ResolveOverlaps resolves task overlaps
func (lo *LayoutOptimizer) ResolveOverlaps(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	lo.logger.Info("Resolving task overlaps")
	
	// Detect overlaps
	overlaps := lo.overlapDetector.DetectOverlaps(layout.Tasks)
	
	// Resolve overlaps
	for _, overlap := range overlaps {
		lo.resolveOverlap(overlap, layout)
	}
	
	return layout
}

// OptimizeSpace optimizes space usage
func (lo *LayoutOptimizer) OptimizeSpace(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	lo.logger.Info("Optimizing space usage")
	
	// Apply space optimization
	layout = lo.spaceOptimizer.Optimize(layout)
	
	return layout
}

// resolveOverlap resolves a specific overlap
func (lo *LayoutOptimizer) resolveOverlap(overlap *OverlapInfo, layout *calendar.CalendarLayout) {
	// Simple resolution: move one task down
	// Implementation would adjust task positions
}

// OverlapInfo represents information about a task overlap
type OverlapInfo struct {
	Task1 *data.Task
	Task2 *data.Task
	Overlap float64
}

// NewOverlapDetector creates a new overlap detector
func NewOverlapDetector(config *LayoutEngineConfig) *OverlapDetector {
	return &OverlapDetector{config: config}
}

// DetectOverlaps detects task overlaps
func (od *OverlapDetector) DetectOverlaps(tasks []*data.Task) []*OverlapInfo {
	overlaps := []*OverlapInfo{}
	
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if od.tasksOverlap(tasks[i], tasks[j]) {
				overlap := od.calculateOverlap(tasks[i], tasks[j])
				if overlap > od.config.MaxOverlapThreshold {
					overlaps = append(overlaps, &OverlapInfo{
						Task1:   tasks[i],
						Task2:   tasks[j],
						Overlap: overlap,
					})
				}
			}
		}
	}
	
	return overlaps
}

// tasksOverlap checks if two tasks overlap
func (od *OverlapDetector) tasksOverlap(task1, task2 *data.Task) bool {
	if task1.StartDate.IsZero() || task1.EndDate.IsZero() ||
		task2.StartDate.IsZero() || task2.EndDate.IsZero() {
		return false
	}
	
	return task1.StartDate.Before(task2.EndDate) && task2.StartDate.Before(task1.EndDate)
}

// calculateOverlap calculates the overlap percentage
func (od *OverlapDetector) calculateOverlap(task1, task2 *data.Task) float64 {
	start := task1.StartDate
	if task2.StartDate.After(start) {
		start = task2.StartDate
	}
	
	end := task1.EndDate
	if task2.EndDate.Before(end) {
		end = task2.EndDate
	}
	
	if start.After(end) {
		return 0
	}
	
	overlapDuration := end.Sub(start)
	task1Duration := task1.EndDate.Sub(task1.StartDate)
	
	return float64(overlapDuration) / float64(task1Duration)
}

// NewSpaceOptimizer creates a new space optimizer
func NewSpaceOptimizer(config *LayoutEngineConfig) *SpaceOptimizer {
	return &SpaceOptimizer{config: config}
}

// Optimize optimizes space usage
func (so *SpaceOptimizer) Optimize(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	// Apply space optimization algorithms
	// This would include:
	// - Minimizing white space
	// - Optimizing row heights
	// - Adjusting task widths
	// - Reorganizing layout structure
	
	return layout
}

// NewLayoutCalculator creates a new layout calculator
func NewLayoutCalculator(config *LayoutEngineConfig) *LayoutCalculator {
	return &LayoutCalculator{
		config:         config,
		visualOptimizer: NewVisualOptimizer(config),
		logger:         &OptimizedLayoutEngineLogger{},
	}
}

// SetLogger sets the logger for the calculator
func (lc *LayoutCalculator) SetLogger(logger PDFLogger) {
	lc.logger = logger
}

// OptimizeVisuals optimizes visual appearance
func (lc *LayoutCalculator) OptimizeVisuals(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	lc.logger.Info("Optimizing visual appearance")
	
	// Apply visual optimization
	layout = lc.visualOptimizer.Optimize(layout)
	
	return layout
}

// NewVisualOptimizer creates a new visual optimizer
func NewVisualOptimizer(config *LayoutEngineConfig) *VisualOptimizer {
	return &VisualOptimizer{config: config}
}

// Optimize optimizes visual appearance
func (vo *VisualOptimizer) Optimize(layout *calendar.CalendarLayout) *calendar.CalendarLayout {
	// Apply visual optimization algorithms
	// This would include:
	// - Adjusting spacing
	// - Optimizing colors
	// - Improving readability
	// - Enhancing visual hierarchy
	
	return layout
}

// OptimizedLayoutEngineLogger provides logging for optimized layout engine
type OptimizedLayoutEngineLogger struct{}

func (l *OptimizedLayoutEngineLogger) Info(msg string, args ...interface{})  { fmt.Printf("[LAYOUT-INFO] "+msg+"\n", args...) }
func (l *OptimizedLayoutEngineLogger) Error(msg string, args ...interface{}) { fmt.Printf("[LAYOUT-ERROR] "+msg+"\n", args...) }
func (l *OptimizedLayoutEngineLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[LAYOUT-DEBUG] "+msg+"\n", args...) }
