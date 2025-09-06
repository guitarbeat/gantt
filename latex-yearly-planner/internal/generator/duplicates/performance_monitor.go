package generator

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// PerformanceMonitor provides comprehensive performance monitoring and analytics
type PerformanceMonitor struct {
	config        *PerformanceMonitorConfig
	metrics       *PerformanceMetrics
	collector     *MetricsCollector
	analyzer      *PerformanceAnalyzer
	reporter      *PerformanceReporter
	logger        PDFLogger
}

// PerformanceMonitorConfig defines configuration for performance monitoring
type PerformanceMonitorConfig struct {
	// Monitoring settings
	EnableMonitoring      bool          `json:"enable_monitoring"`
	EnableRealTimeMetrics bool          `json:"enable_realtime_metrics"`
	MetricsInterval       time.Duration `json:"metrics_interval"`
	MaxMetricsHistory     int           `json:"max_metrics_history"`
	
	// Performance thresholds
	MaxProcessingTime     time.Duration `json:"max_processing_time"`
	MaxMemoryUsage        int64         `json:"max_memory_usage_mb"`
	MaxCPUUsage           float64       `json:"max_cpu_usage_percent"`
	
	// Alerting
	EnableAlerts          bool          `json:"enable_alerts"`
	AlertThresholds       *AlertThresholds `json:"alert_thresholds"`
	
	// Reporting
	EnableReporting       bool          `json:"enable_reporting"`
	ReportInterval        time.Duration `json:"report_interval"`
	EnableDetailedReports bool          `json:"enable_detailed_reports"`
}

// AlertThresholds defines thresholds for performance alerts
type AlertThresholds struct {
	ProcessingTime time.Duration `json:"processing_time"`
	MemoryUsage    int64         `json:"memory_usage_mb"`
	CPUUsage       float64       `json:"cpu_usage_percent"`
	ErrorRate      float64       `json:"error_rate_percent"`
}

// PerformanceMetrics stores performance metrics
type PerformanceMetrics struct {
	config        *PerformanceMonitorConfig
	metrics       map[string]*MetricData
	history       []*MetricsSnapshot
	mu            sync.RWMutex
}

// MetricData represents a single performance metric
type MetricData struct {
	Name        string
	Value       float64
	Unit        string
	Timestamp   time.Time
	Tags        map[string]string
	Description string
}

// MetricsSnapshot represents a snapshot of all metrics at a point in time
type MetricsSnapshot struct {
	Timestamp time.Time
	Metrics   map[string]*MetricData
	System    *SystemMetrics
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	MemoryUsage    int64   `json:"memory_usage_bytes"`
	MemoryPercent  float64 `json:"memory_usage_percent"`
	CPUUsage       float64 `json:"cpu_usage_percent"`
	GoroutineCount int     `json:"goroutine_count"`
	GCCount        int64   `json:"gc_count"`
	GCPause        time.Duration `json:"gc_pause_ns"`
}

// MetricsCollector collects performance metrics
type MetricsCollector struct {
	config        *PerformanceMonitorConfig
	monitor       *PerformanceMonitor
	logger        PDFLogger
}

// PerformanceAnalyzer analyzes performance metrics
type PerformanceAnalyzer struct {
	config        *PerformanceMonitorConfig
	monitor       *PerformanceMonitor
	logger        PDFLogger
}

// PerformanceReporter generates performance reports
type PerformanceReporter struct {
	config        *PerformanceMonitorConfig
	monitor       *PerformanceMonitor
	logger        PDFLogger
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	config := GetDefaultPerformanceMonitorConfig()
	
	return &PerformanceMonitor{
		config:    config,
		metrics:   NewPerformanceMetrics(config),
		collector: NewMetricsCollector(config),
		analyzer:  NewPerformanceAnalyzer(config),
		reporter:  NewPerformanceReporter(config),
		logger:    &PerformanceMonitorLogger{},
	}
}

// GetDefaultPerformanceMonitorConfig returns the default performance monitor configuration
func GetDefaultPerformanceMonitorConfig() *PerformanceMonitorConfig {
	return &PerformanceMonitorConfig{
		EnableMonitoring:       true,
		EnableRealTimeMetrics:  true,
		MetricsInterval:        time.Second * 5,
		MaxMetricsHistory:      1000,
		MaxProcessingTime:      time.Minute * 5,
		MaxMemoryUsage:         512, // 512MB
		MaxCPUUsage:            80.0, // 80%
		EnableAlerts:           true,
		AlertThresholds: &AlertThresholds{
			ProcessingTime: time.Minute * 2,
			MemoryUsage:    256, // 256MB
			CPUUsage:       70.0, // 70%
			ErrorRate:      5.0,  // 5%
		},
		EnableReporting:        true,
		ReportInterval:         time.Minute * 10,
		EnableDetailedReports:  true,
	}
}

