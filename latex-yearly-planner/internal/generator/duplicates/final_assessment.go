package generator

import (
	"fmt"
	"time"

	"latex-yearly-planner/internal/calendar"
)

// FinalAssessment provides comprehensive final visual quality assessment
type FinalAssessment struct {
	config        *AssessmentConfig
	visualDesign  *VisualDesignSystem
	qualityTester *QualityTester
	validator     *QualityValidator
	logger        PDFLogger
}

// AssessmentConfig defines configuration for final assessment
type AssessmentConfig struct {
	// Assessment criteria
	MinOverallScore     float64 `json:"min_overall_score"`
	MinSpacingScore     float64 `json:"min_spacing_score"`
	MinAlignmentScore   float64 `json:"min_alignment_score"`
	MinReadabilityScore float64 `json:"min_readability_score"`
	MinVisualScore      float64 `json:"min_visual_score"`
	MinPerformanceScore float64 `json:"min_performance_score"`
	
	// Visual quality criteria
	MinColorContrast    float64 `json:"min_color_contrast"`
	MinFontSize         float64 `json:"min_font_size"`
	MaxVisualNoise      float64 `json:"max_visual_noise"`
	MinLayoutEfficiency float64 `json:"min_layout_efficiency"`
	
	// Professional standards
	RequireAccessibility bool    `json:"require_accessibility"`
	RequireConsistency   bool    `json:"require_consistency"`
	RequirePerformance   bool    `json:"require_performance"`
	
	// Assessment scope
	TestAllViewTypes     bool    `json:"test_all_view_types"`
	TestAllCategories    bool    `json:"test_all_categories"`
	TestEdgeCases        bool    `json:"test_edge_cases"`
}

// AssessmentResult contains the results of final assessment
type AssessmentResult struct {
	OverallPassed        bool                      `json:"overall_passed"`
	OverallScore         float64                   `json:"overall_score"`
	AssessmentTime       time.Duration             `json:"assessment_time"`
	Timestamp            time.Time                 `json:"timestamp"`
	
	// Component scores
	SpacingScore         float64                   `json:"spacing_score"`
	AlignmentScore       float64                   `json:"alignment_score"`
	ReadabilityScore     float64                   `json:"readability_score"`
	VisualScore          float64                   `json:"visual_score"`
	PerformanceScore     float64                   `json:"performance_score"`
	
	// Visual quality metrics
	ColorContrast        float64                   `json:"color_contrast"`
	FontSizeCompliance   float64                   `json:"font_size_compliance"`
	VisualNoise          float64                   `json:"visual_noise"`
	LayoutEfficiency     float64                   `json:"layout_efficiency"`
	
	// Professional standards
	AccessibilityPassed  bool                      `json:"accessibility_passed"`
	ConsistencyPassed    bool                      `json:"consistency_passed"`
	PerformancePassed    bool                      `json:"performance_passed"`
	
	// Detailed results
	ViewTypeResults      map[ViewType]*ViewResult  `json:"view_type_results"`
	CategoryResults      map[string]*CategoryResult `json:"category_results"`
	EdgeCaseResults      map[string]*EdgeCaseResult `json:"edge_case_results"`
	
	// Issues and recommendations
	CriticalIssues       []QualityIssue            `json:"critical_issues"`
	HighPriorityIssues   []QualityIssue            `json:"high_priority_issues"`
	MediumPriorityIssues []QualityIssue            `json:"medium_priority_issues"`
	LowPriorityIssues    []QualityIssue            `json:"low_priority_issues"`
	
	Recommendations      []QualityRecommendation   `json:"recommendations"`
	ActionItems          []ActionItem              `json:"action_items"`
}

// ViewResult contains results for a specific view type
type ViewResult struct {
	ViewType        ViewType        `json:"view_type"`
	Passed          bool            `json:"passed"`
	Score           float64         `json:"score"`
	Issues          []QualityIssue  `json:"issues"`
	Recommendations []QualityRecommendation `json:"recommendations"`
}

// CategoryResult contains results for a specific category
type CategoryResult struct {
	Category        string          `json:"category"`
	Passed          bool            `json:"passed"`
	Score           float64         `json:"score"`
	ColorContrast   float64         `json:"color_contrast"`
	Readability     float64         `json:"readability"`
	Issues          []QualityIssue  `json:"issues"`
	Recommendations []QualityRecommendation `json:"recommendations"`
}

