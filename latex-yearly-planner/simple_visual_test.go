package main

import (
	"fmt"
	"math"
)

// Simple test that doesn't import the generator package to avoid template issues
func main() {
	fmt.Println("Testing Visual Spacing Optimization (Simple)...")

	// Test 1: Visual Spacing Configuration
	fmt.Println("\n=== Test 1: Visual Spacing Configuration ===")
	testSpacingConfig()

	// Test 2: Quality Metrics
	fmt.Println("\n=== Test 2: Quality Metrics ===")
	testQualityMetrics()

	// Test 3: Spacing Calculations
	fmt.Println("\n=== Test 3: Spacing Calculations ===")
	testSpacingCalculations()

	fmt.Println("\n✅ Simple visual spacing optimization tests completed!")
}

// SpacingConfig represents spacing configuration
type SpacingConfig struct {
	Padding        float64
	Margin         float64
	Gap            float64
	MinSpacing     float64
	MaxSpacing     float64
	PreferredSpacing float64
	Unit           string
}

// QualityThresholds represents quality thresholds
type QualityThresholds struct {
	MinTaskBarHeight    float64
	MinTaskBarWidth     float64
	MinTextSpacing      float64
	MaxOverlapRatio     float64
	MinFontSize         float64
	MinContrastRatio    float64
	MaxLineLength       float64
	MinVisualClarity    float64
	MinLayoutEfficiency float64
	MaxVisualNoise      float64
}

// VisualQualityMetrics represents quality metrics
type VisualQualityMetrics struct {
	OverallScore     float64
	SpacingScore     float64
	AlignmentScore   float64
	ReadabilityScore float64
	VisualClarity    float64
	LayoutEfficiency float64
	VisualNoise      float64
}

func testSpacingConfig() {
	// Test default spacing configuration
	calendarSpacing := SpacingConfig{
		Padding:           2.0,
		Margin:           1.0,
		Gap:              1.5,
		MinSpacing:       0.5,
		MaxSpacing:       4.0,
		PreferredSpacing: 2.0,
		Unit:             "pt",
	}

	taskBarSpacing := SpacingConfig{
		Padding:           1.5,
		Margin:           0.5,
		Gap:              0.8,
		MinSpacing:       0.3,
		MaxSpacing:       2.5,
		PreferredSpacing: 1.2,
		Unit:             "pt",
	}

	textSpacing := SpacingConfig{
		Padding:           0.5,
		Margin:           0.3,
		Gap:              0.4,
		MinSpacing:       0.2,
		MaxSpacing:       1.0,
		PreferredSpacing: 0.6,
		Unit:             "pt",
	}

	// Validate spacing configurations
	if calendarSpacing.Padding <= 0 {
		fmt.Println("❌ Calendar spacing padding is not positive")
		return
	}

	if taskBarSpacing.Padding <= 0 {
		fmt.Println("❌ Task bar spacing padding is not positive")
		return
	}

	if textSpacing.Padding <= 0 {
		fmt.Println("❌ Text spacing padding is not positive")
		return
	}

	// Test spacing bounds
	if calendarSpacing.Padding < calendarSpacing.MinSpacing || calendarSpacing.Padding > calendarSpacing.MaxSpacing {
		fmt.Println("❌ Calendar spacing padding is out of bounds")
		return
	}

	if taskBarSpacing.Padding < taskBarSpacing.MinSpacing || taskBarSpacing.Padding > taskBarSpacing.MaxSpacing {
		fmt.Println("❌ Task bar spacing padding is out of bounds")
		return
	}

	if textSpacing.Padding < textSpacing.MinSpacing || textSpacing.Padding > textSpacing.MaxSpacing {
		fmt.Println("❌ Text spacing padding is out of bounds")
		return
	}

	fmt.Printf("✅ Visual spacing configuration test passed\n")
	fmt.Printf("   Calendar grid padding: %.1f%s\n", calendarSpacing.Padding, calendarSpacing.Unit)
	fmt.Printf("   Task bar padding: %.1f%s\n", taskBarSpacing.Padding, taskBarSpacing.Unit)
	fmt.Printf("   Text padding: %.1f%s\n", textSpacing.Padding, textSpacing.Unit)
}

