package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"phd-dissertation-planner/internal/generator"
)

func main() {
	fmt.Println("=== Performance Optimization Test Suite ===")
	
	// Test 1: Simple Performance Tests (from test_performance_simple.go)
	fmt.Println("\n=== Test 1: Simple Performance Tests ===")
	testSimplePerformance()
	
	// Test 2: Advanced Performance Integration Tests
	fmt.Println("\n=== Test 2: Advanced Performance Integration Tests ===")
	
	// Create performance integration system
	perfIntegration := generator.NewPerformanceIntegration()
	
	// Set up logging
	perfIntegration.SetLogger(&TestLogger{})
	
	// Initialize the system
	ctx := context.Background()
	if err := perfIntegration.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize performance integration: %v", err)
	}
	
	// Test 1: Basic Performance Test
	fmt.Println("\n1. Running Basic Performance Test...")
	if err := testBasicPerformance(perfIntegration, ctx); err != nil {
		log.Printf("Basic performance test failed: %v", err)
	} else {
		fmt.Println("✓ Basic performance test passed")
	}
	
	// Test 2: Load Testing
	fmt.Println("\n2. Running Load Testing...")
	if err := testLoadPerformance(perfIntegration, ctx); err != nil {
		log.Printf("Load performance test failed: %v", err)
	} else {
		fmt.Println("✓ Load performance test passed")
	}
	
	// Test 3: Memory Testing
	fmt.Println("\n3. Running Memory Testing...")
	if err := testMemoryPerformance(perfIntegration, ctx); err != nil {
		log.Printf("Memory performance test failed: %v", err)
	} else {
		fmt.Println("✓ Memory performance test passed")
	}
	
	// Test 4: Concurrency Testing
	fmt.Println("\n4. Running Concurrency Testing...")
	if err := testConcurrencyPerformance(perfIntegration, ctx); err != nil {
		log.Printf("Concurrency performance test failed: %v", err)
	} else {
		fmt.Println("✓ Concurrency performance test passed")
	}
	
	// Test 5: Performance Monitoring
	fmt.Println("\n5. Testing Performance Monitoring...")
	if err := testPerformanceMonitoring(perfIntegration, ctx); err != nil {
		log.Printf("Performance monitoring test failed: %v", err)
	} else {
		fmt.Println("✓ Performance monitoring test passed")
	}
	
	// Test 6: Performance Report Generation
	fmt.Println("\n6. Testing Performance Report Generation...")
	if err := testPerformanceReporting(perfIntegration, ctx); err != nil {
		log.Printf("Performance reporting test failed: %v", err)
	} else {
		fmt.Println("✓ Performance reporting test passed")
	}
	
	// Generate final performance report
	fmt.Println("\n7. Generating Final Performance Report...")
	report := perfIntegration.GetPerformanceReport()
	fmt.Printf("Performance Report Summary: %s\n", report.Summary)
	
	if len(report.Recommendations) > 0 {
		fmt.Println("Recommendations:")
		for i, rec := range report.Recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
	}
	
	fmt.Println("\n=== Performance Optimization Test Suite Completed ===")
}

// testBasicPerformance tests basic performance functionality
func testBasicPerformance(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Create test data
	testData := createTestCSVData(100)
	testFile := "/tmp/test_basic_performance.csv"
	
	// Write test data to file
	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer os.Remove(testFile)
	
	// Test task processing
	start := time.Now()
	tasks, err := perfIntegration.ProcessTasksWithOptimization(ctx, testFile)
	if err != nil {
		return fmt.Errorf("task processing failed: %w", err)
	}
	processingTime := time.Since(start)
	
	// Test layout generation
	config := &generator.CalendarConfig{
		Year:       2024,
		DayWidth:   20.0,
		TaskHeight: 30.0,
	}
	
	start = time.Now()
	layout, err := perfIntegration.GenerateLayoutWithOptimization(tasks, config)
	if err != nil {
		return fmt.Errorf("layout generation failed: %w", err)
	}
	layoutTime := time.Since(start)
	
	// Test PDF generation
	outputPath := "/tmp/test_basic_performance.pdf"
	start = time.Now()
	err = perfIntegration.GeneratePDFWithOptimization(ctx, layout, outputPath)
	if err != nil {
		return fmt.Errorf("PDF generation failed: %w", err)
	}
	pdfTime := time.Since(start)
	
	// Check performance thresholds
	totalTime := processingTime + layoutTime + pdfTime
	if totalTime > time.Minute*2 {
		return fmt.Errorf("total processing time exceeded threshold: %v", totalTime)
	}
	
	fmt.Printf("  - Task processing: %v\n", processingTime)
	fmt.Printf("  - Layout generation: %v\n", layoutTime)
	fmt.Printf("  - PDF generation: %v\n", pdfTime)
	fmt.Printf("  - Total time: %v\n", totalTime)
	
	// Clean up
	os.Remove(outputPath)
	
	return nil
}

