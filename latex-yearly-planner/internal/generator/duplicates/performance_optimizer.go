package generator

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
)

// PerformanceOptimizer provides comprehensive performance optimization for the system
type PerformanceOptimizer struct {
	config        *PerformanceConfig
	cache         *PerformanceCache
	pool          *WorkerPool
	profiler      *PerformanceProfiler
	logger        PDFLogger
}

// PerformanceConfig defines configuration for performance optimization
type PerformanceConfig struct {
	// Caching settings
	EnableCaching        bool          `json:"enable_caching"`
	CacheSize            int           `json:"cache_size"`
	CacheTTL             time.Duration `json:"cache_ttl"`
	CacheCleanupInterval time.Duration `json:"cache_cleanup_interval"`
	
	// Parallel processing settings
	EnableParallelProcessing bool `json:"enable_parallel_processing"`
	MaxWorkers              int  `json:"max_workers"`
	WorkerQueueSize         int  `json:"worker_queue_size"`
	
	// Memory optimization settings
	EnableMemoryOptimization bool `json:"enable_memory_optimization"`
	MaxMemoryUsage          int64 `json:"max_memory_usage_mb"`
	GCThreshold             int64 `json:"gc_threshold_mb"`
	
	// Data processing optimization
	EnableDataOptimization   bool `json:"enable_data_optimization"`
	BatchSize               int  `json:"batch_size"`
	StreamingThreshold      int  `json:"streaming_threshold"`
	
	// Layout optimization
	EnableLayoutOptimization bool `json:"enable_layout_optimization"`
	LayoutCacheSize         int  `json:"layout_cache_size"`
	SkipRedundantCalculations bool `json:"skip_redundant_calculations"`
	
	// PDF generation optimization
	EnablePDFOptimization    bool `json:"enable_pdf_optimization"`
	ParallelPDFGeneration   bool `json:"parallel_pdf_generation"`
	TemplateCaching         bool `json:"template_caching"`
	LaTeXOptimization       bool `json:"latex_optimization"`
}

// PerformanceCache provides intelligent caching for performance optimization
type PerformanceCache struct {
	config        *PerformanceConfig
	cache         map[string]*CacheEntry
	accessTimes   map[string]time.Time
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
	stopCleanup   chan bool
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	AccessCount int
	Size      int64
}

// WorkerPool manages a pool of workers for parallel processing
type WorkerPool struct {
	config     *PerformanceConfig
	workers    []*Worker
	jobQueue   chan Job
	resultChan chan JobResult
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// Worker represents a single worker in the pool
type Worker struct {
	id       int
	jobQueue chan Job
	ctx      context.Context
}

// Job represents a unit of work
type Job struct {
	ID       string
	Type     JobType
	Data     interface{}
	Priority int
	Callback func(interface{}) (interface{}, error)
}

// JobResult represents the result of a job
type JobResult struct {
	JobID    string
	Result   interface{}
	Error    error
	Duration time.Duration
}

// JobType represents the type of job
type JobType int

const (
	JobTypeDataProcessing JobType = iota
	JobTypeLayoutCalculation
	JobTypeVisualOptimization
	JobTypePDFGeneration
	JobTypeTemplateProcessing
	JobTypeCacheOperation
)

// PerformanceProfiler provides performance profiling and monitoring
type PerformanceProfiler struct {
	config     *PerformanceConfig
	metrics    *PerformanceMetrics
	startTime  time.Time
	mu         sync.RWMutex
}

// PerformanceMetrics contains performance measurements
type PerformanceMetrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	AverageResponseTime time.Duration
	MaxResponseTime     time.Duration
	MinResponseTime     time.Duration
	MemoryUsage         int64
	CacheHitRate        float64
	CacheMissRate       float64
	WorkerUtilization   float64
	Throughput          float64
	ErrorRate           float64
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer() *PerformanceOptimizer {
	config := GetDefaultPerformanceConfig()
	
	return &PerformanceOptimizer{
		config:   config,
		cache:    NewPerformanceCache(config),
		pool:     NewWorkerPool(config),
		profiler: NewPerformanceProfiler(config),
		logger:   &PerformanceOptimizerLogger{},
	}
}

// GetDefaultPerformanceConfig returns the default performance configuration
func GetDefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		EnableCaching:            true,
		CacheSize:                1000,
		CacheTTL:                 time.Hour * 24,
		CacheCleanupInterval:     time.Minute * 5,
		EnableParallelProcessing: true,
		MaxWorkers:               runtime.NumCPU(),
		WorkerQueueSize:          100,
		EnableMemoryOptimization: true,
		MaxMemoryUsage:           512, // 512MB
		GCThreshold:              256, // 256MB
		EnableDataOptimization:   true,
		BatchSize:                100,
		StreamingThreshold:       1000,
		EnableLayoutOptimization: true,
		LayoutCacheSize:         500,
		SkipRedundantCalculations: true,
		EnablePDFOptimization:    true,
		ParallelPDFGeneration:   true,
		TemplateCaching:         true,
		LaTeXOptimization:       true,
	}
}