// SetLogger sets the logger for the performance monitor
func (pm *PerformanceMonitor) SetLogger(logger PDFLogger) {
	pm.logger = logger
	pm.collector.SetLogger(logger)
	pm.analyzer.SetLogger(logger)
	pm.reporter.SetLogger(logger)
}

// Start starts performance monitoring
func (pm *PerformanceMonitor) Start(ctx context.Context) error {
	if !pm.config.EnableMonitoring {
		return nil
	}
	
	pm.logger.Info("Starting performance monitoring")
	
	// Start metrics collection
	go pm.collector.StartCollection(ctx)
	
	// Start analysis
	go pm.analyzer.StartAnalysis(ctx)
	
	// Start reporting
	if pm.config.EnableReporting {
		go pm.reporter.StartReporting(ctx)
	}
	
	return nil
}

// Stop stops performance monitoring
func (pm *PerformanceMonitor) Stop() {
	pm.logger.Info("Stopping performance monitoring")
	// Stop all monitoring goroutines
}

// RecordMetric records a performance metric
func (pm *PerformanceMonitor) RecordMetric(name string, value float64, unit string, tags map[string]string) {
	if !pm.config.EnableMonitoring {
		return
	}
	
	metric := &MetricData{
		Name:        name,
		Value:       value,
		Unit:        unit,
		Timestamp:   time.Now(),
		Tags:        tags,
		Description: pm.getMetricDescription(name),
	}
	
	pm.metrics.RecordMetric(metric)
}

// RecordTiming records a timing metric
func (pm *PerformanceMonitor) RecordTiming(name string, duration time.Duration, tags map[string]string) {
	pm.RecordMetric(name, float64(duration.Nanoseconds()), "ns", tags)
}

// RecordCounter records a counter metric
func (pm *PerformanceMonitor) RecordCounter(name string, value float64, tags map[string]string) {
	pm.RecordMetric(name, value, "count", tags)
}

// RecordGauge records a gauge metric
func (pm *PerformanceMonitor) RecordGauge(name string, value float64, unit string, tags map[string]string) {
	pm.RecordMetric(name, value, unit, tags)
}

// GetMetrics returns current metrics
func (pm *PerformanceMonitor) GetMetrics() map[string]*MetricData {
	return pm.metrics.GetCurrentMetrics()
}

// GetMetricsHistory returns metrics history
func (pm *PerformanceMonitor) GetMetricsHistory() []*MetricsSnapshot {
	return pm.metrics.GetHistory()
}

// GetSystemMetrics returns current system metrics
func (pm *PerformanceMonitor) GetSystemMetrics() *SystemMetrics {
	return pm.collector.CollectSystemMetrics()
}

// GetPerformanceReport generates a performance report
func (pm *PerformanceMonitor) GetPerformanceReport() *PerformanceReport {
	return pm.reporter.GenerateReport()
}

// getMetricDescription returns a description for a metric
func (pm *PerformanceMonitor) getMetricDescription(name string) string {
	descriptions := map[string]string{
		"task_processing_time":    "Time taken to process tasks",
		"layout_generation_time":  "Time taken to generate layout",
		"pdf_generation_time":     "Time taken to generate PDF",
		"memory_usage":           "Memory usage in bytes",
		"cpu_usage":              "CPU usage percentage",
		"cache_hit_rate":         "Cache hit rate percentage",
		"error_count":            "Number of errors",
		"task_count":             "Number of tasks processed",
		"pdf_size":               "Generated PDF size in bytes",
		"compilation_time":       "LaTeX compilation time",
	}
	
	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return "Custom metric"
}

// NewPerformanceMetrics creates a new performance metrics store
func NewPerformanceMetrics(config *PerformanceMonitorConfig) *PerformanceMetrics {
	return &PerformanceMetrics{
		config:  config,
		metrics: make(map[string]*MetricData),
		history: make([]*MetricsSnapshot, 0, config.MaxMetricsHistory),
	}
}

// RecordMetric records a metric
func (pm *PerformanceMetrics) RecordMetric(metric *MetricData) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.metrics[metric.Name] = metric
}

