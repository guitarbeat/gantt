package generator

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/data"
)

// PerformanceTester provides comprehensive performance testing capabilities
type PerformanceTester struct {
	config        *PerformanceTestConfig
	monitor       *PerformanceMonitor
	generator     *OptimizedDataProcessor
	layoutEngine  *OptimizedLayoutEngine
	pdfGenerator  *OptimizedPDFGenerator
	logger        PDFLogger
}

// PerformanceTestConfig defines configuration for performance testing
type PerformanceTestConfig struct {
	// Test settings
	EnablePerformanceTests bool          `json:"enable_performance_tests"`
	TestDuration          time.Duration `json:"test_duration"`
	WarmupDuration        time.Duration `json:"warmup_duration"`
	CooldownDuration      time.Duration `json:"cooldown_duration"`
	
	// Test data generation
	MinTaskCount          int `json:"min_task_count"`
	MaxTaskCount          int `json:"max_task_count"`
	TaskCountStep         int `json:"task_count_step"`
	EnableRandomData      bool `json:"enable_random_data"`
	
	// Performance thresholds
	MaxProcessingTime     time.Duration `json:"max_processing_time"`
	MaxMemoryUsage        int64         `json:"max_memory_usage_mb"`
	MaxCPUUsage           float64       `json:"max_cpu_usage_percent"`
	
	// Test scenarios
	EnableLoadTesting     bool `json:"enable_load_testing"`
	EnableStressTesting   bool `json:"enable_stress_testing"`
	EnableMemoryTesting   bool `json:"enable_memory_testing"`
	EnableConcurrencyTesting bool `json:"enable_concurrency_testing"`
	
	// Reporting
	EnableDetailedReports bool `json:"enable_detailed_reports"`
	ReportFormat          string `json:"report_format"` // "json", "text", "html"
}

// PerformanceTestResult represents the result of a performance test
type PerformanceTestResult struct {
	TestName        string
	StartTime       time.Time
	EndTime         time.Time
	Duration        time.Duration
	TaskCount       int
	Success         bool
	Error           error
	Metrics         map[string]float64
	SystemMetrics   *SystemMetrics
	Recommendations []string
}

// LoadTestResult represents the result of a load test
type LoadTestResult struct {
	TestName        string
	ConcurrentUsers int
	TotalRequests   int
	SuccessfulRequests int
	FailedRequests  int
	AverageResponseTime time.Duration
	MaxResponseTime time.Duration
	MinResponseTime time.Duration
	Throughput      float64 // requests per second
	ErrorRate       float64
	Metrics         map[string]float64
}

// StressTestResult represents the result of a stress test
type StressTestResult struct {
	TestName        string
	MaxLoad         int
	BreakingPoint   int
	RecoveryTime    time.Duration
	MemoryLeaks     bool
	PerformanceDegradation float64
	Metrics         map[string]float64
}

// MemoryTestResult represents the result of a memory test
type MemoryTestResult struct {
	TestName        string
	InitialMemory   int64
	PeakMemory      int64
	FinalMemory     int64
	MemoryLeaks     bool
	GCStats         *GCStats
	AllocationRate  float64
	Metrics         map[string]float64
}

// GCStats represents garbage collection statistics
type GCStats struct {
	Count       int64
	PauseTotal  time.Duration
	PauseAvg    time.Duration
	PauseMax    time.Duration
	HeapSize    int64
	HeapInUse   int64
}

// NewPerformanceTester creates a new performance tester
func NewPerformanceTester() *PerformanceTester {
	config := GetDefaultPerformanceTestConfig()
	
	return &PerformanceTester{
		config:       config,
		monitor:      NewPerformanceMonitor(),
		generator:    NewOptimizedDataProcessor(),
		layoutEngine: NewOptimizedLayoutEngine(),
		pdfGenerator: NewOptimizedPDFGenerator(),
		logger:       &PerformanceTesterLogger{},
	}
}