// SetLogger sets the logger for the performance optimizer
func (po *PerformanceOptimizer) SetLogger(logger PDFLogger) {
	po.logger = logger
	po.cache.SetLogger(logger)
	po.pool.SetLogger(logger)
	po.profiler.SetLogger(logger)
}

// OptimizeDataProcessing optimizes data processing operations
func (po *PerformanceOptimizer) OptimizeDataProcessing(tasks []*data.Task) ([]*data.Task, error) {
	start := time.Now()
	po.profiler.StartOperation("data_processing")
	
	// Check if we should use streaming for large datasets
	if len(tasks) > po.config.StreamingThreshold {
		return po.optimizeDataProcessingStreaming(tasks)
	}
	
	// Use parallel processing for medium datasets
	if len(tasks) > po.config.BatchSize && po.config.EnableParallelProcessing {
		return po.optimizeDataProcessingParallel(tasks)
	}
	
	// Use optimized single-threaded processing for small datasets
	return po.optimizeDataProcessingSingle(tasks)
}

// optimizeDataProcessingStreaming optimizes data processing using streaming
func (po *PerformanceOptimizer) optimizeDataProcessingStreaming(tasks []*data.Task) ([]*data.Task, error) {
	po.logger.Info("Using streaming optimization for %d tasks", len(tasks))
	
	// Process tasks in batches
	batchSize := po.config.BatchSize
	results := make([]*data.Task, 0, len(tasks))
	
	for i := 0; i < len(tasks); i += batchSize {
		end := i + batchSize
		if end > len(tasks) {
			end = len(tasks)
		}
		
		batch := tasks[i:end]
		optimizedBatch, err := po.optimizeDataProcessingSingle(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
		}
		
		results = append(results, optimizedBatch...)
		
		// Check memory usage and trigger GC if needed
		if po.shouldTriggerGC() {
			runtime.GC()
		}
	}
	
	po.profiler.EndOperation("data_processing", time.Since(start))
	return results, nil
}

// optimizeDataProcessingParallel optimizes data processing using parallel processing
func (po *PerformanceOptimizer) optimizeDataProcessingParallel(tasks []*data.Task) ([]*data.Task, error) {
	po.logger.Info("Using parallel optimization for %d tasks", len(tasks))
	
	// Create jobs for parallel processing
	jobs := make([]Job, 0, len(tasks))
	for i, task := range tasks {
		jobs = append(jobs, Job{
			ID:       fmt.Sprintf("task_%d", i),
			Type:     JobTypeDataProcessing,
			Data:     task,
			Priority: 1,
			Callback: po.optimizeTask,
		})
	}
	
	// Process jobs in parallel
	results, err := po.pool.ProcessJobs(jobs)
	if err != nil {
		return nil, fmt.Errorf("parallel processing failed: %w", err)
	}
	
	// Convert results back to tasks
	optimizedTasks := make([]*data.Task, len(results))
	for i, result := range results {
		if result.Error != nil {
			return nil, fmt.Errorf("task %d processing failed: %w", i, result.Error)
		}
		optimizedTasks[i] = result.Result.(*data.Task)
	}
	
	po.profiler.EndOperation("data_processing", time.Since(start))
	return optimizedTasks, nil
}