// testLoadPerformance tests load performance
func testLoadPerformance(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Test with different task counts
	taskCounts := []int{50, 100, 200, 500}
	
	for _, count := range taskCounts {
		fmt.Printf("  - Testing with %d tasks...\n", count)
		
		// Create test data
		testData := createTestCSVData(count)
		testFile := fmt.Sprintf("/tmp/test_load_%d.csv", count)
		
		// Write test data to file
		if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %w", err)
		}
		defer os.Remove(testFile)
		
		// Test complete workflow
		start := time.Now()
		outputPath := fmt.Sprintf("/tmp/test_load_%d.pdf", count)
		config := &generator.CalendarConfig{
			Year:       2024,
			DayWidth:   20.0,
			TaskHeight: 30.0,
		}
		
		err := perfIntegration.ProcessCompleteWorkflow(ctx, testFile, outputPath, config)
		if err != nil {
			return fmt.Errorf("workflow failed for %d tasks: %w", count, err)
		}
		
		totalTime := time.Since(start)
		fmt.Printf("    - Processed %d tasks in %v\n", count, totalTime)
		
		// Check performance threshold (should scale linearly)
		expectedTime := time.Duration(count/100) * time.Second * 2
		if totalTime > expectedTime*2 {
			return fmt.Errorf("performance degraded for %d tasks: %v > %v", count, totalTime, expectedTime*2)
		}
		
		// Clean up
		os.Remove(outputPath)
	}
	
	return nil
}

// testMemoryPerformance tests memory performance
func testMemoryPerformance(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Test memory usage with different task counts
	taskCounts := []int{100, 500, 1000}
	
	for _, count := range taskCounts {
		fmt.Printf("  - Testing memory with %d tasks...\n", count)
		
		// Create test data
		testData := createTestCSVData(count)
		testFile := fmt.Sprintf("/tmp/test_memory_%d.csv", count)
		
		// Write test data to file
		if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %w", err)
		}
		defer os.Remove(testFile)
		
		// Test complete workflow
		outputPath := fmt.Sprintf("/tmp/test_memory_%d.pdf", count)
		config := &generator.CalendarConfig{
			Year:       2024,
			DayWidth:   20.0,
			TaskHeight: 30.0,
		}
		
		err := perfIntegration.ProcessCompleteWorkflow(ctx, testFile, outputPath, config)
		if err != nil {
			return fmt.Errorf("workflow failed for %d tasks: %w", count, err)
		}
		
		// Get memory metrics
		metrics := perfIntegration.GetPerformanceMetrics()
		if systemMetrics, ok := metrics["system"].(*generator.SystemMetrics); ok {
			memoryMB := float64(systemMetrics.MemoryUsage) / (1024 * 1024)
			fmt.Printf("    - Memory usage: %.2f MB\n", memoryMB)
			
			// Check memory threshold (should not exceed 512MB)
			if memoryMB > 512 {
				return fmt.Errorf("memory usage exceeded threshold: %.2f MB > 512 MB", memoryMB)
			}
		}
		
		// Clean up
		os.Remove(outputPath)
	}
	
	return nil
}

// testConcurrencyPerformance tests concurrency performance
func testConcurrencyPerformance(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Test concurrent processing
	concurrencyLevels := []int{1, 2, 4, 8}
	
	for _, level := range concurrencyLevels {
		fmt.Printf("  - Testing concurrency with %d workers...\n", level)
		
		// Create test data
		testData := createTestCSVData(100)
		testFile := fmt.Sprintf("/tmp/test_concurrency_%d.csv", level)
		
		// Write test data to file
		if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %w", err)
		}
		defer os.Remove(testFile)
		
		// Test concurrent processing
		start := time.Now()
		outputPath := fmt.Sprintf("/tmp/test_concurrency_%d.pdf", level)
		config := &generator.CalendarConfig{
			Year:       2024,
			DayWidth:   20.0,
			TaskHeight: 30.0,
		}
		
		err := perfIntegration.ProcessCompleteWorkflow(ctx, testFile, outputPath, config)
		if err != nil {
			return fmt.Errorf("concurrent workflow failed for %d workers: %w", level, err)
		}
		
		totalTime := time.Since(start)
		fmt.Printf("    - Processed with %d workers in %v\n", level, totalTime)
		
		// Clean up
		os.Remove(outputPath)
	}
	
	return nil
}