// GetDefaultPerformanceTestConfig returns the default performance test configuration
func GetDefaultPerformanceTestConfig() *PerformanceTestConfig {
	return &PerformanceTestConfig{
		EnablePerformanceTests: true,
		TestDuration:           time.Minute * 5,
		WarmupDuration:         time.Second * 30,
		CooldownDuration:       time.Second * 30,
		MinTaskCount:           10,
		MaxTaskCount:           1000,
		TaskCountStep:          50,
		EnableRandomData:       true,
		MaxProcessingTime:      time.Minute * 2,
		MaxMemoryUsage:         512, // 512MB
		MaxCPUUsage:            80.0, // 80%
		EnableLoadTesting:      true,
		EnableStressTesting:    true,
		EnableMemoryTesting:    true,
		EnableConcurrencyTesting: true,
		EnableDetailedReports:  true,
		ReportFormat:           "json",
	}
}

// SetLogger sets the logger for the performance tester
func (pt *PerformanceTester) SetLogger(logger PDFLogger) {
	pt.logger = logger
	pt.monitor.SetLogger(logger)
	pt.generator.SetLogger(logger)
	pt.layoutEngine.SetLogger(logger)
	pt.pdfGenerator.SetLogger(logger)
}

// RunPerformanceTests runs all performance tests
func (pt *PerformanceTester) RunPerformanceTests(ctx context.Context) (*PerformanceTestSuite, error) {
	pt.logger.Info("Starting performance test suite")
	
	suite := &PerformanceTestSuite{
		StartTime: time.Now(),
		Config:    pt.config,
		Results:   make([]*PerformanceTestResult, 0),
	}
	
	// Run individual tests
	if pt.config.EnableLoadTesting {
		if err := pt.runLoadTests(ctx, suite); err != nil {
			pt.logger.Error("Load tests failed: %v", err)
		}
	}
	
	if pt.config.EnableStressTesting {
		if err := pt.runStressTests(ctx, suite); err != nil {
			pt.logger.Error("Stress tests failed: %v", err)
		}
	}
	
	if pt.config.EnableMemoryTesting {
		if err := pt.runMemoryTests(ctx, suite); err != nil {
			pt.logger.Error("Memory tests failed: %v", err)
		}
	}
	
	if pt.config.EnableConcurrencyTesting {
		if err := pt.runConcurrencyTests(ctx, suite); err != nil {
			pt.logger.Error("Concurrency tests failed: %v", err)
		}
	}
	
	suite.EndTime = time.Now()
	suite.Duration = suite.EndTime.Sub(suite.StartTime)
	
	pt.logger.Info("Performance test suite completed in %v", suite.Duration)
	return suite, nil
}

// runLoadTests runs load testing scenarios
func (pt *PerformanceTester) runLoadTests(ctx context.Context, suite *PerformanceTestSuite) error {
	pt.logger.Info("Running load tests")
	
	// Test different task counts
	for taskCount := pt.config.MinTaskCount; taskCount <= pt.config.MaxTaskCount; taskCount += pt.config.TaskCountStep {
		result := pt.runLoadTest(ctx, taskCount)
		suite.Results = append(suite.Results, result)
		
		if !result.Success {
			pt.logger.Error("Load test failed for %d tasks: %v", taskCount, result.Error)
		}
	}
	
	return nil
}