// GetCurrentMetrics returns current metrics
func (pm *PerformanceMetrics) GetCurrentMetrics() map[string]*MetricData {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	metrics := make(map[string]*MetricData)
	for name, metric := range pm.metrics {
		metrics[name] = metric
	}
	
	return metrics
}

// GetHistory returns metrics history
func (pm *PerformanceMetrics) GetHistory() []*MetricsSnapshot {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	history := make([]*MetricsSnapshot, len(pm.history))
	copy(history, pm.history)
	
	return history
}

// AddSnapshot adds a metrics snapshot
func (pm *PerformanceMetrics) AddSnapshot(snapshot *MetricsSnapshot) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.history = append(pm.history, snapshot)
	
	// Trim history if too long
	if len(pm.history) > pm.config.MaxMetricsHistory {
		pm.history = pm.history[1:]
	}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config *PerformanceMonitorConfig) *MetricsCollector {
	return &MetricsCollector{
		config: config,
		logger: &PerformanceMonitorLogger{},
	}
}

// SetLogger sets the logger for the collector
func (mc *MetricsCollector) SetLogger(logger PDFLogger) {
	mc.logger = logger
}

// StartCollection starts metrics collection
func (mc *MetricsCollector) StartCollection(ctx context.Context) {
	ticker := time.NewTicker(mc.config.MetricsInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mc.collectMetrics()
		}
	}
}

// collectMetrics collects current metrics
func (mc *MetricsCollector) collectMetrics() {
	// Collect system metrics
	systemMetrics := mc.CollectSystemMetrics()
	
	// Create snapshot
	snapshot := &MetricsSnapshot{
		Timestamp: time.Now(),
		Metrics:   make(map[string]*MetricData),
		System:    systemMetrics,
	}
	
	// Add system metrics to snapshot
	snapshot.Metrics["memory_usage"] = &MetricData{
		Name:  "memory_usage",
		Value: float64(systemMetrics.MemoryUsage),
		Unit:  "bytes",
		Timestamp: time.Now(),
	}
	
	snapshot.Metrics["cpu_usage"] = &MetricData{
		Name:  "cpu_usage",
		Value: systemMetrics.CPUUsage,
		Unit:  "percent",
		Timestamp: time.Now(),
	}
	
	// Add snapshot to metrics store
	if mc.monitor != nil {
		mc.monitor.metrics.AddSnapshot(snapshot)
	}
}

// CollectSystemMetrics collects current system metrics
func (mc *MetricsCollector) CollectSystemMetrics() *SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Calculate memory usage percentage (rough estimate)
	totalMemory := int64(8 * 1024 * 1024 * 1024) // Assume 8GB total
	memoryPercent := float64(m.Alloc) / float64(totalMemory) * 100
	
	return &SystemMetrics{
		MemoryUsage:    int64(m.Alloc),
		MemoryPercent:  memoryPercent,
		CPUUsage:       0.0, // Would need external library for accurate CPU usage
		GoroutineCount: runtime.NumGoroutine(),
		GCCount:        int64(m.NumGC),
		GCPause:        time.Duration(m.PauseTotalNs),
	}
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer(config *PerformanceMonitorConfig) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		config: config,
		logger: &PerformanceMonitorLogger{},
	}
}

// SetLogger sets the logger for the analyzer
func (pa *PerformanceAnalyzer) SetLogger(logger PDFLogger) {
	pa.logger = logger
}

// StartAnalysis starts performance analysis
func (pa *PerformanceAnalyzer) StartAnalysis(ctx context.Context) {
	ticker := time.NewTicker(pa.config.MetricsInterval * 2)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pa.analyzePerformance()
		}
	}
}

// analyzePerformance analyzes current performance
func (pa *PerformanceAnalyzer) analyzePerformance() {
	if pa.monitor == nil {
		return
	}
	
	// Get current metrics
	metrics := pa.monitor.GetMetrics()
	
	// Check for performance issues
	pa.checkPerformanceThresholds(metrics)
	
	// Analyze trends
	pa.analyzeTrends()
}

// checkPerformanceThresholds checks if performance thresholds are exceeded
func (pa *PerformanceAnalyzer) checkPerformanceThresholds(metrics map[string]*MetricData) {
	// Check memory usage
	if memoryMetric, exists := metrics["memory_usage"]; exists {
		memoryMB := memoryMetric.Value / (1024 * 1024)
		if memoryMB > float64(pa.config.AlertThresholds.MemoryUsage) {
			pa.logger.Error("Memory usage exceeded threshold: %.2f MB > %d MB", 
				memoryMB, pa.config.AlertThresholds.MemoryUsage)
		}
	}
	
	// Check CPU usage
	if cpuMetric, exists := metrics["cpu_usage"]; exists {
		if cpuMetric.Value > pa.config.AlertThresholds.CPUUsage {
			pa.logger.Error("CPU usage exceeded threshold: %.2f%% > %.2f%%", 
				cpuMetric.Value, pa.config.AlertThresholds.CPUUsage)
		}
	}
}

