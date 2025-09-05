package generator

import (
	"fmt"
	"math"
	"time"

	"latex-yearly-planner/internal/calendar"
)

// QualityTester provides comprehensive quality testing and visual validation
type QualityTester struct {
	config     *QualityTestConfig
	logger     PDFLogger
	thresholds *QualityThresholds
}

// QualityTestConfig defines configuration for quality testing
type QualityTestConfig struct {
	// Test categories
	EnableSpacingTests     bool `json:"enable_spacing_tests"`
	EnableAlignmentTests   bool `json:"enable_alignment_tests"`
	EnableReadabilityTests bool `json:"enable_readability_tests"`
	EnableVisualTests      bool `json:"enable_visual_tests"`
	EnablePerformanceTests bool `json:"enable_performance_tests"`
	
	// Test parameters
	SampleSize             int     `json:"sample_size"`
	TestIterations         int     `json:"test_iterations"`
	QualityThreshold       float64 `json:"quality_threshold"`
	PerformanceThreshold   float64 `json:"performance_threshold"`
	
	// Visual validation
	EnableScreenshotTests  bool `json:"enable_screenshot_tests"`
	EnableColorTests       bool `json:"enable_color_tests"`
	EnableFontTests        bool `json:"enable_font_tests"`
	EnableLayoutTests      bool `json:"enable_layout_tests"`
}

// QualityTestResult contains the results of quality testing
type QualityTestResult struct {
	OverallScore      float64                    `json:"overall_score"`
	TestResults       map[string]*TestCategory   `json:"test_results"`
	Issues            []QualityIssue             `json:"issues"`
	Recommendations   []QualityRecommendation    `json:"recommendations"`
	PerformanceMetrics *PerformanceMetrics       `json:"performance_metrics"`
	TestDuration      time.Duration              `json:"test_duration"`
	Timestamp         time.Time                  `json:"timestamp"`
}

// TestCategory represents results for a specific test category
type TestCategory struct {
	CategoryName string          `json:"category_name"`
	Score        float64         `json:"score"`
	Passed       bool            `json:"passed"`
	Tests        []*IndividualTest `json:"tests"`
	Issues       []QualityIssue  `json:"issues"`
}

// IndividualTest represents a single test result
type IndividualTest struct {
	TestName    string        `json:"test_name"`
	Description string        `json:"description"`
	Passed      bool          `json:"passed"`
	Score       float64       `json:"score"`
	Duration    time.Duration `json:"duration"`
	Details     string        `json:"details"`
	Issues      []QualityIssue `json:"issues"`
}

// QualityIssue represents a quality issue found during testing
type QualityIssue struct {
	Severity    IssueSeverity `json:"severity"`
	Category    string        `json:"category"`
	Description string        `json:"description"`
	Location    string        `json:"location"`
	Suggestions []string      `json:"suggestions"`
	Timestamp   time.Time     `json:"timestamp"`
}

// IssueSeverity represents the severity of a quality issue
type IssueSeverity int