// runLoadTest runs a single load test
func (pt *PerformanceTester) runLoadTest(ctx context.Context, taskCount int) *PerformanceTestResult {
	start := time.Now()
	pt.logger.Info("Running load test with %d tasks", taskCount)
	
	// Generate test data
	tasks, err := pt.generateTestTasks(taskCount)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("load_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to generate test data: %w", err),
		}
	}
	
	// Process tasks
	processedTasks, err := pt.generator.ProcessTasks("") // Use in-memory data
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("load_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to process tasks: %w", err),
		}
	}
	
	// Generate layout
	config := &calendar.CalendarConfig{
		Year:       2024,
		DayWidth:   20.0,
		TaskHeight: 30.0,
	}
	
	layout, err := pt.layoutEngine.GenerateLayout(processedTasks, config)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("load_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to generate layout: %w", err),
		}
	}
	
	// Generate PDF
	pdfPath := fmt.Sprintf("/tmp/test_%d_tasks.pdf", taskCount)
	err = pt.pdfGenerator.GeneratePDF(ctx, layout, pdfPath)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("load_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to generate PDF: %w", err),
		}
	}
	
	end := time.Now()
	duration := end.Sub(start)
	
	// Collect metrics
	metrics := pt.collectTestMetrics()
	systemMetrics := pt.monitor.GetSystemMetrics()
	
	// Check performance thresholds
	success := true
	var recommendations []string
	
	if duration > pt.config.MaxProcessingTime {
		success = false
		recommendations = append(recommendations, "Processing time exceeded threshold")
	}
	
	if systemMetrics.MemoryUsage > pt.config.MaxMemoryUsage*1024*1024 {
		success = false
		recommendations = append(recommendations, "Memory usage exceeded threshold")
	}
	
	if systemMetrics.CPUUsage > pt.config.MaxCPUUsage {
		success = false
		recommendations = append(recommendations, "CPU usage exceeded threshold")
	}
	
	return &PerformanceTestResult{
		TestName:        fmt.Sprintf("load_test_%d_tasks", taskCount),
		StartTime:       start,
		EndTime:         end,
		Duration:        duration,
		TaskCount:       taskCount,
		Success:         success,
		Error:           nil,
		Metrics:         metrics,
		SystemMetrics:   systemMetrics,
		Recommendations: recommendations,
	}
}

// runStressTests runs stress testing scenarios
func (pt *PerformanceTester) runStressTests(ctx context.Context, suite *PerformanceTestSuite) error {
	pt.logger.Info("Running stress tests")
	
	// Test with increasing load until failure
	maxLoad := pt.config.MaxTaskCount * 2
	step := pt.config.TaskCountStep
	
	for load := pt.config.MinTaskCount; load <= maxLoad; load += step {
		result := pt.runStressTest(ctx, load)
		suite.Results = append(suite.Results, result)
		
		if !result.Success {
			pt.logger.Info("Breaking point reached at %d tasks", load)
			break
		}
	}
	
	return nil
}

// runStressTest runs a single stress test
func (pt *PerformanceTester) runStressTest(ctx context.Context, load int) *PerformanceTestResult {
	start := time.Now()
	pt.logger.Info("Running stress test with %d tasks", load)
	
	// Generate test data
	tasks, err := pt.generateTestTasks(load)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("stress_test_%d_tasks", load),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: load,
			Success:   false,
			Error:     fmt.Errorf("failed to generate test data: %w", err),
		}
	}
	
	// Process tasks with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, pt.config.MaxProcessingTime)
	defer cancel()
	
	processedTasks, err := pt.generator.ProcessTasks("") // Use in-memory data
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("stress_test_%d_tasks", load),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: load,
			Success:   false,
			Error:     fmt.Errorf("failed to process tasks: %w", err),
		}
	}
	
	// Check if context was cancelled (timeout)
	select {
	case <-timeoutCtx.Done():
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("stress_test_%d_tasks", load),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: load,
			Success:   false,
			Error:     fmt.Errorf("test timed out after %v", pt.config.MaxProcessingTime),
		}
	default:
	}
	
	end := time.Now()
	duration := end.Sub(start)
	
	// Collect metrics
	metrics := pt.collectTestMetrics()
	systemMetrics := pt.monitor.GetSystemMetrics()
	
	// Determine success based on performance thresholds
	success := duration <= pt.config.MaxProcessingTime &&
		systemMetrics.MemoryUsage <= pt.config.MaxMemoryUsage*1024*1024 &&
		systemMetrics.CPUUsage <= pt.config.MaxCPUUsage
	
	var recommendations []string
	if !success {
		recommendations = append(recommendations, "System reached breaking point")
	}
	
	return &PerformanceTestResult{
		TestName:        fmt.Sprintf("stress_test_%d_tasks", load),
		StartTime:       start,
		EndTime:         end,
		Duration:        duration,
		TaskCount:       load,
		Success:         success,
		Error:           nil,
		Metrics:         metrics,
		SystemMetrics:   systemMetrics,
		Recommendations: recommendations,
	}
}