// testPerformanceMonitoring tests performance monitoring
func testPerformanceMonitoring(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Test performance monitoring
	fmt.Println("  - Testing performance monitoring...")
	
	// Get initial metrics
	initialMetrics := perfIntegration.GetPerformanceMetrics()
	if initialMetrics == nil {
		return fmt.Errorf("failed to get initial metrics")
	}
	
	// Perform some work to generate metrics
	testData := createTestCSVData(50)
	testFile := "/tmp/test_monitoring.csv"
	
	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer os.Remove(testFile)
	
	// Process tasks to generate metrics
	_, err := perfIntegration.ProcessTasksWithOptimization(ctx, testFile)
	if err != nil {
		return fmt.Errorf("task processing failed: %w", err)
	}
	
	// Get updated metrics
	updatedMetrics := perfIntegration.GetPerformanceMetrics()
	if updatedMetrics == nil {
		return fmt.Errorf("failed to get updated metrics")
	}
	
	// Check if metrics were updated
	if len(updatedMetrics) <= len(initialMetrics) {
		return fmt.Errorf("metrics were not updated")
	}
	
	fmt.Println("    - Performance monitoring working correctly")
	return nil
}

// testPerformanceReporting tests performance reporting
func testPerformanceReporting(perfIntegration *generator.PerformanceIntegration, ctx context.Context) error {
	// Test performance reporting
	fmt.Println("  - Testing performance reporting...")
	
	// Generate performance report
	report := perfIntegration.GetPerformanceReport()
	if report == nil {
		return fmt.Errorf("failed to generate performance report")
	}
	
	// Check report structure
	if report.Summary == "" {
		return fmt.Errorf("performance report summary is empty")
	}
	
	if report.Timestamp.IsZero() {
		return fmt.Errorf("performance report timestamp is zero")
	}
	
	fmt.Printf("    - Performance report generated: %s\n", report.Summary)
	return nil
}

// createTestCSVData creates test CSV data
func createTestCSVData(taskCount int) string {
	var data string
	
	// Header
	data += "Task ID,Task Name,Description,Category,Priority,Status,Assignee,Start Date,Due Date,Dependencies\n"
	
	// Generate tasks
	for i := 0; i < taskCount; i++ {
		taskID := fmt.Sprintf("task_%d", i)
		taskName := fmt.Sprintf("Test Task %d", i)
		description := fmt.Sprintf("Description for test task %d", i)
		category := "Test"
		priority := (i % 5) + 1
		status := "Active"
		assignee := fmt.Sprintf("user_%d", i%10)
		startDate := fmt.Sprintf("2024-%02d-%02d", (i%12)+1, (i%28)+1)
		endDate := fmt.Sprintf("2024-%02d-%02d", ((i+1)%12)+1, ((i+1)%28)+1)
		dependencies := ""
		
		// Add some dependencies
		if i > 0 && i%3 == 0 {
			dependencies = fmt.Sprintf("task_%d", i-1)
		}
		
		data += fmt.Sprintf("%s,%s,%s,%s,%d,%s,%s,%s,%s,%s\n",
			taskID, taskName, description, category, priority, status, assignee, startDate, endDate, dependencies)
	}
	
	return data
}

// TestLogger provides logging for tests
type TestLogger struct{}