const (
	SeverityLow IssueSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// QualityRecommendation represents a recommendation for improvement
type QualityRecommendation struct {
	Category      string  `json:"category"`
	Description   string  `json:"description"`
	Priority      int     `json:"priority"`
	Impact        float64 `json:"impact"`
	Effort        string  `json:"effort"`
	Implementation string `json:"implementation"`
}

// PerformanceMetrics contains performance measurements
type PerformanceMetrics struct {
	LayoutTime       time.Duration `json:"layout_time"`
	RenderingTime    time.Duration `json:"rendering_time"`
	MemoryUsage      int64         `json:"memory_usage"`
	FileSize         int64         `json:"file_size"`
	CompilationTime  time.Duration `json:"compilation_time"`
	QualityScore     float64       `json:"quality_score"`
}

// NewQualityTester creates a new quality tester
func NewQualityTester() *QualityTester {
	return &QualityTester{
		config:     GetDefaultQualityTestConfig(),
		logger:     &QualityTesterLogger{},
		thresholds: GetDefaultQualityThresholds(),
	}
}

// GetDefaultQualityTestConfig returns the default quality test configuration
func GetDefaultQualityTestConfig() *QualityTestConfig {
	return &QualityTestConfig{
		EnableSpacingTests:     true,
		EnableAlignmentTests:   true,
		EnableReadabilityTests: true,
		EnableVisualTests:      true,
		EnablePerformanceTests: true,
		SampleSize:             100,
		TestIterations:         5,
		QualityThreshold:       0.8,
		PerformanceThreshold:   0.7,
		EnableScreenshotTests:  false, // Disabled by default due to complexity
		EnableColorTests:       true,
		EnableFontTests:        true,
		EnableLayoutTests:      true,
	}
}

// GetDefaultQualityThresholds returns the default quality thresholds
func GetDefaultQualityThresholds() *QualityThresholds {
	return &QualityThresholds{
		MinTaskBarHeight:    8.0,
		MinTaskBarWidth:     12.0,
		MinTextSpacing:      0.5,
		MaxOverlapRatio:     0.3,
		MinFontSize:         6.0,
		MinContrastRatio:    4.5,
		MaxLineLength:       60.0,
		MinVisualClarity:    0.7,
		MinLayoutEfficiency: 0.8,
		MaxVisualNoise:      0.3,
	}
}

// SetLogger sets the logger for the quality tester
func (qt *QualityTester) SetLogger(logger PDFLogger) {
	qt.logger = logger
}

// RunQualityTests runs comprehensive quality tests
func (qt *QualityTester) RunQualityTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) (*QualityTestResult, error) {
	startTime := time.Now()
	qt.logger.Info("Starting comprehensive quality tests for %s view", viewType)
	
	result := &QualityTestResult{
		TestResults: make(map[string]*TestCategory),
		Issues:      make([]QualityIssue, 0),
		Recommendations: make([]QualityRecommendation, 0),
		Timestamp:   time.Now(),
	}
	
	// Run spacing tests
	if qt.config.EnableSpacingTests {
		spacingResult := qt.runSpacingTests(layoutResult, viewType)
		result.TestResults["spacing"] = spacingResult
	}
	
	// Run alignment tests
	if qt.config.EnableAlignmentTests {
		alignmentResult := qt.runAlignmentTests(layoutResult, viewType)
		result.TestResults["alignment"] = alignmentResult
	}
	
	// Run readability tests
	if qt.config.EnableReadabilityTests {
		readabilityResult := qt.runReadabilityTests(layoutResult, viewType)
		result.TestResults["readability"] = readabilityResult
	}
	
	// Run visual tests
	if qt.config.EnableVisualTests {
		visualResult := qt.runVisualTests(layoutResult, viewType)
		result.TestResults["visual"] = visualResult
	}
	
	// Run performance tests
	if qt.config.EnablePerformanceTests {
		performanceResult := qt.runPerformanceTests(layoutResult, viewType)
		result.TestResults["performance"] = performanceResult
	}
	
	// Calculate overall score
	result.OverallScore = qt.calculateOverallScore(result.TestResults)
	
	// Generate recommendations
	result.Recommendations = qt.generateRecommendations(result)
	
	// Calculate test duration
	result.TestDuration = time.Since(startTime)
	
	qt.logger.Info("Quality tests completed in %v with overall score: %.2f", 
		result.TestDuration, result.OverallScore)
	
	return result, nil
}