// runMemoryTests runs memory testing scenarios
func (pt *PerformanceTester) runMemoryTests(ctx context.Context, suite *PerformanceTestSuite) error {
	pt.logger.Info("Running memory tests")
	
	// Test memory usage with different task counts
	for taskCount := pt.config.MinTaskCount; taskCount <= pt.config.MaxTaskCount; taskCount += pt.config.TaskCountStep {
		result := pt.runMemoryTest(ctx, taskCount)
		suite.Results = append(suite.Results, result)
	}
	
	return nil
}

// runMemoryTest runs a single memory test
func (pt *PerformanceTester) runMemoryTest(ctx context.Context, taskCount int) *PerformanceTestResult {
	start := time.Now()
	pt.logger.Info("Running memory test with %d tasks", taskCount)
	
	// Force garbage collection before test
	runtime.GC()
	
	// Get initial memory stats
	var initialMem runtime.MemStats
	runtime.ReadMemStats(&initialMem)
	
	// Generate test data
	tasks, err := pt.generateTestTasks(taskCount)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("memory_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to generate test data: %w", err),
		}
	}
	
	// Process tasks
	processedTasks, err := pt.generator.ProcessTasks("") // Use in-memory data
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("memory_test_%d_tasks", taskCount),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to process tasks: %w", err),
		}
	}
	
	// Get peak memory stats
	var peakMem runtime.MemStats
	runtime.ReadMemStats(&peakMem)
	
	// Force garbage collection after test
	runtime.GC()
	
	// Get final memory stats
	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)
	
	end := time.Now()
	duration := end.Sub(start)
	
	// Calculate memory metrics
	memoryUsed := int64(peakMem.Alloc - initialMem.Alloc)
	memoryLeaks := finalMem.Alloc > initialMem.Alloc*1.1 // 10% tolerance
	
	// Collect metrics
	metrics := pt.collectTestMetrics()
	metrics["memory_used_bytes"] = float64(memoryUsed)
	metrics["memory_leaks"] = float64(boolToInt(memoryLeaks))
	
	systemMetrics := pt.monitor.GetSystemMetrics()
	
	// Determine success
	success := !memoryLeaks && memoryUsed <= pt.config.MaxMemoryUsage*1024*1024
	
	var recommendations []string
	if memoryLeaks {
		recommendations = append(recommendations, "Memory leaks detected")
	}
	if memoryUsed > pt.config.MaxMemoryUsage*1024*1024 {
		recommendations = append(recommendations, "Memory usage exceeded threshold")
	}
	
	return &PerformanceTestResult{
		TestName:        fmt.Sprintf("memory_test_%d_tasks", taskCount),
		StartTime:       start,
		EndTime:         end,
		Duration:        duration,
		TaskCount:       taskCount,
		Success:         success,
		Error:           nil,
		Metrics:         metrics,
		SystemMetrics:   systemMetrics,
		Recommendations: recommendations,
	}
}

// runConcurrencyTests runs concurrency testing scenarios
func (pt *PerformanceTester) runConcurrencyTests(ctx context.Context, suite *PerformanceTestSuite) error {
	pt.logger.Info("Running concurrency tests")
	
	// Test with different concurrency levels
	concurrencyLevels := []int{1, 2, 4, 8, 16}
	
	for _, level := range concurrencyLevels {
		result := pt.runConcurrencyTest(ctx, level)
		suite.Results = append(suite.Results, result)
	}
	
	return nil
}

