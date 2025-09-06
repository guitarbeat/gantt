package generator

import (
	"context"
	"fmt"
	"time"

	"latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/data"
)

// PerformanceIntegration provides comprehensive performance optimization integration
type PerformanceIntegration struct {
	config           *PerformanceIntegrationConfig
	dataProcessor    *OptimizedDataProcessor
	layoutEngine     *OptimizedLayoutEngine
	pdfGenerator     *OptimizedPDFGenerator
	performanceMonitor *PerformanceMonitor
	performanceTester *PerformanceTester
	logger           PDFLogger
}

// PerformanceIntegrationConfig defines configuration for performance integration
type PerformanceIntegrationConfig struct {
	// Component settings
	EnableDataOptimization    bool `json:"enable_data_optimization"`
	EnableLayoutOptimization  bool `json:"enable_layout_optimization"`
	EnablePDFOptimization     bool `json:"enable_pdf_optimization"`
	EnablePerformanceMonitoring bool `json:"enable_performance_monitoring"`
	EnablePerformanceTesting  bool `json:"enable_performance_testing"`
	
	// Integration settings
	EnableCaching             bool `json:"enable_caching"`
	EnableParallelProcessing  bool `json:"enable_parallel_processing"`
	EnableMemoryOptimization  bool `json:"enable_memory_optimization"`
	EnableQualityOptimization bool `json:"enable_quality_optimization"`
	
	// Performance thresholds
	MaxProcessingTime         time.Duration `json:"max_processing_time"`
	MaxMemoryUsage            int64         `json:"max_memory_usage_mb"`
	MaxCPUUsage               float64       `json:"max_cpu_usage_percent"`
	
	// Monitoring settings
	EnableRealTimeMonitoring  bool          `json:"enable_realtime_monitoring"`
	MonitoringInterval        time.Duration `json:"monitoring_interval"`
	EnablePerformanceAlerts   bool          `json:"enable_performance_alerts"`
	
	// Testing settings
	EnableAutomatedTesting    bool          `json:"enable_automated_testing"`
	TestInterval              time.Duration `json:"test_interval"`
	EnableContinuousTesting   bool          `json:"enable_continuous_testing"`
}

// NewPerformanceIntegration creates a new performance integration system
func NewPerformanceIntegration() *PerformanceIntegration {
	config := GetDefaultPerformanceIntegrationConfig()
	
	return &PerformanceIntegration{
		config:            config,
		dataProcessor:     NewOptimizedDataProcessor(),
		layoutEngine:      NewOptimizedLayoutEngine(),
		pdfGenerator:      NewOptimizedPDFGenerator(),
		performanceMonitor: NewPerformanceMonitor(),
		performanceTester: NewPerformanceTester(),
		logger:            &PerformanceIntegrationLogger{},
	}
}

// GetDefaultPerformanceIntegrationConfig returns the default performance integration configuration
func GetDefaultPerformanceIntegrationConfig() *PerformanceIntegrationConfig {
	return &PerformanceIntegrationConfig{
		EnableDataOptimization:      true,
		EnableLayoutOptimization:    true,
		EnablePDFOptimization:       true,
		EnablePerformanceMonitoring: true,
		EnablePerformanceTesting:    true,
		EnableCaching:               true,
		EnableParallelProcessing:    true,
		EnableMemoryOptimization:    true,
		EnableQualityOptimization:   true,
		MaxProcessingTime:           time.Minute * 5,
		MaxMemoryUsage:              512, // 512MB
		MaxCPUUsage:                 80.0, // 80%
		EnableRealTimeMonitoring:    true,
		MonitoringInterval:          time.Second * 5,
		EnablePerformanceAlerts:     true,
		EnableAutomatedTesting:      true,
		TestInterval:                time.Minute * 30,
		EnableContinuousTesting:     false,
	}
}

// SetLogger sets the logger for the performance integration
func (pi *PerformanceIntegration) SetLogger(logger PDFLogger) {
	pi.logger = logger
	pi.dataProcessor.SetLogger(logger)
	pi.layoutEngine.SetLogger(logger)
	pi.pdfGenerator.SetLogger(logger)
	pi.performanceMonitor.SetLogger(logger)
	pi.performanceTester.SetLogger(logger)
}