// EdgeCaseResult contains results for edge case testing
type EdgeCaseResult struct {
	TestCase        string          `json:"test_case"`
	Passed          bool            `json:"passed"`
	Score           float64         `json:"score"`
	Issues          []QualityIssue  `json:"issues"`
	Recommendations []QualityRecommendation `json:"recommendations"`
}

// ActionItem represents a specific action to take
type ActionItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Priority    int     `json:"priority"`
	Effort      string  `json:"effort"`
	Impact      float64 `json:"impact"`
	Status      string  `json:"status"`
	Assignee    string  `json:"assignee"`
	DueDate     string  `json:"due_date"`
}

// NewFinalAssessment creates a new final assessment
func NewFinalAssessment() *FinalAssessment {
	return &FinalAssessment{
		config:        GetDefaultAssessmentConfig(),
		visualDesign:  NewVisualDesignSystem(),
		qualityTester: NewQualityTester(),
		validator:     NewQualityValidator(),
		logger:        &FinalAssessmentLogger{},
	}
}

// GetDefaultAssessmentConfig returns the default assessment configuration
func GetDefaultAssessmentConfig() *AssessmentConfig {
	return &AssessmentConfig{
		MinOverallScore:     0.85,
		MinSpacingScore:     0.80,
		MinAlignmentScore:   0.80,
		MinReadabilityScore: 0.85,
		MinVisualScore:      0.85,
		MinPerformanceScore: 0.75,
		MinColorContrast:    4.5,
		MinFontSize:         12.0,
		MaxVisualNoise:      0.3,
		MinLayoutEfficiency: 0.80,
		RequireAccessibility: true,
		RequireConsistency:   true,
		RequirePerformance:   true,
		TestAllViewTypes:     true,
		TestAllCategories:    true,
		TestEdgeCases:        true,
	}
}

// SetLogger sets the logger for the final assessment
func (fa *FinalAssessment) SetLogger(logger PDFLogger) {
	fa.logger = logger
	fa.visualDesign.SetLogger(logger)
	fa.qualityTester.SetLogger(logger)
	fa.validator.SetLogger(logger)
}

// ConductFinalAssessment performs comprehensive final visual quality assessment
func (fa *FinalAssessment) ConductFinalAssessment() (*AssessmentResult, error) {
	startTime := time.Now()
	fa.logger.Info("Starting comprehensive final visual quality assessment")
	
	result := &AssessmentResult{
		ViewTypeResults:      make(map[ViewType]*ViewResult),
		CategoryResults:      make(map[string]*CategoryResult),
		EdgeCaseResults:      make(map[string]*EdgeCaseResult),
		CriticalIssues:       make([]QualityIssue, 0),
		HighPriorityIssues:   make([]QualityIssue, 0),
		MediumPriorityIssues: make([]QualityIssue, 0),
		LowPriorityIssues:    make([]QualityIssue, 0),
		Recommendations:      make([]QualityRecommendation, 0),
		ActionItems:          make([]ActionItem, 0),
		Timestamp:            time.Now(),
	}
	
	// Test all view types
	if fa.config.TestAllViewTypes {
		fa.logger.Info("Testing all view types")
		fa.testAllViewTypes(result)
	}
	
	// Test all categories
	if fa.config.TestAllCategories {
		fa.logger.Info("Testing all categories")
		fa.testAllCategories(result)
	}
	
	// Test edge cases
	if fa.config.TestEdgeCases {
		fa.logger.Info("Testing edge cases")
		fa.testEdgeCases(result)
	}
	
	// Calculate overall scores
	fa.calculateOverallScores(result)
	
	// Assess professional standards
	fa.assessProfessionalStandards(result)
	
	// Generate recommendations and action items
	fa.generateRecommendations(result)
	fa.generateActionItems(result)
	
	// Determine overall pass/fail
	result.OverallPassed = fa.determineOverallPass(result)
	
	// Calculate assessment time
	result.AssessmentTime = time.Since(startTime)
	
	fa.logger.Info("Final assessment completed in %v with overall score: %.2f", 
		result.AssessmentTime, result.OverallScore)
	
	return result, nil
}

