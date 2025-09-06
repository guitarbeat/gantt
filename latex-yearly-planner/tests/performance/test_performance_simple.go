package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"latex-yearly-planner/internal/generator"
)

func main() {
	fmt.Println("=== Simple Performance Optimization Test ===")
	
	// Create performance optimizer
	optimizer := generator.NewSimplePerformanceOptimizer()
	
	// Set up logging
	optimizer.SetLogger(&TestLogger{})
	
	// Test 1: Basic Performance Test
	fmt.Println("\n1. Running Basic Performance Test...")
	if err := testBasicPerformance(optimizer); err != nil {
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
	
	// Generate final performance report
	fmt.Println("\n5. Generating Final Performance Report...")
	metrics := optimizer.GetPerformanceMetrics()
	systemMetrics := optimizer.GetSystemMetrics()
	
	fmt.Printf("Performance Metrics:\n")
	for name, value := range metrics {
		fmt.Printf("  %s: %.2f\n", name, value)
	}
	
	fmt.Printf("\nSystem Metrics:\n")
	for name, value := range systemMetrics {
		fmt.Printf("  %s: %.2f\n", name, value)
	}
	
	fmt.Println("\n=== Simple Performance Optimization Test Completed ===")
}

// testBasicPerformance tests basic performance functionality
func testBasicPerformance(optimizer *generator.SimplePerformanceOptimizer) error {
	ctx := context.Background()
	
	// Test data processing
	fmt.Println("  - Testing data processing...")
	start := time.Now()
	_, err := optimizer.OptimizeDataProcessing(ctx, "test data")
	if err != nil {
		return fmt.Errorf("data processing failed: %w", err)
	}
	processingTime := time.Since(start)
	fmt.Printf("    - Data processing time: %v\n", processingTime)
	
	// Test layout generation
	fmt.Println("  - Testing layout generation...")
	start = time.Now()
	_, err = optimizer.OptimizeLayoutGeneration(ctx, "test layout")
	if err != nil {
		return fmt.Errorf("layout generation failed: %w", err)
	}
	layoutTime := time.Since(start)
	fmt.Printf("    - Layout generation time: %v\n", layoutTime)
	
	// Test PDF generation
	fmt.Println("  - Testing PDF generation...")
	start = time.Now()
	_, err = optimizer.OptimizePDFGeneration(ctx, "test pdf")
	if err != nil {
		return fmt.Errorf("PDF generation failed: %w", err)
	}
	pdfTime := time.Since(start)
	fmt.Printf("    - PDF generation time: %v\n", pdfTime)
	
	// Check performance thresholds
	totalTime := processingTime + layoutTime + pdfTime
	if totalTime > time.Second*2 {
		return fmt.Errorf("total processing time exceeded threshold: %v", totalTime)
	}
	
	return nil
}

// testCaching tests caching functionality
func testCaching(optimizer *generator.SimplePerformanceOptimizer) error {
	ctx := context.Background()
	
	// First call (should be slow)
	fmt.Println("  - Testing first call (cache miss)...")
	start := time.Now()
	_, err := optimizer.OptimizeDataProcessing(ctx, "test data")
	if err != nil {
		return fmt.Errorf("first data processing failed: %w", err)
	}
	firstTime := time.Since(start)
	fmt.Printf("    - First call time: %v\n", firstTime)
	
	// Second call (should be fast due to caching)
	fmt.Println("  - Testing second call (cache hit)...")
	start = time.Now()
	_, err = optimizer.OptimizeDataProcessing(ctx, "test data")
	if err != nil {
		return fmt.Errorf("second data processing failed: %w", err)
	}
	secondTime := time.Since(start)
	fmt.Printf("    - Second call time: %v\n", secondTime)
	
	// Check if caching improved performance
	if secondTime >= firstTime {
		return fmt.Errorf("caching did not improve performance: %v >= %v", secondTime, firstTime)
	}
	
	return nil
}

// testMemory tests memory usage
func testMemory(optimizer *generator.SimplePerformanceOptimizer) error {
	ctx := context.Background()
	
	// Get initial memory
	initialMetrics := optimizer.GetSystemMetrics()
	initialMemory := initialMetrics["memory_alloc_bytes"]
	fmt.Printf("  - Initial memory: %.2f MB\n", initialMemory/(1024*1024))
	
	// Process multiple items to test memory usage
	fmt.Println("  - Processing multiple items...")
	for i := 0; i < 10; i++ {
		_, err := optimizer.OptimizeDataProcessing(ctx, fmt.Sprintf("test data %d", i))
		if err != nil {
			return fmt.Errorf("data processing %d failed: %w", i, err)
		}
	}
	
	// Get final memory
	finalMetrics := optimizer.GetSystemMetrics()
	finalMemory := finalMetrics["memory_alloc_bytes"]
	fmt.Printf("  - Final memory: %.2f MB\n", finalMemory/(1024*1024))
	
	// Check memory usage
	memoryIncrease := finalMemory - initialMemory
	if memoryIncrease > 100*1024*1024 { // 100MB
		return fmt.Errorf("memory usage increased too much: %.2f MB", memoryIncrease/(1024*1024))
	}
	
	return nil
}

// testPerformanceMetrics tests performance metrics collection
func testPerformanceMetrics(optimizer *generator.SimplePerformanceOptimizer) error {
	ctx := context.Background()
	
	// Process some data to generate metrics
	_, err := optimizer.OptimizeDataProcessing(ctx, "test data")
	if err != nil {
		return fmt.Errorf("data processing failed: %w", err)
	}
	
	// Get metrics
	metrics := optimizer.GetPerformanceMetrics()
	if len(metrics) == 0 {
		return fmt.Errorf("no metrics collected")
	}
	
	// Check for expected metrics
	expectedMetrics := []string{"data_processing_time", "data_processing_success"}
	for _, expected := range expectedMetrics {
		if _, exists := metrics[expected]; !exists {
			return fmt.Errorf("expected metric %s not found", expected)
		}
	}
	
	fmt.Printf("  - Collected %d metrics\n", len(metrics))
	return nil
}

// TestLogger provides logging for tests
type TestLogger struct{}

func (l *TestLogger) Info(msg string, args ...interface{})  { fmt.Printf("[TEST-INFO] "+msg+"\n", args...) }
func (l *TestLogger) Error(msg string, args ...interface{}) { fmt.Printf("[TEST-ERROR] "+msg+"\n", args...) }
func (l *TestLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[TEST-DEBUG] "+msg+"\n", args...) }