// Initialize initializes the performance integration system
func (pi *PerformanceIntegration) Initialize(ctx context.Context) error {
	pi.logger.Info("Initializing performance integration system")
	
	// Initialize performance monitoring
	if pi.config.EnablePerformanceMonitoring {
		if err := pi.performanceMonitor.Start(ctx); err != nil {
			return fmt.Errorf("failed to start performance monitoring: %w", err)
		}
	}
	
	// Initialize performance testing
	if pi.config.EnablePerformanceTesting {
		// Run initial performance tests
		go pi.runInitialPerformanceTests(ctx)
	}
	
	pi.logger.Info("Performance integration system initialized successfully")
	return nil
}

// ProcessTasksWithOptimization processes tasks with full performance optimization
func (pi *PerformanceIntegration) ProcessTasksWithOptimization(ctx context.Context, filePath string) ([]*data.Task, error) {
	start := time.Now()
	pi.logger.Info("Starting optimized task processing")
	
	// Record start metrics
	pi.performanceMonitor.RecordTiming("task_processing_start", time.Since(start), map[string]string{"file": filePath})
	
	// Process tasks with optimization
	tasks, err := pi.dataProcessor.ProcessTasks(filePath)
	if err != nil {
		pi.performanceMonitor.RecordCounter("task_processing_errors", 1, map[string]string{"file": filePath})
		return nil, fmt.Errorf("task processing failed: %w", err)
	}
	
	// Record success metrics
	processingTime := time.Since(start)
	pi.performanceMonitor.RecordTiming("task_processing_time", processingTime, map[string]string{"file": filePath})
	pi.performanceMonitor.RecordCounter("task_processing_success", 1, map[string]string{"file": filePath})
	pi.performanceMonitor.RecordGauge("task_count", float64(len(tasks)), "count", map[string]string{"file": filePath})
	
	// Check performance thresholds
	if processingTime > pi.config.MaxProcessingTime {
		pi.logger.Warn("Task processing time exceeded threshold: %v > %v", processingTime, pi.config.MaxProcessingTime)
	}
	
	pi.logger.Info("Processed %d tasks in %v", len(tasks), processingTime)
	return tasks, nil
}

// GenerateLayoutWithOptimization generates layout with full performance optimization
func (pi *PerformanceIntegration) GenerateLayoutWithOptimization(tasks []*data.Task, config *calendar.CalendarConfig) (*calendar.CalendarLayout, error) {
	start := time.Now()
	pi.logger.Info("Starting optimized layout generation")
	
	// Record start metrics
	pi.performanceMonitor.RecordTiming("layout_generation_start", time.Since(start), map[string]string{"task_count": fmt.Sprintf("%d", len(tasks))})
	
	// Generate layout with optimization
	layout, err := pi.layoutEngine.GenerateLayout(tasks, config)
	if err != nil {
		pi.performanceMonitor.RecordCounter("layout_generation_errors", 1, map[string]string{"task_count": fmt.Sprintf("%d", len(tasks))})
		return nil, fmt.Errorf("layout generation failed: %w", err)
	}
	
	// Record success metrics
	generationTime := time.Since(start)
	pi.performanceMonitor.RecordTiming("layout_generation_time", generationTime, map[string]string{"task_count": fmt.Sprintf("%d", len(tasks))})
	pi.performanceMonitor.RecordCounter("layout_generation_success", 1, map[string]string{"task_count": fmt.Sprintf("%d", len(tasks))})
	
	// Check performance thresholds
	if generationTime > pi.config.MaxProcessingTime {
		pi.logger.Warn("Layout generation time exceeded threshold: %v > %v", generationTime, pi.config.MaxProcessingTime)
	}
	
	pi.logger.Info("Generated layout in %v", generationTime)
	return layout, nil
}

// GeneratePDFWithOptimization generates PDF with full performance optimization
func (pi *PerformanceIntegration) GeneratePDFWithOptimization(ctx context.Context, layout *calendar.CalendarLayout, outputPath string) error {
	start := time.Now()
	pi.logger.Info("Starting optimized PDF generation")
	
	// Record start metrics
	pi.performanceMonitor.RecordTiming("pdf_generation_start", time.Since(start), map[string]string{"output": outputPath})
	
	// Generate PDF with optimization
	err := pi.pdfGenerator.GeneratePDF(ctx, layout, outputPath)
	if err != nil {
		pi.performanceMonitor.RecordCounter("pdf_generation_errors", 1, map[string]string{"output": outputPath})
		return fmt.Errorf("PDF generation failed: %w", err)
	}
	
	// Record success metrics
	generationTime := time.Since(start)
	pi.performanceMonitor.RecordTiming("pdf_generation_time", generationTime, map[string]string{"output": outputPath})
	pi.performanceMonitor.RecordCounter("pdf_generation_success", 1, map[string]string{"output": outputPath})
	
	// Check performance thresholds
	if generationTime > pi.config.MaxProcessingTime {
		pi.logger.Warn("PDF generation time exceeded threshold: %v > %v", generationTime, pi.config.MaxProcessingTime)
	}
	
	pi.logger.Info("Generated PDF in %v", generationTime)
	return nil
}