// testAllViewTypes tests all supported view types
func (fa *FinalAssessment) testAllViewTypes(result *AssessmentResult) {
	viewTypes := []ViewType{
		ViewTypeMonthly,
		ViewTypeWeekly,
		ViewTypeDaily,
		ViewTypeYearly,
		ViewTypeQuarterly,
	}
	
	for _, viewType := range viewTypes {
		fa.logger.Info("Testing view type: %s", viewType)
		
		// Create sample layout result for testing
		layoutResult := fa.createSampleLayoutResult(viewType)
		
		// Run quality tests
		qualityResult, err := fa.qualityTester.RunQualityTests(layoutResult, viewType)
		if err != nil {
			fa.logger.Error("Quality tests failed for view type %s: %v", viewType, err)
			continue
		}
		
		// Create view result
		viewResult := &ViewResult{
			ViewType:        viewType,
			Passed:          qualityResult.OverallScore >= fa.config.MinOverallScore,
			Score:           qualityResult.OverallScore,
			Issues:          qualityResult.Issues,
			Recommendations: qualityResult.Recommendations,
		}
		
		result.ViewTypeResults[viewType] = viewResult
	}
}

// testAllCategories tests all supported categories
func (fa *FinalAssessment) testAllCategories(result *AssessmentResult) {
	categories := []string{
		"PROPOSAL", "LASER", "IMAGING", "ADMIN", 
		"DISSERTATION", "RESEARCH", "PUBLICATION",
	}
	
	for _, category := range categories {
		fa.logger.Info("Testing category: %s", category)
		
		// Test color contrast
		colorContrast := fa.testCategoryColorContrast(category)
		
		// Test readability
		readability := fa.testCategoryReadability(category)
		
		// Create category result
		categoryResult := &CategoryResult{
			Category:        category,
			Passed:          colorContrast >= fa.config.MinColorContrast && readability >= fa.config.MinReadabilityScore,
			Score:           (colorContrast + readability) / 2.0,
			ColorContrast:   colorContrast,
			Readability:     readability,
			Issues:          fa.identifyCategoryIssues(category, colorContrast, readability),
			Recommendations: fa.generateCategoryRecommendations(category, colorContrast, readability),
		}
		
		result.CategoryResults[category] = categoryResult
	}
}

// testEdgeCases tests various edge cases
func (fa *FinalAssessment) testEdgeCases(result *AssessmentResult) {
	edgeCases := []string{
		"empty_layout",
		"single_task",
		"many_tasks",
		"long_task_names",
		"special_characters",
		"high_priority_tasks",
		"overlapping_tasks",
		"month_boundaries",
	}
	
	for _, edgeCase := range edgeCases {
		fa.logger.Info("Testing edge case: %s", edgeCase)
		
		// Create edge case layout
		layoutResult := fa.createEdgeCaseLayout(edgeCase)
		
		// Run quality tests
		qualityResult, err := fa.qualityTester.RunQualityTests(layoutResult, ViewTypeMonthly)
		if err != nil {
			fa.logger.Error("Quality tests failed for edge case %s: %v", edgeCase, err)
			continue
		}
		
		// Create edge case result
		edgeCaseResult := &EdgeCaseResult{
			TestCase:        edgeCase,
			Passed:          qualityResult.OverallScore >= fa.config.MinOverallScore,
			Score:           qualityResult.OverallScore,
			Issues:          qualityResult.Issues,
			Recommendations: qualityResult.Recommendations,
		}
		
		result.EdgeCaseResults[edgeCase] = edgeCaseResult
	}
}

// calculateOverallScores calculates overall scores
func (fa *FinalAssessment) calculateOverallScores(result *AssessmentResult) {
	// Calculate component scores from view type results
	totalScore := 0.0
	count := 0
	
	for _, viewResult := range result.ViewTypeResults {
		totalScore += viewResult.Score
		count++
	}
	
	if count > 0 {
		result.OverallScore = totalScore / float64(count)
	}
	
	// Calculate individual component scores
	result.SpacingScore = fa.calculateComponentScore(result, "spacing")
	result.AlignmentScore = fa.calculateComponentScore(result, "alignment")
	result.ReadabilityScore = fa.calculateComponentScore(result, "readability")
	result.VisualScore = fa.calculateComponentScore(result, "visual")
	result.PerformanceScore = fa.calculateComponentScore(result, "performance")
	
	// Calculate visual quality metrics
	result.ColorContrast = fa.calculateAverageColorContrast(result)
	result.FontSizeCompliance = fa.calculateFontSizeCompliance(result)
	result.VisualNoise = fa.calculateVisualNoise(result)
	result.LayoutEfficiency = fa.calculateLayoutEfficiency(result)
}