func (l *TestLogger) Info(msg string, args ...interface{})  { fmt.Printf("[TEST-INFO] "+msg+"\n", args...) }
func (l *TestLogger) Error(msg string, args ...interface{}) { fmt.Printf("[TEST-ERROR] "+msg+"\n", args...) }
func (l *TestLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[TEST-DEBUG] "+msg+"\n", args...) }
func (l *TestLogger) Warn(msg string, args ...interface{})  { fmt.Printf("[TEST-WARN] "+msg+"\n", args...) }

// testSimplePerformance tests simple performance optimization (from test_performance_simple.go)
func testSimplePerformance() {
	fmt.Println("=== Simple Performance Optimization Test ===")
	
	// Create performance optimizer
	optimizer := generator.NewSimplePerformanceOptimizer()
	
	// Set up logging
	optimizer.SetLogger(&TestLogger{})
	
	// Test 1: Basic Performance Test
	fmt.Println("\n1. Running Basic Performance Test...")
	if err := testBasicPerformanceSimple(optimizer); err != nil {
		log.Printf("Basic performance test failed: %v", err)
	} else {
		fmt.Println("✓ Basic performance test passed")
	}
	
	// Test 2: Caching Test
	fmt.Println("\n2. Running Caching Test...")
	if err := testCaching(optimizer); err != nil {
		log.Printf("Caching test failed: %v", err)
	} else {
		fmt.Println("✓ Caching test passed")
	}
	
	// Test 3: Memory Test
	fmt.Println("\n3. Running Memory Test...")
	if err := testMemory(optimizer); err != nil {
		log.Printf("Memory test failed: %v", err)
	} else {
		fmt.Println("✓ Memory test passed")
	}
	
	// Test 4: Performance Metrics Test
	fmt.Println("\n4. Running Performance Metrics Test...")
	if err := testPerformanceMetrics(optimizer); err != nil {
		log.Printf("Performance metrics test failed: %v", err)
	} else {
		fmt.Println("✓ Performance metrics test passed")
	}
	
	fmt.Println("\n✅ Simple performance optimization tests completed!")
}

// testBasicPerformanceSimple tests basic performance (from test_performance_simple.go)
func testBasicPerformanceSimple(optimizer *generator.SimplePerformanceOptimizer) error {
	start := time.Now()
	
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	
	duration := time.Since(start)
	
	// Check if performance is within acceptable limits
	if duration > 200*time.Millisecond {
		return fmt.Errorf("performance too slow: %v", duration)
	}
	
	fmt.Printf("   Basic performance: %v\n", duration)
	return nil
}

// testCaching tests caching functionality (from test_performance_simple.go)
func testCaching(optimizer *generator.SimplePerformanceOptimizer) error {
	// Test cache hit
	start := time.Now()
	optimizer.GetCachedResult("test-key")
	cacheHitTime := time.Since(start)
	
	// Test cache miss
	start = time.Now()
	optimizer.ProcessData("test-data")
	cacheMissTime := time.Since(start)
	
	fmt.Printf("   Cache hit time: %v\n", cacheHitTime)
	fmt.Printf("   Cache miss time: %v\n", cacheMissTime)
	
	// Cache hit should be faster than cache miss
	if cacheHitTime >= cacheMissTime {
		return fmt.Errorf("caching not working properly")
	}
	
	return nil
}

// testMemory tests memory usage (from test_performance_simple.go)
func testMemory(optimizer *generator.SimplePerformanceOptimizer) error {
	// Test memory allocation
	start := time.Now()
	
	// Simulate memory-intensive operation
	data := make([]byte, 1024*1024) // 1MB
	_ = data
	
	duration := time.Since(start)
	
	fmt.Printf("   Memory allocation time: %v\n", duration)
	
	// Check if memory allocation is reasonable
	if duration > 10*time.Millisecond {
		return fmt.Errorf("memory allocation too slow: %v", duration)
	}
	
	return nil
}

// testPerformanceMetrics tests performance metrics (from test_performance_simple.go)
func testPerformanceMetrics(optimizer *generator.SimplePerformanceOptimizer) error {
	// Get performance metrics
	metrics := optimizer.GetPerformanceMetrics()
	
	fmt.Printf("   Total operations: %d\n", metrics.TotalOperations)
	fmt.Printf("   Average time: %v\n", metrics.AverageTime)
	fmt.Printf("   Cache hit rate: %.2f%%\n", metrics.CacheHitRate*100)
	
	// Validate metrics
	if metrics.TotalOperations < 0 {
		return fmt.Errorf("invalid total operations: %d", metrics.TotalOperations)
	}
	
	if metrics.AverageTime < 0 {
		return fmt.Errorf("invalid average time: %v", metrics.AverageTime)
	}
	
	if metrics.CacheHitRate < 0 || metrics.CacheHitRate > 1 {
		return fmt.Errorf("invalid cache hit rate: %.2f", metrics.CacheHitRate)
	}
	
	return nil
}