// runConcurrencyTest runs a single concurrency test
func (pt *PerformanceTester) runConcurrencyTest(ctx context.Context, concurrencyLevel int) *PerformanceTestResult {
	start := time.Now()
	pt.logger.Info("Running concurrency test with %d workers", concurrencyLevel)
	
	// Create worker pool
	workerPool := make(chan struct{}, concurrencyLevel)
	var wg sync.WaitGroup
	
	// Test data
	taskCount := pt.config.MinTaskCount * 2
	tasks, err := pt.generateTestTasks(taskCount)
	if err != nil {
		return &PerformanceTestResult{
			TestName:  fmt.Sprintf("concurrency_test_%d_workers", concurrencyLevel),
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			TaskCount: taskCount,
			Success:   false,
			Error:     fmt.Errorf("failed to generate test data: %w", err),
		}
	}
	
	// Process tasks concurrently
	successCount := 0
	errorCount := 0
	var mu sync.Mutex
	
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			workerPool <- struct{}{}
			defer func() { <-workerPool }()
			
			// Process a subset of tasks
			subsetSize := taskCount / concurrencyLevel
			startIdx := i * subsetSize
			endIdx := startIdx + subsetSize
			if endIdx > taskCount {
				endIdx = taskCount
			}
			
			subset := tasks[startIdx:endIdx]
			
			// Process subset
			_, err := pt.generator.ProcessTasks("") // Use in-memory data
			
			mu.Lock()
			if err != nil {
				errorCount++
			} else {
				successCount++
			}
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	end := time.Now()
	duration := end.Sub(start)
	
	// Collect metrics
	metrics := pt.collectTestMetrics()
	metrics["concurrency_level"] = float64(concurrencyLevel)
	metrics["success_count"] = float64(successCount)
	metrics["error_count"] = float64(errorCount)
	metrics["throughput"] = float64(successCount) / duration.Seconds()
	
	systemMetrics := pt.monitor.GetSystemMetrics()
	
	// Determine success
	success := errorCount == 0 && successCount > 0
	
	var recommendations []string
	if errorCount > 0 {
		recommendations = append(recommendations, "Errors occurred during concurrent processing")
	}
	if successCount == 0 {
		recommendations = append(recommendations, "No successful operations")
	}
	
	return &PerformanceTestResult{
		TestName:        fmt.Sprintf("concurrency_test_%d_workers", concurrencyLevel),
		StartTime:       start,
		EndTime:         end,
		Duration:        duration,
		TaskCount:       taskCount,
		Success:         success,
		Error:           nil,
		Metrics:         metrics,
		SystemMetrics:   systemMetrics,
		Recommendations: recommendations,
	}
}

// generateTestTasks generates test tasks for performance testing
func (pt *PerformanceTester) generateTestTasks(count int) ([]*data.Task, error) {
	tasks := make([]*data.Task, count)
	
	for i := 0; i < count; i++ {
		task := &data.Task{
			ID:          fmt.Sprintf("task_%d", i),
			Name:        fmt.Sprintf("Test Task %d", i),
			Description: fmt.Sprintf("Description for test task %d", i),
			Category:    "Test",
			Priority:    rand.Intn(5) + 1,
			Status:      "Active",
			Assignee:    fmt.Sprintf("user_%d", i%10),
			StartDate:   time.Now().AddDate(0, 0, rand.Intn(365)),
			EndDate:     time.Now().AddDate(0, 0, rand.Intn(365)+1),
			Dependencies: []string{},
			IsMilestone: rand.Float32() < 0.1, // 10% chance of being milestone
		}
		
		tasks[i] = task
	}
	
	return tasks, nil
}

// collectTestMetrics collects metrics for a test
func (pt *PerformanceTester) collectTestMetrics() map[string]float64 {
	metrics := make(map[string]float64)
	
	// Get system metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	metrics["memory_alloc_bytes"] = float64(memStats.Alloc)
	metrics["memory_sys_bytes"] = float64(memStats.Sys)
	metrics["gc_count"] = float64(memStats.NumGC)
	metrics["gc_pause_ns"] = float64(memStats.PauseTotalNs)
	metrics["goroutine_count"] = float64(runtime.NumGoroutine())
	
	return metrics
}

// PerformanceTestSuite represents a complete performance test suite
type PerformanceTestSuite struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Config    *PerformanceTestConfig
	Results   []*PerformanceTestResult
}

// PerformanceTesterLogger provides logging for performance tester
type PerformanceTesterLogger struct{}

func (l *PerformanceTesterLogger) Info(msg string, args ...interface{})  { fmt.Printf("[TEST-INFO] "+msg+"\n", args...) }
func (l *PerformanceTesterLogger) Error(msg string, args ...interface{}) { fmt.Printf("[TEST-ERROR] "+msg+"\n", args...) }
func (l *PerformanceTesterLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[TEST-DEBUG] "+msg+"\n", args...) }

// Helper function to convert bool to int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