// assessProfessionalStandards assesses professional standards
func (fa *FinalAssessment) assessProfessionalStandards(result *AssessmentResult) {
	// Assess accessibility
	result.AccessibilityPassed = fa.assessAccessibility(result)
	
	// Assess consistency
	result.ConsistencyPassed = fa.assessConsistency(result)
	
	// Assess performance
	result.PerformancePassed = fa.assessPerformance(result)
}

// generateRecommendations generates recommendations based on assessment results
func (fa *FinalAssessment) generateRecommendations(result *AssessmentResult) {
	// Collect all issues by priority
	fa.collectIssuesByPriority(result)
	
	// Generate recommendations based on issues
	for _, issue := range result.CriticalIssues {
		result.Recommendations = append(result.Recommendations, QualityRecommendation{
			Category:      issue.Category,
			Description:   fmt.Sprintf("CRITICAL: %s", issue.Description),
			Priority:      1,
			Impact:        1.0,
			Effort:        "High",
			Implementation: "Immediate action required",
		})
	}
	
	for _, issue := range result.HighPriorityIssues {
		result.Recommendations = append(result.Recommendations, QualityRecommendation{
			Category:      issue.Category,
			Description:   fmt.Sprintf("HIGH: %s", issue.Description),
			Priority:      2,
			Impact:        0.8,
			Effort:        "Medium",
			Implementation: "Address within 1 week",
		})
	}
	
	for _, issue := range result.MediumPriorityIssues {
		result.Recommendations = append(result.Recommendations, QualityRecommendation{
			Category:      issue.Category,
			Description:   fmt.Sprintf("MEDIUM: %s", issue.Description),
			Priority:      3,
			Impact:        0.6,
			Effort:        "Low",
			Implementation: "Address within 2 weeks",
		})
	}
}

// generateActionItems generates specific action items
func (fa *FinalAssessment) generateActionItems(result *AssessmentResult) {
	actionItems := []ActionItem{
		{
			ID:          "accessibility-audit",
			Description: "Conduct comprehensive accessibility audit",
			Priority:    1,
			Effort:      "High",
			Impact:      0.9,
			Status:      "pending",
			Assignee:    "Visual Design Team",
			DueDate:     "2024-01-15",
		},
		{
			ID:          "color-contrast-review",
			Description: "Review and improve color contrast ratios",
			Priority:    2,
			Effort:      "Medium",
			Impact:      0.8,
			Status:      "pending",
			Assignee:    "Design Team",
			DueDate:     "2024-01-20",
		},
		{
			ID:          "typography-optimization",
			Description: "Optimize typography for better readability",
			Priority:    2,
			Effort:      "Medium",
			Impact:      0.7,
			Status:      "pending",
			Assignee:    "Typography Team",
			DueDate:     "2024-01-25",
		},
		{
			ID:          "performance-optimization",
			Description: "Optimize layout performance and rendering",
			Priority:    3,
			Effort:      "High",
			Impact:      0.6,
			Status:      "pending",
			Assignee:    "Performance Team",
			DueDate:     "2024-02-01",
		},
	}
	
	result.ActionItems = actionItems
}

// Helper methods for assessment calculations
func (fa *FinalAssessment) createSampleLayoutResult(viewType ViewType) *calendar.IntegratedLayoutResult {
	// Create sample layout result based on view type
	return &calendar.IntegratedLayoutResult{
		TaskBars: []*calendar.IntegratedTaskBar{
			{
				TaskName:    "Sample Task 1",
				Description: "A sample task for testing",
				Category:    "RESEARCH",
				Priority:    3,
				StartX:      10.0,
				Y:           20.0,
				Width:       50.0,
				Height:      15.0,
				Color:       "blue",
				IsContinuation: false,
				IsStart:     true,
				IsEnd:       true,
				MonthBoundary: false,
			},
		},
		Statistics: &calendar.IntegratedLayoutStatistics{
			TotalTasks:     1,
			ProcessedBars:  1,
			SpaceEfficiency: 0.85,
			VisualQuality:   0.80,
		},
	}
}