// ProcessCompleteWorkflow processes the complete workflow with optimization
func (pi *PerformanceIntegration) ProcessCompleteWorkflow(ctx context.Context, filePath string, outputPath string, config *calendar.CalendarConfig) error {
	start := time.Now()
	pi.logger.Info("Starting complete optimized workflow")
	
	// Record start metrics
	pi.performanceMonitor.RecordTiming("complete_workflow_start", time.Since(start), map[string]string{"file": filePath, "output": outputPath})
	
	// Step 1: Process tasks
	tasks, err := pi.ProcessTasksWithOptimization(ctx, filePath)
	if err != nil {
		return fmt.Errorf("task processing failed: %w", err)
	}
	
	// Step 2: Generate layout
	layout, err := pi.GenerateLayoutWithOptimization(tasks, config)
	if err != nil {
		return fmt.Errorf("layout generation failed: %w", err)
	}
	
	// Step 3: Generate PDF
	err = pi.GeneratePDFWithOptimization(ctx, layout, outputPath)
	if err != nil {
		return fmt.Errorf("PDF generation failed: %w", err)
	}
	
	// Record success metrics
	totalTime := time.Since(start)
	pi.performanceMonitor.RecordTiming("complete_workflow_time", totalTime, map[string]string{"file": filePath, "output": outputPath})
	pi.performanceMonitor.RecordCounter("complete_workflow_success", 1, map[string]string{"file": filePath, "output": outputPath})
	
	// Check performance thresholds
	if totalTime > pi.config.MaxProcessingTime {
		pi.logger.Warn("Complete workflow time exceeded threshold: %v > %v", totalTime, pi.config.MaxProcessingTime)
	}
	
	pi.logger.Info("Complete workflow finished in %v", totalTime)
	return nil
}

// GetPerformanceMetrics returns current performance metrics
func (pi *PerformanceIntegration) GetPerformanceMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})
	
	// Get system metrics
	systemMetrics := pi.performanceMonitor.GetSystemMetrics()
	metrics["system"] = systemMetrics
	
	// Get current metrics
	currentMetrics := pi.performanceMonitor.GetMetrics()
	metrics["current"] = currentMetrics
	
	// Get performance report
	report := pi.performanceMonitor.GetPerformanceReport()
	metrics["report"] = report
	
	return metrics
}

// GetPerformanceReport returns a comprehensive performance report
func (pi *PerformanceIntegration) GetPerformanceReport() *PerformanceReport {
	return pi.performanceMonitor.GetPerformanceReport()
}

// RunPerformanceTests runs performance tests
func (pi *PerformanceIntegration) RunPerformanceTests(ctx context.Context) (*PerformanceTestSuite, error) {
	pi.logger.Info("Running performance tests")
	
	suite, err := pi.performanceTester.RunPerformanceTests(ctx)
	if err != nil {
		return nil, fmt.Errorf("performance testing failed: %w", err)
	}
	
	pi.logger.Info("Performance tests completed: %d tests, %d successful", len(suite.Results), pi.countSuccessfulTests(suite.Results))
	return suite, nil
}

// countSuccessfulTests counts successful tests
func (pi *PerformanceIntegration) countSuccessfulTests(results []*PerformanceTestResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}

// runInitialPerformanceTests runs initial performance tests
func (pi *PerformanceIntegration) runInitialPerformanceTests(ctx context.Context) {
	pi.logger.Info("Running initial performance tests")
	
	suite, err := pi.performanceTester.RunPerformanceTests(ctx)
	if err != nil {
		pi.logger.Error("Initial performance tests failed: %v", err)
		return
	}
	
	pi.logger.Info("Initial performance tests completed: %d tests, %d successful", 
		len(suite.Results), pi.countSuccessfulTests(suite.Results))
}