// optimizeDataProcessingSingle optimizes data processing using single-threaded processing
func (po *PerformanceOptimizer) optimizeDataProcessingSingle(tasks []*data.Task) ([]*data.Task, error) {
	po.logger.Info("Using single-threaded optimization for %d tasks", len(tasks))
	
	optimizedTasks := make([]*data.Task, len(tasks))
	
	for i, task := range tasks {
		optimized, err := po.optimizeTask(task)
		if err != nil {
			return nil, fmt.Errorf("failed to optimize task %d: %w", i, err)
		}
		optimizedTasks[i] = optimized.(*data.Task)
	}
	
	po.profiler.EndOperation("data_processing", time.Since(start))
	return optimizedTasks, nil
}

// optimizeTask optimizes a single task
func (po *PerformanceOptimizer) optimizeTask(task interface{}) (interface{}, error) {
	t, ok := task.(*data.Task)
	if !ok {
		return nil, fmt.Errorf("invalid task type")
	}
	
	// Create optimized task copy
	optimized := &data.Task{
		ID:           t.ID,
		Name:         t.Name,
		StartDate:    t.StartDate,
		EndDate:      t.EndDate,
		Category:     t.Category,
		Description:  t.Description,
		Priority:     t.Priority,
		Status:       t.Status,
		Assignee:     t.Assignee,
		ParentID:     t.ParentID,
		Dependencies: make([]string, len(t.Dependencies)),
		IsMilestone:  t.IsMilestone,
	}
	
	// Copy dependencies efficiently
	copy(optimized.Dependencies, t.Dependencies)
	
	return optimized, nil
}

// OptimizeLayoutProcessing optimizes layout processing operations
func (po *PerformanceOptimizer) OptimizeLayoutProcessing(tasks []*data.Task, config *calendar.GridConfig) (*calendar.IntegratedLayoutResult, error) {
	start := time.Now()
	po.profiler.StartOperation("layout_processing")
	
	// Check cache first
	cacheKey := po.generateLayoutCacheKey(tasks, config)
	if po.config.EnableCaching {
		if cached, found := po.cache.Get(cacheKey); found {
			po.logger.Info("Layout cache hit for key: %s", cacheKey)
			po.profiler.EndOperation("layout_processing", time.Since(start))
			return cached.(*calendar.IntegratedLayoutResult), nil
		}
	}
	
	// Create optimized layout integration
	layoutIntegration := calendar.NewCalendarGridIntegration(config)
	
	// Process tasks with optimized layout
	result, err := layoutIntegration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		return nil, fmt.Errorf("layout processing failed: %w", err)
	}
	
	// Cache the result
	if po.config.EnableCaching {
		po.cache.Set(cacheKey, result)
	}
	
	po.profiler.EndOperation("layout_processing", time.Since(start))
	return result, nil
}

// OptimizePDFGeneration optimizes PDF generation operations
func (po *PerformanceOptimizer) OptimizePDFGeneration(cfg config.Config, options PDFGenerationOptions) (*PDFGenerationResult, error) {
	start := time.Now()
	po.profiler.StartOperation("pdf_generation")
	
	// Create optimized PDF pipeline
	pipeline := NewPDFPipeline("", "")
	pipeline.SetLogger(po.logger)
	
	// Apply PDF optimizations
	if po.config.LaTeXOptimization {
		options = po.optimizeLaTeXOptions(options)
	}
	
	// Generate PDF with optimizations
	result, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}
	
	po.profiler.EndOperation("pdf_generation", time.Since(start))
	return result, nil
}

// optimizeLaTeXOptions optimizes LaTeX compilation options
func (po *PerformanceOptimizer) optimizeLaTeXOptions(options PDFGenerationOptions) PDFGenerationOptions {
	// Use faster LaTeX engine if available
	if options.CompilationEngine == "" {
		options.CompilationEngine = "pdflatex" // Fastest engine
	}
	
	// Reduce retries for faster failure
	if options.MaxRetries == 0 {
		options.MaxRetries = 2 // Reduced from default 3
	}
	
	// Enable cleanup for memory efficiency
	options.CleanupTempFiles = true
	
	return options
}

// generateLayoutCacheKey generates a cache key for layout processing
func (po *PerformanceOptimizer) generateLayoutCacheKey(tasks []*data.Task, config *calendar.GridConfig) string {
	// Create a hash of task IDs and config
	var key string
	for _, task := range tasks {
		key += task.ID + "_"
	}
	key += fmt.Sprintf("%v_%v_%.2f_%.2f", 
		config.CalendarStart, config.CalendarEnd, config.DayWidth, config.DayHeight)
	return key
}