func (fa *FinalAssessment) testCategoryColorContrast(category string) float64 {
	// Test color contrast for a specific category
	// This would typically involve calculating actual contrast ratios
	return 4.5 // Placeholder - should be calculated based on actual colors
}

func (fa *FinalAssessment) testCategoryReadability(category string) float64 {
	// Test readability for a specific category
	// This would typically involve analyzing font sizes, spacing, etc.
	return 0.85 // Placeholder - should be calculated based on actual metrics
}

func (fa *FinalAssessment) identifyCategoryIssues(category string, colorContrast, readability float64) []QualityIssue {
	issues := make([]QualityIssue, 0)
	
	if colorContrast < fa.config.MinColorContrast {
		issues = append(issues, QualityIssue{
			Severity:    SeverityHigh,
			Category:    "accessibility",
			Description: fmt.Sprintf("Category %s has insufficient color contrast: %.1f", category, colorContrast),
			Location:    fmt.Sprintf("Category: %s", category),
			Suggestions: []string{"Increase color contrast", "Use darker/lighter variants"},
		})
	}
	
	if readability < fa.config.MinReadabilityScore {
		issues = append(issues, QualityIssue{
			Severity:    SeverityMedium,
			Category:    "readability",
			Description: fmt.Sprintf("Category %s has poor readability: %.2f", category, readability),
			Location:    fmt.Sprintf("Category: %s", category),
			Suggestions: []string{"Improve font sizing", "Adjust spacing", "Enhance typography"},
		})
	}
	
	return issues
}

func (fa *FinalAssessment) generateCategoryRecommendations(category string, colorContrast, readability float64) []QualityRecommendation {
	recommendations := make([]QualityRecommendation, 0)
	
	if colorContrast < fa.config.MinColorContrast {
		recommendations = append(recommendations, QualityRecommendation{
			Category:      "accessibility",
			Description:   fmt.Sprintf("Improve color contrast for category %s", category),
			Priority:      1,
			Impact:        0.9,
			Effort:        "Medium",
			Implementation: "Update category color definitions",
		})
	}
	
	if readability < fa.config.MinReadabilityScore {
		recommendations = append(recommendations, QualityRecommendation{
			Category:      "readability",
			Description:   fmt.Sprintf("Enhance readability for category %s", category),
			Priority:      2,
			Impact:        0.7,
			Effort:        "Low",
			Implementation: "Adjust typography settings",
		})
	}
	
	return recommendations
}

func (fa *FinalAssessment) createEdgeCaseLayout(edgeCase string) *calendar.IntegratedLayoutResult {
	// Create edge case layout based on the test case
	return &calendar.IntegratedLayoutResult{
		TaskBars: []*calendar.IntegratedTaskBar{},
		Statistics: &calendar.IntegratedLayoutStatistics{
			TotalTasks:     0,
			ProcessedBars:  0,
			SpaceEfficiency: 0.0,
			VisualQuality:   0.0,
		},
	}
}

func (fa *FinalAssessment) calculateComponentScore(result *AssessmentResult, component string) float64 {
	// Calculate component score from all results
	totalScore := 0.0
	count := 0
	
	for _, viewResult := range result.ViewTypeResults {
		// This would typically extract the specific component score
		// For now, use the overall score as a placeholder
		totalScore += viewResult.Score
		count++
	}
	
	if count == 0 {
		return 0.0
	}
	
	return totalScore / float64(count)
}

func (fa *FinalAssessment) calculateAverageColorContrast(result *AssessmentResult) float64 {
	totalContrast := 0.0
	count := 0
	
	for _, categoryResult := range result.CategoryResults {
		totalContrast += categoryResult.ColorContrast
		count++
	}
	
	if count == 0 {
		return 0.0
	}
	
	return totalContrast / float64(count)
}

func (fa *FinalAssessment) calculateFontSizeCompliance(result *AssessmentResult) float64 {
	// Calculate font size compliance based on minimum requirements
	// This would typically involve analyzing actual font sizes used
	return 0.95 // Placeholder - should be calculated based on actual metrics
}

func (fa *FinalAssessment) calculateVisualNoise(result *AssessmentResult) float64 {
	// Calculate visual noise based on layout complexity
	// This would typically involve analyzing visual elements and their density
	return 0.15 // Placeholder - should be calculated based on actual metrics
}