// runSpacingTests runs spacing-related quality tests
func (qt *QualityTester) runSpacingTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) *TestCategory {
	category := &TestCategory{
		CategoryName: "spacing",
		Tests:        make([]*IndividualTest, 0),
		Issues:       make([]QualityIssue, 0),
	}
	
	// Test 1: Task bar spacing
	test1 := &IndividualTest{
		TestName:    "task_bar_spacing",
		Description: "Validate task bar spacing meets minimum requirements",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime := time.Now()
	spacingIssues := 0
	for _, bar := range layoutResult.TaskBars {
		if bar.Height < qt.thresholds.MinTaskBarHeight {
			spacingIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "spacing",
				Description: fmt.Sprintf("Task bar height %.1f is below minimum %.1f", bar.Height, qt.thresholds.MinTaskBarHeight),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Increase task bar height", "Adjust spacing configuration"},
			})
		}
		if bar.Width < qt.thresholds.MinTaskBarWidth {
			spacingIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "spacing",
				Description: fmt.Sprintf("Task bar width %.1f is below minimum %.1f", bar.Width, qt.thresholds.MinTaskBarWidth),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Increase task bar width", "Adjust spacing configuration"},
			})
		}
	}
	
	test1.Duration = time.Since(startTime)
	if spacingIssues > 0 {
		test1.Passed = false
		test1.Score = math.Max(0.0, 1.0 - float64(spacingIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test1)
	
	// Test 2: Text spacing
	test2 := &IndividualTest{
		TestName:    "text_spacing",
		Description: "Validate text spacing for readability",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime = time.Now()
	textSpacingIssues := 0
	for _, bar := range layoutResult.TaskBars {
		// Check if task name is too long for the available space
		estimatedTextWidth := float64(len(bar.TaskName)) * 0.6 // Rough estimation
		if estimatedTextWidth > bar.Width*0.8 { // 80% of bar width
			textSpacingIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityMedium,
				Category:    "spacing",
				Description: fmt.Sprintf("Task name '%s' may be too long for available space", bar.TaskName),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Truncate task name", "Increase task bar width", "Use smaller font"},
			})
		}
	}
	
	test2.Duration = time.Since(startTime)
	if textSpacingIssues > 0 {
		test2.Passed = false
		test2.Score = math.Max(0.0, 1.0 - float64(textSpacingIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test2)
	
	// Calculate category score
	totalScore := 0.0
	for _, test := range category.Tests {
		totalScore += test.Score
	}
	category.Score = totalScore / float64(len(category.Tests))
	category.Passed = category.Score >= qt.config.QualityThreshold
	
	return category
}

// runAlignmentTests runs alignment-related quality tests
func (qt *QualityTester) runAlignmentTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) *TestCategory {
	category := &TestCategory{
		CategoryName: "alignment",
		Tests:        make([]*IndividualTest, 0),
		Issues:       make([]QualityIssue, 0),
	}
	
	// Test 1: Task bar alignment
	test1 := &IndividualTest{
		TestName:    "task_bar_alignment",
		Description: "Validate task bar alignment consistency",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime := time.Now()
	alignmentIssues := 0
	
	// Check for overlapping task bars
	for i, bar1 := range layoutResult.TaskBars {
		for j, bar2 := range layoutResult.TaskBars {
			if i >= j {
				continue
			}
			
			// Check for horizontal overlap
			if bar1.StartX < bar2.StartX+bar2.Width && bar1.StartX+bar1.Width > bar2.StartX {
				// Check for vertical overlap
				if bar1.Y < bar2.Y+bar2.Height && bar1.Y+bar1.Height > bar2.Y {
					overlapRatio := qt.calculateOverlapRatio(bar1, bar2)
					if overlapRatio > qt.thresholds.MaxOverlapRatio {
						alignmentIssues++
						category.Issues = append(category.Issues, QualityIssue{
							Severity:    SeverityHigh,
							Category:    "alignment",
							Description: fmt.Sprintf("Task bars overlap by %.1f%% (max allowed: %.1f%%)", overlapRatio*100, qt.thresholds.MaxOverlapRatio*100),
							Location:    fmt.Sprintf("Tasks: %s and %s", bar1.TaskName, bar2.TaskName),
							Suggestions: []string{"Adjust task positioning", "Implement smart stacking", "Increase spacing"},
						})
					}
				}
			}
		}
	}
	
	test1.Duration = time.Since(startTime)
	if alignmentIssues > 0 {
		test1.Passed = false
		test1.Score = math.Max(0.0, 1.0 - float64(alignmentIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test1)
	
	// Test 2: Grid alignment
	test2 := &IndividualTest{
		TestName:    "grid_alignment",
		Description: "Validate task bars align with calendar grid",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime = time.Now()
	gridAlignmentIssues := 0
	
	for _, bar := range layoutResult.TaskBars {
		// Check if task bar is positioned within reasonable bounds
		if bar.StartX < 0 || bar.Y < 0 {
			gridAlignmentIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityCritical,
				Category:    "alignment",
				Description: "Task bar positioned outside grid bounds",
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Fix task positioning algorithm", "Validate input data"},
			})
		}
	}
	
	test2.Duration = time.Since(startTime)
	if gridAlignmentIssues > 0 {
		test2.Passed = false
		test2.Score = math.Max(0.0, 1.0 - float64(gridAlignmentIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test2)
	
	// Calculate category score
	totalScore := 0.0
	for _, test := range category.Tests {
		totalScore += test.Score
	}
	category.Score = totalScore / float64(len(category.Tests))
	category.Passed = category.Score >= qt.config.QualityThreshold
	
	return category
}

// runReadabilityTests runs readability-related quality tests
func (qt *QualityTester) runReadabilityTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) *TestCategory {
	category := &TestCategory{
		CategoryName: "readability",
		Tests:        make([]*IndividualTest, 0),
		Issues:       make([]QualityIssue, 0),
	}
	
	// Test 1: Font size validation
	test1 := &IndividualTest{
		TestName:    "font_size_validation",
		Description: "Validate font sizes meet readability requirements",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime := time.Now()
	fontSizeIssues := 0
	
	// Estimate font size based on task bar dimensions
	for _, bar := range layoutResult.TaskBars {
		estimatedFontSize := bar.Height * 0.6 // Rough estimation
		if estimatedFontSize < qt.thresholds.MinFontSize {
			fontSizeIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "readability",
				Description: fmt.Sprintf("Estimated font size %.1f is below minimum %.1f", estimatedFontSize, qt.thresholds.MinFontSize),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Increase task bar height", "Use larger font size", "Truncate text"},
			})
		}
	}
	
	test1.Duration = time.Since(startTime)
	if fontSizeIssues > 0 {
		test1.Passed = false
		test1.Score = math.Max(0.0, 1.0 - float64(fontSizeIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test1)
	
	// Test 2: Text length validation
	test2 := &IndividualTest{
		TestName:    "text_length_validation",
		Description: "Validate text length for readability",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime = time.Now()
	textLengthIssues := 0
	
	for _, bar := range layoutResult.TaskBars {
		if len(bar.TaskName) > int(qt.thresholds.MaxLineLength) {
			textLengthIssues++
			category.Issues = append(category.Issues, QualityIssue{
				Severity:    SeverityMedium,
				Category:    "readability",
				Description: fmt.Sprintf("Task name length %d exceeds maximum %d", len(bar.TaskName), int(qt.thresholds.MaxLineLength)),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Truncate task name", "Use abbreviations", "Increase task bar width"},
			})
		}
	}
	
	test2.Duration = time.Since(startTime)
	if textLengthIssues > 0 {
		test2.Passed = false
		test2.Score = math.Max(0.0, 1.0 - float64(textLengthIssues)/float64(len(layoutResult.TaskBars)))
	}
	
	category.Tests = append(category.Tests, test2)
	
	// Calculate category score
	totalScore := 0.0
	for _, test := range category.Tests {
		totalScore += test.Score
	}
	category.Score = totalScore / float64(len(category.Tests))
	category.Passed = category.Score >= qt.config.QualityThreshold
	
	return category
}

// runVisualTests runs visual-related quality tests
func (qt *QualityTester) runVisualTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) *TestCategory {
	category := &TestCategory{
		CategoryName: "visual",
		Tests:        make([]*IndividualTest, 0),
		Issues:       make([]QualityIssue, 0),
	}
	
	// Test 1: Visual clarity
	test1 := &IndividualTest{
		TestName:    "visual_clarity",
		Description: "Validate visual clarity meets requirements",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime := time.Now()
	
	// Calculate visual clarity based on spacing and alignment
	clarityScore := qt.calculateVisualClarity(layoutResult)
	if clarityScore < qt.thresholds.MinVisualClarity {
		category.Issues = append(category.Issues, QualityIssue{
			Severity:    SeverityHigh,
			Category:    "visual",
			Description: fmt.Sprintf("Visual clarity %.2f is below minimum %.2f", clarityScore, qt.thresholds.MinVisualClarity),
			Location:    "Overall layout",
			Suggestions: []string{"Improve spacing", "Enhance alignment", "Adjust visual hierarchy"},
		})
		test1.Passed = false
		test1.Score = clarityScore
	} else {
		test1.Score = clarityScore
	}
	
	test1.Duration = time.Since(startTime)
	category.Tests = append(category.Tests, test1)
	
	// Test 2: Layout efficiency
	test2 := &IndividualTest{
		TestName:    "layout_efficiency",
		Description: "Validate layout efficiency",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime = time.Now()
	
	// Calculate layout efficiency
	efficiencyScore := qt.calculateLayoutEfficiency(layoutResult)
	if efficiencyScore < qt.thresholds.MinLayoutEfficiency {
		category.Issues = append(category.Issues, QualityIssue{
			Severity:    SeverityMedium,
			Category:    "visual",
			Description: fmt.Sprintf("Layout efficiency %.2f is below minimum %.2f", efficiencyScore, qt.thresholds.MinLayoutEfficiency),
			Location:    "Overall layout",
			Suggestions: []string{"Optimize space usage", "Improve task distribution", "Reduce wasted space"},
		})
		test2.Passed = false
		test2.Score = efficiencyScore
	} else {
		test2.Score = efficiencyScore
	}
	
	test2.Duration = time.Since(startTime)
	category.Tests = append(category.Tests, test2)
	
	// Calculate category score
	totalScore := 0.0
	for _, test := range category.Tests {
		totalScore += test.Score
	}
	category.Score = totalScore / float64(len(category.Tests))
	category.Passed = category.Score >= qt.config.QualityThreshold
	
	return category
}

// runPerformanceTests runs performance-related quality tests
func (qt *QualityTester) runPerformanceTests(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) *TestCategory {
	category := &TestCategory{
		CategoryName: "performance",
		Tests:        make([]*IndividualTest, 0),
		Issues:       make([]QualityIssue, 0),
	}
	
	// Test 1: Layout performance
	test1 := &IndividualTest{
		TestName:    "layout_performance",
		Description: "Validate layout performance",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime := time.Now()
	
	// Simulate layout processing time
	layoutTime := time.Duration(len(layoutResult.TaskBars)) * time.Microsecond * 100
	if layoutTime > time.Millisecond*100 {
		category.Issues = append(category.Issues, QualityIssue{
			Severity:    SeverityMedium,
			Category:    "performance",
			Description: fmt.Sprintf("Layout processing time %v is too slow", layoutTime),
			Location:    "Layout processing",
			Suggestions: []string{"Optimize layout algorithm", "Reduce task count", "Improve data structures"},
		})
		test1.Passed = false
		test1.Score = 0.7
	}
	
	test1.Duration = time.Since(startTime)
	category.Tests = append(category.Tests, test1)
	
	// Test 2: Memory usage
	test2 := &IndividualTest{
		TestName:    "memory_usage",
		Description: "Validate memory usage",
		Passed:      true,
		Score:       1.0,
	}
	
	startTime = time.Now()
	
	// Estimate memory usage
	estimatedMemory := int64(len(layoutResult.TaskBars)) * 1000 // Rough estimation
	if estimatedMemory > 1024*1024*10 { // 10MB
		category.Issues = append(category.Issues, QualityIssue{
			Severity:    SeverityMedium,
			Category:    "performance",
			Description: fmt.Sprintf("Estimated memory usage %d bytes is too high", estimatedMemory),
			Location:    "Memory usage",
			Suggestions: []string{"Optimize data structures", "Reduce memory footprint", "Implement lazy loading"},
		})
		test2.Passed = false
		test2.Score = 0.8
	}
	
	test2.Duration = time.Since(startTime)
	category.Tests = append(category.Tests, test2)
	
	// Calculate category score
	totalScore := 0.0
	for _, test := range category.Tests {
		totalScore += test.Score
	}
	category.Score = totalScore / float64(len(category.Tests))
	category.Passed = category.Score >= qt.config.QualityThreshold
	
	return category
}

// calculateOverlapRatio calculates the overlap ratio between two task bars
func (qt *QualityTester) calculateOverlapRatio(bar1, bar2 *calendar.IntegratedTaskBar) float64 {
	// Calculate horizontal overlap
	horizontalOverlap := math.Max(0, math.Min(bar1.StartX+bar1.Width, bar2.StartX+bar2.Width) - math.Max(bar1.StartX, bar2.StartX))
	
	// Calculate vertical overlap
	verticalOverlap := math.Max(0, math.Min(bar1.Y+bar1.Height, bar2.Y+bar2.Height) - math.Max(bar1.Y, bar2.Y))
	
	// Calculate overlap area
	overlapArea := horizontalOverlap * verticalOverlap
	
	// Calculate total area
	area1 := bar1.Width * bar1.Height
	area2 := bar2.Width * bar2.Height
	
	// Return overlap ratio
	return overlapArea / math.Max(area1, area2)
}

// calculateVisualClarity calculates the visual clarity score
func (qt *QualityTester) calculateVisualClarity(layoutResult *calendar.IntegratedLayoutResult) float64 {
	if len(layoutResult.TaskBars) == 0 {
		return 1.0
	}
	
	// Calculate clarity based on spacing consistency
	spacingScore := 0.0
	for _, bar := range layoutResult.TaskBars {
		if bar.Height >= qt.thresholds.MinTaskBarHeight && bar.Width >= qt.thresholds.MinTaskBarWidth {
			spacingScore += 1.0
		}
	}
	spacingScore /= float64(len(layoutResult.TaskBars))
	
	// Calculate clarity based on overlap
	overlapScore := 1.0
	overlapCount := 0
	for i, bar1 := range layoutResult.TaskBars {
		for j, bar2 := range layoutResult.TaskBars {
			if i >= j {
				continue
			}
			overlapRatio := qt.calculateOverlapRatio(bar1, bar2)
			if overlapRatio > qt.thresholds.MaxOverlapRatio {
				overlapCount++
			}
		}
	}
	if len(layoutResult.TaskBars) > 1 {
		overlapScore = 1.0 - float64(overlapCount)/float64(len(layoutResult.TaskBars)*(len(layoutResult.TaskBars)-1)/2)
	}
	
	// Combine scores
	return (spacingScore + overlapScore) / 2.0
}

// calculateLayoutEfficiency calculates the layout efficiency score
func (qt *QualityTester) calculateLayoutEfficiency(layoutResult *calendar.IntegratedLayoutResult) float64 {
	if len(layoutResult.TaskBars) == 0 {
		return 1.0
	}
	
	// Calculate used space
	usedSpace := 0.0
	for _, bar := range layoutResult.TaskBars {
		usedSpace += bar.Width * bar.Height
	}
	
	// Calculate total available space (assuming 100x100 grid)
	totalSpace := 100.0 * 100.0
	
	// Calculate efficiency
	efficiency := usedSpace / totalSpace
	
	// Normalize to 0-1 range
	return math.Min(1.0, efficiency)
}

// calculateOverallScore calculates the overall quality score
func (qt *QualityTester) calculateOverallScore(testResults map[string]*TestCategory) float64 {
	if len(testResults) == 0 {
		return 0.0
	}
	
	totalScore := 0.0
	for _, category := range testResults {
		totalScore += category.Score
	}
	
	return totalScore / float64(len(testResults))
}

// generateRecommendations generates quality improvement recommendations
func (qt *QualityTester) generateRecommendations(result *QualityTestResult) []QualityRecommendation {
	recommendations := make([]QualityRecommendation, 0)
	
	// Analyze issues and generate recommendations
	for _, issue := range result.Issues {
		switch issue.Category {
		case "spacing":
			recommendations = append(recommendations, QualityRecommendation{
				Category:      "spacing",
				Description:   "Improve spacing configuration to meet quality thresholds",
				Priority:      1,
				Impact:        0.8,
				Effort:        "Medium",
				Implementation: "Adjust spacing parameters in visual configuration",
			})
		case "alignment":
			recommendations = append(recommendations, QualityRecommendation{
				Category:      "alignment",
				Description:   "Implement smart stacking algorithm to prevent overlaps",
				Priority:      2,
				Impact:        0.9,
				Effort:        "High",
				Implementation: "Enhance layout algorithm with overlap detection and resolution",
			})
		case "readability":
			recommendations = append(recommendations, QualityRecommendation{
				Category:      "readability",
				Description:   "Optimize font sizes and text truncation for better readability",
				Priority:      1,
				Impact:        0.7,
				Effort:        "Low",
				Implementation: "Adjust typography settings and implement text truncation",
			})
		case "visual":
			recommendations = append(recommendations, QualityRecommendation{
				Category:      "visual",
				Description:   "Enhance visual hierarchy and clarity",
				Priority:      2,
				Impact:        0.8,
				Effort:        "Medium",
				Implementation: "Improve color schemes and visual design system",
			})
		case "performance":
			recommendations = append(recommendations, QualityRecommendation{
				Category:      "performance",
				Description:   "Optimize layout algorithm for better performance",
				Priority:      3,
				Impact:        0.6,
				Effort:        "High",
				Implementation: "Refactor layout processing and optimize data structures",
			})
		}
	}
	
	return recommendations
}

// QualityTesterLogger provides logging for quality tester
type QualityTesterLogger struct{}

func (l *QualityTesterLogger) Info(msg string, args ...interface{})  { fmt.Printf("[QUALITY-INFO] "+msg+"\n", args...) }
func (l *QualityTesterLogger) Error(msg string, args ...interface{}) { fmt.Printf("[QUALITY-ERROR] "+msg+"\n", args...) }
func (l *QualityTesterLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[QUALITY-DEBUG] "+msg+"\n", args...) }