// StartContinuousMonitoring starts continuous performance monitoring
func (pi *PerformanceIntegration) StartContinuousMonitoring(ctx context.Context) {
	if !pi.config.EnableRealTimeMonitoring {
		return
	}
	
	pi.logger.Info("Starting continuous performance monitoring")
	
	ticker := time.NewTicker(pi.config.MonitoringInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pi.performContinuousMonitoring()
		}
	}
}

// performContinuousMonitoring performs continuous monitoring
func (pi *PerformanceIntegration) performContinuousMonitoring() {
	// Get current metrics
	metrics := pi.performanceMonitor.GetMetrics()
	systemMetrics := pi.performanceMonitor.GetSystemMetrics()
	
	// Check for performance issues
	pi.checkPerformanceThresholds(metrics, systemMetrics)
	
	// Log performance status
	pi.logger.Debug("Performance monitoring: Memory: %.2f MB, CPU: %.2f%%, Goroutines: %d", 
		float64(systemMetrics.MemoryUsage)/(1024*1024), 
		systemMetrics.CPUUsage, 
		systemMetrics.GoroutineCount)
}

// checkPerformanceThresholds checks performance thresholds
func (pi *PerformanceIntegration) checkPerformanceThresholds(metrics map[string]*MetricData, systemMetrics *SystemMetrics) {
	// Check memory usage
	memoryMB := float64(systemMetrics.MemoryUsage) / (1024 * 1024)
	if memoryMB > float64(pi.config.MaxMemoryUsage) {
		pi.logger.Warn("Memory usage exceeded threshold: %.2f MB > %d MB", memoryMB, pi.config.MaxMemoryUsage)
	}
	
	// Check CPU usage
	if systemMetrics.CPUUsage > pi.config.MaxCPUUsage {
		pi.logger.Warn("CPU usage exceeded threshold: %.2f%% > %.2f%%", systemMetrics.CPUUsage, pi.config.MaxCPUUsage)
	}
	
	// Check goroutine count
	if systemMetrics.GoroutineCount > 1000 {
		pi.logger.Warn("High goroutine count: %d", systemMetrics.GoroutineCount)
	}
}

// StartContinuousTesting starts continuous performance testing
func (pi *PerformanceIntegration) StartContinuousTesting(ctx context.Context) {
	if !pi.config.EnableContinuousTesting {
		return
	}
	
	pi.logger.Info("Starting continuous performance testing")
	
	ticker := time.NewTicker(pi.config.TestInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pi.runContinuousTests(ctx)
		}
	}
}

// runContinuousTests runs continuous tests
func (pi *PerformanceIntegration) runContinuousTests(ctx context.Context) {
	pi.logger.Info("Running continuous performance tests")
	
	suite, err := pi.performanceTester.RunPerformanceTests(ctx)
	if err != nil {
		pi.logger.Error("Continuous performance tests failed: %v", err)
		return
	}
	
	successCount := pi.countSuccessfulTests(suite.Results)
	pi.logger.Info("Continuous performance tests completed: %d tests, %d successful", 
		len(suite.Results), successCount)
	
	// Alert if performance degraded
	if successCount < len(suite.Results)*0.8 { // Less than 80% success rate
		pi.logger.Warn("Performance degradation detected: %d/%d tests successful", successCount, len(suite.Results))
	}
}

// Stop stops the performance integration system
func (pi *PerformanceIntegration) Stop() {
	pi.logger.Info("Stopping performance integration system")
	
	// Stop performance monitoring
	pi.performanceMonitor.Stop()
	
	pi.logger.Info("Performance integration system stopped")
}

// PerformanceIntegrationLogger provides logging for performance integration
type PerformanceIntegrationLogger struct{}

func (l *PerformanceIntegrationLogger) Info(msg string, args ...interface{})  { fmt.Printf("[PERF-INTEGRATION-INFO] "+msg+"\n", args...) }
func (l *PerformanceIntegrationLogger) Error(msg string, args ...interface{}) { fmt.Printf("[PERF-INTEGRATION-ERROR] "+msg+"\n", args...) }
func (l *PerformanceIntegrationLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[PERF-INTEGRATION-DEBUG] "+msg+"\n", args...) }
func (l *PerformanceIntegrationLogger) Warn(msg string, args ...interface{})  { fmt.Printf("[PERF-INTEGRATION-WARN] "+msg+"\n", args...) }