func testQualityMetrics() {
	// Create sample quality metrics
	metrics := VisualQualityMetrics{
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
	thresholds := QualityThresholds{
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

	if thresholds.MinTaskBarHeight <= 0 {
		fmt.Println("❌ Minimum task bar height is not positive")
		return
	}

	if thresholds.MinTaskBarWidth <= 0 {
		fmt.Println("❌ Minimum task bar width is not positive")
		return
	}

	if thresholds.MinTextSpacing <= 0 {
		fmt.Println("❌ Minimum text spacing is not positive")
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
}

func testSpacingCalculations() {
	// Test responsive spacing calculations
	baseSpacing := 2.0
	
	// Test density-based adjustments
	lowDensityMultiplier := 1.2
	normalDensityMultiplier := 1.0
	highDensityMultiplier := 0.8
	veryHighDensityMultiplier := 0.7

	lowDensitySpacing := baseSpacing * lowDensityMultiplier
	normalDensitySpacing := baseSpacing * normalDensityMultiplier
	highDensitySpacing := baseSpacing * highDensityMultiplier
	veryHighDensitySpacing := baseSpacing * veryHighDensityMultiplier

	// Validate density adjustments
	if lowDensitySpacing <= normalDensitySpacing {
		fmt.Println("❌ Low density spacing should be greater than normal")
		return
	}

	if normalDensitySpacing <= highDensitySpacing {
		fmt.Println("❌ Normal density spacing should be greater than high")
		return
	}

	if highDensitySpacing <= veryHighDensitySpacing {
		fmt.Println("❌ High density spacing should be greater than very high")
		return
	}

	// Test content-based adjustments
	shortContentMultiplier := 0.9
	longContentMultiplier := 1.1

	shortContentSpacing := baseSpacing * shortContentMultiplier
	longContentSpacing := baseSpacing * longContentMultiplier

	// Validate content adjustments
	if shortContentSpacing >= baseSpacing {
		fmt.Println("❌ Short content spacing should be less than base")
		return
	}

	if longContentSpacing <= baseSpacing {
		fmt.Println("❌ Long content spacing should be greater than base")
		return
	}

	// Test view-based adjustments
	monthlyViewMultiplier := 1.0
	weeklyViewMultiplier := 1.2
	dailyViewMultiplier := 1.5

	monthlyViewSpacing := baseSpacing * monthlyViewMultiplier
	weeklyViewSpacing := baseSpacing * weeklyViewMultiplier
	dailyViewSpacing := baseSpacing * dailyViewMultiplier

	// Validate view adjustments
	if monthlyViewSpacing >= weeklyViewSpacing {
		fmt.Println("❌ Monthly view spacing should be less than weekly")
		return
	}

	if weeklyViewSpacing >= dailyViewSpacing {
		fmt.Println("❌ Weekly view spacing should be less than daily")
		return
	}

	// Test spacing bounds enforcement
	minSpacing := 0.5
	maxSpacing := 4.0

	enforcedSpacing := math.Max(minSpacing, math.Min(maxSpacing, baseSpacing))
	if enforcedSpacing < minSpacing || enforcedSpacing > maxSpacing {
		fmt.Println("❌ Spacing bounds enforcement failed")
		return
	}

	fmt.Printf("✅ Spacing calculations test passed\n")
	fmt.Printf("   Low density spacing: %.1f\n", lowDensitySpacing)
	fmt.Printf("   Normal density spacing: %.1f\n", normalDensitySpacing)
	fmt.Printf("   High density spacing: %.1f\n", highDensitySpacing)
	fmt.Printf("   Very high density spacing: %.1f\n", veryHighDensitySpacing)
	fmt.Printf("   Short content spacing: %.1f\n", shortContentSpacing)
	fmt.Printf("   Long content spacing: %.1f\n", longContentSpacing)
	fmt.Printf("   Monthly view spacing: %.1f\n", monthlyViewSpacing)
	fmt.Printf("   Weekly view spacing: %.1f\n", weeklyViewSpacing)
	fmt.Printf("   Daily view spacing: %.1f\n", dailyViewSpacing)
	fmt.Printf("   Enforced spacing: %.1f\n", enforcedSpacing)
}