// shouldTriggerGC checks if garbage collection should be triggered
func (po *PerformanceOptimizer) shouldTriggerGC() bool {
	if !po.config.EnableMemoryOptimization {
		return false
	}
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Convert bytes to MB
	memoryUsageMB := int64(m.Alloc / 1024 / 1024)
	return memoryUsageMB > po.config.GCThreshold
}

// GetPerformanceMetrics returns current performance metrics
func (po *PerformanceOptimizer) GetPerformanceMetrics() *PerformanceMetrics {
	return po.profiler.GetMetrics()
}

// ResetMetrics resets performance metrics
func (po *PerformanceOptimizer) ResetMetrics() {
	po.profiler.ResetMetrics()
}

// NewPerformanceCache creates a new performance cache
func NewPerformanceCache(config *PerformanceConfig) *PerformanceCache {
	cache := &PerformanceCache{
		config:      config,
		cache:       make(map[string]*CacheEntry),
		accessTimes: make(map[string]time.Time),
		stopCleanup: make(chan bool),
	}
	
	// Start cleanup goroutine
	if config.EnableCaching {
		cache.startCleanup()
	}
	
	return cache
}

// SetLogger sets the logger for the cache
func (pc *PerformanceCache) SetLogger(logger PDFLogger) {
	// Cache doesn't need logging for now
}

// Get retrieves an item from the cache
func (pc *PerformanceCache) Get(key string) (interface{}, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	
	entry, exists := pc.cache[key]
	if !exists {
		return nil, false
	}
	
	// Check if entry has expired
	if time.Since(entry.CreatedAt) > pc.config.CacheTTL {
		return nil, false
	}
	
	// Update access time and count
	entry.AccessCount++
	pc.accessTimes[key] = time.Now()
	
	return entry.Data, true
}

// Set stores an item in the cache
func (pc *PerformanceCache) Set(key string, value interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	
	// Check cache size limit
	if len(pc.cache) >= pc.config.CacheSize {
		pc.evictLRU()
	}
	
	// Create cache entry
	entry := &CacheEntry{
		Data:        value,
		CreatedAt:   time.Now(),
		AccessCount: 1,
		Size:        pc.calculateSize(value),
	}
	
	pc.cache[key] = entry
	pc.accessTimes[key] = time.Now()
}

// evictLRU evicts the least recently used item
func (pc *PerformanceCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time
	
	for key, accessTime := range pc.accessTimes {
		if oldestKey == "" || accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = accessTime
		}
	}
	
	if oldestKey != "" {
		delete(pc.cache, oldestKey)
		delete(pc.accessTimes, oldestKey)
	}
}

// calculateSize calculates the approximate size of a value
func (pc *PerformanceCache) calculateSize(value interface{}) int64 {
	// Simple size estimation
	return int64(1024) // 1KB default
}

// startCleanup starts the cache cleanup goroutine
func (pc *PerformanceCache) startCleanup() {
	pc.cleanupTicker = time.NewTicker(pc.config.CacheCleanupInterval)
	
	go func() {
		for {
			select {
			case <-pc.cleanupTicker.C:
				pc.cleanup()
			case <-pc.stopCleanup:
				pc.cleanupTicker.Stop()
				return
			}
		}
	}()
}

// cleanup removes expired entries from the cache
func (pc *PerformanceCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	
	now := time.Now()
	for key, entry := range pc.cache {
		if now.Sub(entry.CreatedAt) > pc.config.CacheTTL {
			delete(pc.cache, key)
			delete(pc.accessTimes, key)
		}
	}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(config *PerformanceConfig) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	pool := &WorkerPool{
		config:     config,
		workers:    make([]*Worker, config.MaxWorkers),
		jobQueue:   make(chan Job, config.WorkerQueueSize),
		resultChan: make(chan JobResult, config.WorkerQueueSize),
		ctx:        ctx,
		cancel:     cancel,
	}
	
	// Create workers
	for i := 0; i < config.MaxWorkers; i++ {
		pool.workers[i] = &Worker{
			id:       i,
			jobQueue: pool.jobQueue,
			ctx:      ctx,
		}
	}
	
	// Start workers
	pool.startWorkers()
	
	return pool
}