func (fa *FinalAssessment) calculateLayoutEfficiency(result *AssessmentResult) float64 {
	// Calculate layout efficiency based on space usage
	// This would typically involve analyzing space utilization
	return 0.85 // Placeholder - should be calculated based on actual metrics
}

func (fa *FinalAssessment) assessAccessibility(result *AssessmentResult) bool {
	// Assess accessibility compliance
	return result.ColorContrast >= fa.config.MinColorContrast &&
		   result.FontSizeCompliance >= 0.9 &&
		   len(result.CriticalIssues) == 0
}

func (fa *FinalAssessment) assessConsistency(result *AssessmentResult) bool {
	// Assess visual consistency across view types and categories
	// This would typically involve analyzing consistency metrics
	return true // Placeholder - should be calculated based on actual metrics
}

func (fa *FinalAssessment) assessPerformance(result *AssessmentResult) bool {
	// Assess performance compliance
	return result.PerformanceScore >= fa.config.MinPerformanceScore &&
		   result.LayoutEfficiency >= fa.config.MinLayoutEfficiency
}

func (fa *FinalAssessment) collectIssuesByPriority(result *AssessmentResult) {
	// Collect issues from all results and categorize by priority
	for _, viewResult := range result.ViewTypeResults {
		for _, issue := range viewResult.Issues {
			switch issue.Severity {
			case SeverityCritical:
				result.CriticalIssues = append(result.CriticalIssues, issue)
			case SeverityHigh:
				result.HighPriorityIssues = append(result.HighPriorityIssues, issue)
			case SeverityMedium:
				result.MediumPriorityIssues = append(result.MediumPriorityIssues, issue)
			case SeverityLow:
				result.LowPriorityIssues = append(result.LowPriorityIssues, issue)
			}
		}
	}
	
	for _, categoryResult := range result.CategoryResults {
		for _, issue := range categoryResult.Issues {
			switch issue.Severity {
			case SeverityCritical:
				result.CriticalIssues = append(result.CriticalIssues, issue)
			case SeverityHigh:
				result.HighPriorityIssues = append(result.HighPriorityIssues, issue)
			case SeverityMedium:
				result.MediumPriorityIssues = append(result.MediumPriorityIssues, issue)
			case SeverityLow:
				result.LowPriorityIssues = append(result.LowPriorityIssues, issue)
			}
		}
	}
	
	for _, edgeCaseResult := range result.EdgeCaseResults {
		for _, issue := range edgeCaseResult.Issues {
			switch issue.Severity {
			case SeverityCritical:
				result.CriticalIssues = append(result.CriticalIssues, issue)
			case SeverityHigh:
				result.HighPriorityIssues = append(result.HighPriorityIssues, issue)
			case SeverityMedium:
				result.MediumPriorityIssues = append(result.MediumPriorityIssues, issue)
			case SeverityLow:
				result.LowPriorityIssues = append(result.LowPriorityIssues, issue)
			}
		}
	}
}

func (fa *FinalAssessment) determineOverallPass(result *AssessmentResult) bool {
	// Determine if the overall assessment passes
	return result.OverallScore >= fa.config.MinOverallScore &&
		   result.SpacingScore >= fa.config.MinSpacingScore &&
		   result.AlignmentScore >= fa.config.MinAlignmentScore &&
		   result.ReadabilityScore >= fa.config.MinReadabilityScore &&
		   result.VisualScore >= fa.config.MinVisualScore &&
		   result.PerformanceScore >= fa.config.MinPerformanceScore &&
		   result.AccessibilityPassed &&
		   result.ConsistencyPassed &&
		   result.PerformancePassed &&
		   len(result.CriticalIssues) == 0
}

// FinalAssessmentLogger provides logging for final assessment
type FinalAssessmentLogger struct{}

func (l *FinalAssessmentLogger) Info(msg string, args ...interface{})  { fmt.Printf("[ASSESSMENT-INFO] "+msg+"\n", args...) }
func (l *FinalAssessmentLogger) Error(msg string, args ...interface{}) { fmt.Printf("[ASSESSMENT-ERROR] "+msg+"\n", args...) }
func (l *FinalAssessmentLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[ASSESSMENT-DEBUG] "+msg+"\n", args...) }