// analyzeTrends analyzes performance trends
func (pa *PerformanceAnalyzer) analyzeTrends() {
	// This would analyze trends in the metrics history
	// For now, just log that analysis is happening
	pa.logger.Debug("Analyzing performance trends")
}

// NewPerformanceReporter creates a new performance reporter
func NewPerformanceReporter(config *PerformanceMonitorConfig) *PerformanceReporter {
	return &PerformanceReporter{
		config: config,
		logger: &PerformanceMonitorLogger{},
	}
}

// SetLogger sets the logger for the reporter
func (pr *PerformanceReporter) SetLogger(logger PDFLogger) {
	pr.logger = logger
}

// StartReporting starts performance reporting
func (pr *PerformanceReporter) StartReporting(ctx context.Context) {
	ticker := time.NewTicker(pr.config.ReportInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			report := pr.GenerateReport()
			pr.logger.Info("Performance Report: %s", report.Summary)
		}
	}
}

// GenerateReport generates a performance report
func (pr *PerformanceReporter) GenerateReport() *PerformanceReport {
	if pr.monitor == nil {
		return &PerformanceReport{Summary: "No data available"}
	}
	
	metrics := pr.monitor.GetMetrics()
	systemMetrics := pr.monitor.GetSystemMetrics()
	
	// Calculate summary statistics
	report := &PerformanceReport{
		Timestamp:     time.Now(),
		Summary:       pr.generateSummary(metrics, systemMetrics),
		Metrics:       metrics,
		SystemMetrics: systemMetrics,
		Recommendations: pr.generateRecommendations(metrics, systemMetrics),
	}
	
	return report
}

// generateSummary generates a summary of performance
func (pr *PerformanceReporter) generateSummary(metrics map[string]*MetricData, systemMetrics *SystemMetrics) string {
	// Count metrics
	metricCount := len(metrics)
	
	// Get key metrics
	memoryMB := float64(systemMetrics.MemoryUsage) / (1024 * 1024)
	cpuPercent := systemMetrics.CPUUsage
	
	summary := fmt.Sprintf("Performance Summary: %d metrics, Memory: %.2f MB, CPU: %.2f%%, Goroutines: %d", 
		metricCount, memoryMB, cpuPercent, systemMetrics.GoroutineCount)
	
	return summary
}

// generateRecommendations generates performance recommendations
func (pr *PerformanceReporter) generateRecommendations(metrics map[string]*MetricData, systemMetrics *SystemMetrics) []string {
	var recommendations []string
	
	// Check memory usage
	memoryMB := float64(systemMetrics.MemoryUsage) / (1024 * 1024)
	if memoryMB > float64(pr.config.AlertThresholds.MemoryUsage) {
		recommendations = append(recommendations, "Consider optimizing memory usage or increasing memory limits")
	}
	
	// Check CPU usage
	if systemMetrics.CPUUsage > pr.config.AlertThresholds.CPUUsage {
		recommendations = append(recommendations, "Consider optimizing CPU usage or reducing processing load")
	}
	
	// Check goroutine count
	if systemMetrics.GoroutineCount > 1000 {
		recommendations = append(recommendations, "High goroutine count detected, consider optimizing concurrency")
	}
	
	// Check GC pressure
	if systemMetrics.GCCount > 100 {
		recommendations = append(recommendations, "High GC count detected, consider optimizing memory allocation")
	}
	
	return recommendations
}

// PerformanceReport represents a performance report
type PerformanceReport struct {
	Timestamp        time.Time
	Summary          string
	Metrics          map[string]*MetricData
	SystemMetrics    *SystemMetrics
	Recommendations  []string
}

// PerformanceMonitorLogger provides logging for performance monitor
type PerformanceMonitorLogger struct{}

func (l *PerformanceMonitorLogger) Info(msg string, args ...interface{})  { fmt.Printf("[PERF-INFO] "+msg+"\n", args...) }
func (l *PerformanceMonitorLogger) Error(msg string, args ...interface{}) { fmt.Printf("[PERF-ERROR] "+msg+"\n", args...) }
func (l *PerformanceMonitorLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[PERF-DEBUG] "+msg+"\n", args...) }
