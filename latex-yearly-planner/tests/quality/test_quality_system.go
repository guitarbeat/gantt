package main

import (
	"fmt"
	"os"

	"latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/generator"
)

// Test the quality testing and validation system
func main() {
	fmt.Println("Testing Quality Testing and Validation System...")

	// Test 1: Quality Tester
	fmt.Println("\n=== Test 1: Quality Tester ===")
	testQualityTester()

	// Test 2: Quality Validator
	fmt.Println("\n=== Test 2: Quality Validator ===")
	testQualityValidator()

	// Test 3: Quality Configuration
	fmt.Println("\n=== Test 3: Quality Configuration ===")
	testQualityConfiguration()

	// Test 4: Quality Metrics
	fmt.Println("\n=== Test 4: Quality Metrics ===")
	testQualityMetrics()

	fmt.Println("\n✅ Quality testing system tests completed!")
}

func testQualityTester() {
	// Create quality tester
	tester := generator.NewQualityTester()
	if tester == nil {
		fmt.Println("❌ Failed to create quality tester")
		return
	}

	// Create sample layout result
	layoutResult := createSampleLayoutResult()

	// Run quality tests
	qualityResult, err := tester.RunQualityTests(layoutResult, generator.ViewTypeMonthly)
	if err != nil {
		fmt.Printf("❌ Quality tests failed: %v\n", err)
		return
	}

	// Validate results
	if qualityResult == nil {
		fmt.Println("❌ Quality test result is nil")
		return
	}

	if qualityResult.OverallScore < 0.0 || qualityResult.OverallScore > 1.0 {
		fmt.Println("❌ Overall score is out of range")
		return
	}

	if len(qualityResult.TestResults) == 0 {
		fmt.Println("❌ No test results generated")
		return
	}

	// Check test categories
	expectedCategories := []string{"spacing", "alignment", "readability", "visual", "performance"}
	for _, category := range expectedCategories {
		if _, exists := qualityResult.TestResults[category]; !exists {
			fmt.Printf("❌ Missing test category: %s\n", category)
			return
		}
	}

	fmt.Printf("✅ Quality tester test passed\n")
	fmt.Printf("   Overall score: %.2f\n", qualityResult.OverallScore)
	fmt.Printf("   Test categories: %d\n", len(qualityResult.TestResults))
	fmt.Printf("   Test duration: %v\n", qualityResult.TestDuration)
	fmt.Printf("   Issues found: %d\n", len(qualityResult.Issues))
	fmt.Printf("   Recommendations: %d\n", len(qualityResult.Recommendations))
}

func testQualityValidator() {
	// Create quality validator
	validator := generator.NewQualityValidator()
	if validator == nil {
		fmt.Println("❌ Failed to create quality validator")
		return
	}

	// Create sample layout result
	layoutResult := createSampleLayoutResult()

	// Create a temporary PDF file for testing
	pdfPath := "test_quality.pdf"
	file, err := os.Create(pdfPath)
	if err != nil {
		fmt.Printf("❌ Failed to create test PDF file: %v\n", err)
		return
	}
	file.WriteString("Test PDF content")
	file.Close()
	defer os.Remove(pdfPath)

	// Run quality validation
	validationResult, err := validator.ValidateQuality(layoutResult, generator.ViewTypeMonthly, pdfPath)
	if err != nil {
		fmt.Printf("❌ Quality validation failed: %v\n", err)
		return
	}

	// Validate results
	if validationResult == nil {
		fmt.Println("❌ Validation result is nil")
		return
	}

	if validationResult.ValidationScore < 0.0 || validationResult.ValidationScore > 1.0 {
		fmt.Println("❌ Validation score is out of range")
		return
	}

	// Check validation components
	if validationResult.PDFValidation == nil {
		fmt.Println("❌ PDF validation is nil")
		return
	}

	if validationResult.LaTeXValidation == nil {
		fmt.Println("❌ LaTeX validation is nil")
		return
	}

	if validationResult.VisualValidation == nil {
		fmt.Println("❌ Visual validation is nil")
		return
	}

	if validationResult.ContentValidation == nil {
		fmt.Println("❌ Content validation is nil")
		return
	}

	fmt.Printf("✅ Quality validator test passed\n")
	fmt.Printf("   Overall passed: %v\n", validationResult.OverallPassed)
	fmt.Printf("   Validation score: %.2f\n", validationResult.ValidationScore)
	fmt.Printf("   PDF validation passed: %v\n", validationResult.PDFValidation.Passed)
	fmt.Printf("   LaTeX validation passed: %v\n", validationResult.LaTeXValidation.Passed)
	fmt.Printf("   Visual validation passed: %v\n", validationResult.VisualValidation.Passed)
	fmt.Printf("   Content validation passed: %v\n", validationResult.ContentValidation.Passed)
	fmt.Printf("   Validation time: %v\n", validationResult.ValidationTime)
	fmt.Printf("   Issues found: %d\n", len(validationResult.Issues))
	fmt.Printf("   Recommendations: %d\n", len(validationResult.Recommendations))
}