// SetLogger sets the logger for the worker pool
func (wp *WorkerPool) SetLogger(logger PDFLogger) {
	// Worker pool doesn't need logging for now
}

// ProcessJobs processes a batch of jobs
func (wp *WorkerPool) ProcessJobs(jobs []Job) ([]JobResult, error) {
	// Send jobs to queue
	for _, job := range jobs {
		select {
		case wp.jobQueue <- job:
		case <-wp.ctx.Done():
			return nil, wp.ctx.Err()
		}
	}
	
	// Collect results
	results := make([]JobResult, len(jobs))
	for i := 0; i < len(jobs); i++ {
		select {
		case result := <-wp.resultChan:
			results[i] = result
		case <-wp.ctx.Done():
			return nil, wp.ctx.Err()
		}
	}
	
	return results, nil
}

// startWorkers starts all workers
func (wp *WorkerPool) startWorkers() {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		go worker.run(wp.resultChan, &wp.wg)
	}
}

// run runs a worker
func (w *Worker) run(resultChan chan<- JobResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for {
		select {
		case job := <-w.jobQueue:
			start := time.Now()
			result, err := job.Callback(job.Data)
			duration := time.Since(start)
			
			resultChan <- JobResult{
				JobID:    job.ID,
				Result:   result,
				Error:    err,
				Duration: duration,
			}
		case <-w.ctx.Done():
			return
		}
	}
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler(config *PerformanceConfig) *PerformanceProfiler {
	return &PerformanceProfiler{
		config:    config,
		metrics:   &PerformanceMetrics{},
		startTime: time.Now(),
	}
}

// SetLogger sets the logger for the profiler
func (pp *PerformanceProfiler) SetLogger(logger PDFLogger) {
	// Profiler doesn't need logging for now
}

// StartOperation starts timing an operation
func (pp *PerformanceProfiler) StartOperation(operation string) {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	pp.metrics.TotalRequests++
}

// EndOperation ends timing an operation
func (pp *PerformanceProfiler) EndOperation(operation string, duration time.Duration) {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	pp.metrics.SuccessfulRequests++
	
	// Update timing metrics
	if pp.metrics.AverageResponseTime == 0 {
		pp.metrics.AverageResponseTime = duration
	} else {
		pp.metrics.AverageResponseTime = (pp.metrics.AverageResponseTime + duration) / 2
	}
	
	if duration > pp.metrics.MaxResponseTime {
		pp.metrics.MaxResponseTime = duration
	}
	
	if pp.metrics.MinResponseTime == 0 || duration < pp.metrics.MinResponseTime {
		pp.metrics.MinResponseTime = duration
	}
	
	// Update memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	pp.metrics.MemoryUsage = int64(m.Alloc / 1024 / 1024) // MB
}

// GetMetrics returns current performance metrics
func (pp *PerformanceProfiler) GetMetrics() *PerformanceMetrics {
	pp.mu.RLock()
	defer pp.mu.RUnlock()
	
	// Calculate derived metrics
	if pp.metrics.TotalRequests > 0 {
		pp.metrics.ErrorRate = float64(pp.metrics.FailedRequests) / float64(pp.metrics.TotalRequests)
		pp.metrics.Throughput = float64(pp.metrics.SuccessfulRequests) / time.Since(pp.startTime).Seconds()
	}
	
	return pp.metrics
}

// ResetMetrics resets performance metrics
func (pp *PerformanceProfiler) ResetMetrics() {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	pp.metrics = &PerformanceMetrics{}
	pp.startTime = time.Now()
}

// PerformanceOptimizerLogger provides logging for performance optimizer
type PerformanceOptimizerLogger struct{}

func (l *PerformanceOptimizerLogger) Info(msg string, args ...interface{})  { fmt.Printf("[PERF-INFO] "+msg+"\n", args...) }
func (l *PerformanceOptimizerLogger) Error(msg string, args ...interface{}) { fmt.Printf("[PERF-ERROR] "+msg+"\n", args...) }
func (l *PerformanceOptimizerLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[PERF-DEBUG] "+msg+"\n", args...) }