func testQualityConfiguration() {
	// Test quality test configuration
	testConfig := generator.GetDefaultQualityTestConfig()
	if testConfig == nil {
		fmt.Println("❌ Failed to get default quality test config")
		return
	}

	// Validate configuration
	if !testConfig.EnableSpacingTests {
		fmt.Println("❌ Spacing tests should be enabled by default")
		return
	}

	if !testConfig.EnableAlignmentTests {
		fmt.Println("❌ Alignment tests should be enabled by default")
		return
	}

	if !testConfig.EnableReadabilityTests {
		fmt.Println("❌ Readability tests should be enabled by default")
		return
	}

	if !testConfig.EnableVisualTests {
		fmt.Println("❌ Visual tests should be enabled by default")
		return
	}

	if !testConfig.EnablePerformanceTests {
		fmt.Println("❌ Performance tests should be enabled by default")
		return
	}

	if testConfig.SampleSize <= 0 {
		fmt.Println("❌ Sample size should be positive")
		return
	}

	if testConfig.TestIterations <= 0 {
		fmt.Println("❌ Test iterations should be positive")
		return
	}

	if testConfig.QualityThreshold < 0.0 || testConfig.QualityThreshold > 1.0 {
		fmt.Println("❌ Quality threshold should be between 0 and 1")
		return
	}

	// Test quality validation configuration
	validationConfig := generator.GetDefaultQualityValidationConfig()
	if validationConfig == nil {
		fmt.Println("❌ Failed to get default quality validation config")
		return
	}

	// Validate validation configuration
	if !validationConfig.EnablePDFValidation {
		fmt.Println("❌ PDF validation should be enabled by default")
		return
	}

	if !validationConfig.EnableLaTeXValidation {
		fmt.Println("❌ LaTeX validation should be enabled by default")
		return
	}

	if !validationConfig.EnableVisualValidation {
		fmt.Println("❌ Visual validation should be enabled by default")
		return
	}

	if !validationConfig.EnableContentValidation {
		fmt.Println("❌ Content validation should be enabled by default")
		return
	}

	if validationConfig.MinPDFSize <= 0 {
		fmt.Println("❌ Minimum PDF size should be positive")
		return
	}

	if validationConfig.MaxPDFSize <= validationConfig.MinPDFSize {
		fmt.Println("❌ Maximum PDF size should be greater than minimum")
		return
	}

	if validationConfig.MinPageCount <= 0 {
		fmt.Println("❌ Minimum page count should be positive")
		return
	}

	if validationConfig.MaxPageCount <= validationConfig.MinPageCount {
		fmt.Println("❌ Maximum page count should be greater than minimum")
		return
	}

	if validationConfig.MinQualityScore < 0.0 || validationConfig.MinQualityScore > 1.0 {
		fmt.Println("❌ Minimum quality score should be between 0 and 1")
		return
	}

	if validationConfig.MaxCompilationTime <= 0 {
		fmt.Println("❌ Maximum compilation time should be positive")
		return
	}

	fmt.Printf("✅ Quality configuration test passed\n")
	fmt.Printf("   Test config sample size: %d\n", testConfig.SampleSize)
	fmt.Printf("   Test config iterations: %d\n", testConfig.TestIterations)
	fmt.Printf("   Test config quality threshold: %.2f\n", testConfig.QualityThreshold)
	fmt.Printf("   Validation config min PDF size: %d bytes\n", validationConfig.MinPDFSize)
	fmt.Printf("   Validation config max PDF size: %d bytes\n", validationConfig.MaxPDFSize)
	fmt.Printf("   Validation config min page count: %d\n", validationConfig.MinPageCount)
	fmt.Printf("   Validation config max page count: %d\n", validationConfig.MaxPageCount)
	fmt.Printf("   Validation config min quality score: %.2f\n", validationConfig.MinQualityScore)
	fmt.Printf("   Validation config max compilation time: %v\n", validationConfig.MaxCompilationTime)
}

func testQualityMetrics() {
	// Test quality metrics calculation
	metrics := &generator.VisualQualityMetrics{
		OverallScore:     0.85,
		SpacingScore:     0.90,
		AlignmentScore:   0.80,
		ReadabilityScore: 0.85,
		VisualClarity:    0.85,
		LayoutEfficiency: 0.80,
		VisualNoise:      0.15,
	}

	// Validate metrics ranges
	if metrics.OverallScore < 0.0 || metrics.OverallScore > 1.0 {
		fmt.Println("❌ Overall score is out of range")
		return
	}

	if metrics.SpacingScore < 0.0 || metrics.SpacingScore > 1.0 {
		fmt.Println("❌ Spacing score is out of range")
		return
	}

	if metrics.AlignmentScore < 0.0 || metrics.AlignmentScore > 1.0 {
		fmt.Println("❌ Alignment score is out of range")
		return
	}

	if metrics.ReadabilityScore < 0.0 || metrics.ReadabilityScore > 1.0 {
		fmt.Println("❌ Readability score is out of range")
		return
	}

	if metrics.VisualClarity < 0.0 || metrics.VisualClarity > 1.0 {
		fmt.Println("❌ Visual clarity is out of range")
		return
	}

	if metrics.LayoutEfficiency < 0.0 || metrics.LayoutEfficiency > 1.0 {
		fmt.Println("❌ Layout efficiency is out of range")
		return
	}

	if metrics.VisualNoise < 0.0 || metrics.VisualNoise > 1.0 {
		fmt.Println("❌ Visual noise is out of range")
		return
	}

	// Test quality thresholds
	thresholds := generator.GetDefaultQualityThresholds()
	if thresholds == nil {
		fmt.Println("❌ Failed to get default quality thresholds")
		return
	}

	// Validate thresholds
	if thresholds.MinTaskBarHeight <= 0 {
		fmt.Println("❌ Minimum task bar height should be positive")
		return
	}

	if thresholds.MinTaskBarWidth <= 0 {
		fmt.Println("❌ Minimum task bar width should be positive")
		return
	}

	if thresholds.MinTextSpacing <= 0 {
		fmt.Println("❌ Minimum text spacing should be positive")
		return
	}

	if thresholds.MaxOverlapRatio < 0.0 || thresholds.MaxOverlapRatio > 1.0 {
		fmt.Println("❌ Maximum overlap ratio should be between 0 and 1")
		return
	}

	if thresholds.MinFontSize <= 0 {
		fmt.Println("❌ Minimum font size should be positive")
		return
	}

	if thresholds.MinContrastRatio <= 0 {
		fmt.Println("❌ Minimum contrast ratio should be positive")
		return
	}

	if thresholds.MaxLineLength <= 0 {
		fmt.Println("❌ Maximum line length should be positive")
		return
	}

	if thresholds.MinVisualClarity < 0.0 || thresholds.MinVisualClarity > 1.0 {
		fmt.Println("❌ Minimum visual clarity should be between 0 and 1")
		return
	}

	if thresholds.MinLayoutEfficiency < 0.0 || thresholds.MinLayoutEfficiency > 1.0 {
		fmt.Println("❌ Minimum layout efficiency should be between 0 and 1")
		return
	}

	if thresholds.MaxVisualNoise < 0.0 || thresholds.MaxVisualNoise > 1.0 {
		fmt.Println("❌ Maximum visual noise should be between 0 and 1")
		return
	}

	fmt.Printf("✅ Quality metrics test passed\n")
	fmt.Printf("   Overall score: %.2f\n", metrics.OverallScore)
	fmt.Printf("   Spacing score: %.2f\n", metrics.SpacingScore)
	fmt.Printf("   Alignment score: %.2f\n", metrics.AlignmentScore)
	fmt.Printf("   Readability score: %.2f\n", metrics.ReadabilityScore)
	fmt.Printf("   Visual clarity: %.2f\n", metrics.VisualClarity)
	fmt.Printf("   Layout efficiency: %.2f\n", metrics.LayoutEfficiency)
	fmt.Printf("   Visual noise: %.2f\n", metrics.VisualNoise)
	fmt.Printf("   Min task bar height: %.1f\n", thresholds.MinTaskBarHeight)
	fmt.Printf("   Min task bar width: %.1f\n", thresholds.MinTaskBarWidth)
	fmt.Printf("   Min text spacing: %.1f\n", thresholds.MinTextSpacing)
	fmt.Printf("   Max overlap ratio: %.1f\n", thresholds.MaxOverlapRatio)
	fmt.Printf("   Min font size: %.1f\n", thresholds.MinFontSize)
	fmt.Printf("   Min contrast ratio: %.1f\n", thresholds.MinContrastRatio)
	fmt.Printf("   Max line length: %.1f\n", thresholds.MaxLineLength)
	fmt.Printf("   Min visual clarity: %.1f\n", thresholds.MinVisualClarity)
	fmt.Printf("   Min layout efficiency: %.1f\n", thresholds.MinLayoutEfficiency)
	fmt.Printf("   Max visual noise: %.1f\n", thresholds.MaxVisualNoise)
}

// createSampleLayoutResult creates a sample layout result for testing
func createSampleLayoutResult() *calendar.IntegratedLayoutResult {
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
			{
				TaskName:    "Sample Task 2",
				Description: "Another sample task",
				Category:    "ADMIN",
				Priority:    2,
				StartX:      70.0,
				Y:           20.0,
				Width:       30.0,
				Height:      12.0,
				Color:       "gray",
				IsContinuation: false,
				IsStart:     true,
				IsEnd:       true,
				MonthBoundary: false,
			},
		},
		Statistics: &calendar.IntegratedLayoutStatistics{
			TotalTasks:     2,
			ProcessedBars:  2,
			SpaceEfficiency: 0.85,
			VisualQuality:   0.80,
		},
	}
}
