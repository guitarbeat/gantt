package generator

import (
	"fmt"
	"math"
	"strings"
	"time"

	"latex-yearly-planner/internal/calendar"
)

// VisualOptimizer provides visual quality optimization for PDF generation
type VisualOptimizer struct {
	spacingConfig *VisualSpacingConfig
	logger        PDFLogger
}

// NewVisualOptimizer creates a new visual optimizer
func NewVisualOptimizer() *VisualOptimizer {
	return &VisualOptimizer{
		spacingConfig: GetDefaultVisualSpacingConfig(),
		logger:        &VisualOptimizerLogger{},
	}
}

// SetLogger sets the logger for the visual optimizer
func (vo *VisualOptimizer) SetLogger(logger PDFLogger) {
	vo.logger = logger
}

// GetSpacingConfig returns the spacing configuration
func (vo *VisualOptimizer) GetSpacingConfig() *VisualSpacingConfig {
	return vo.spacingConfig
}

// OptimizeLayout optimizes the layout for visual quality
func (vo *VisualOptimizer) OptimizeLayout(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) (*OptimizedLayoutResult, error) {
	vo.logger.Info("Starting visual optimization for %s view", viewType)
	
	// Analyze the layout context
	context := vo.analyzeLayoutContext(layoutResult, viewType)
	
	// Calculate optimal spacing
	optimalSpacing := vo.spacingConfig.CalculateOptimalSpacing(context)
	
	// Validate spacing
	validation := optimalSpacing.ValidateSpacing()
	if !validation.IsValid {
		vo.logger.Error("Spacing validation failed: %v", validation.Issues)
		return nil, fmt.Errorf("spacing validation failed: %s", strings.Join(validation.Issues, ", "))
	}
	
	// Optimize task bars
	optimizedBars := vo.optimizeTaskBars(layoutResult.TaskBars, optimalSpacing)
	
	// Calculate visual quality metrics
	qualityMetrics := vo.calculateQualityMetrics(optimizedBars, optimalSpacing)
	
	// Create optimized layout result
	result := &OptimizedLayoutResult{
		OriginalResult:    layoutResult,
		OptimizedBars:     optimizedBars,
		OptimalSpacing:    optimalSpacing,
		QualityMetrics:    qualityMetrics,
		ValidationResult:  validation,
		ViewType:          viewType,
		OptimizationTime:  time.Now(),
	}
	
	vo.logger.Info("Visual optimization completed with quality score: %.2f", qualityMetrics.OverallScore)
	
	return result, nil
}

// analyzeLayoutContext analyzes the layout context for optimization
func (vo *VisualOptimizer) analyzeLayoutContext(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) AnalysisContext {
	context := AnalysisContext{
		ViewType:              viewType,
		TaskCount:             len(layoutResult.TaskBars),
		CategoryDistribution:  make(map[string]int),
		PriorityDistribution:  make(map[string]int),
		ContentComplexity:     0.0,
		AvgTaskNameLength:     0.0,
		AvailableSpace:        100.0, // Default assumption
	}
	
	if context.TaskCount == 0 {
		context.TaskDensity = DensityLow
		return context
	}
	
	// Calculate task density
	taskDensity := float64(context.TaskCount) / 30.0 // Assuming 30-day month
	switch {
	case taskDensity < 0.3:
		context.TaskDensity = DensityLow
	case taskDensity < 0.7:
		context.TaskDensity = DensityNormal
	case taskDensity < 1.2:
		context.TaskDensity = DensityHigh
	default:
		context.TaskDensity = DensityVeryHigh
	}
	
	// Analyze task bars
	totalNameLength := 0.0
	for _, bar := range layoutResult.TaskBars {
		// Category distribution
		if bar.Category != "" {
			context.CategoryDistribution[bar.Category]++
		}
		
		// Priority distribution
		priority := vo.mapPriorityToString(bar.Priority)
		context.PriorityDistribution[priority]++
		
		// Content analysis
		nameLength := float64(len(bar.TaskName))
		totalNameLength += nameLength
		
		// Content complexity (based on description length and special characters)
		descLength := float64(len(bar.Description))
		complexity := (nameLength + descLength) / 50.0 // Normalize
		context.ContentComplexity += complexity
	}
	
	// Calculate averages
	context.AvgTaskNameLength = totalNameLength / float64(context.TaskCount)
	context.ContentComplexity /= float64(context.TaskCount)
	
	// Adjust available space based on view type
	switch viewType {
	case ViewTypeMonthly:
		context.AvailableSpace = 100.0
	case ViewTypeWeekly:
		context.AvailableSpace = 120.0
	case ViewTypeDaily:
		context.AvailableSpace = 150.0
	default:
		context.AvailableSpace = 100.0
	}
	
	return context
}

// mapPriorityToString maps priority number to string
func (vo *VisualOptimizer) mapPriorityToString(priority int) string {
	switch {
	case priority >= 5:
		return "CRITICAL"
	case priority >= 4:
		return "HIGH"
	case priority >= 3:
		return "MEDIUM"
	case priority >= 2:
		return "LOW"
	default:
		return "MINIMAL"
	}
}

// optimizeTaskBars optimizes task bars for visual quality
func (vo *VisualOptimizer) optimizeTaskBars(originalBars []*calendar.IntegratedTaskBar, spacing *OptimizedSpacing) []*OptimizedTaskBar {
	optimizedBars := make([]*OptimizedTaskBar, len(originalBars))
	
	for i, bar := range originalBars {
		optimized := &OptimizedTaskBar{
			OriginalBar: bar,
			OptimizedX:  bar.StartX,
			OptimizedY:  bar.Y,
			OptimizedWidth:  bar.Width,
			OptimizedHeight: bar.Height,
			VisualQuality:   1.0,
		}
		
		// Apply spacing optimizations
		vo.optimizeTaskBarSpacing(optimized, spacing)
		
		// Apply alignment optimizations
		vo.optimizeTaskBarAlignment(optimized, spacing)
		
		// Apply visual enhancements
		vo.enhanceTaskBarVisuals(optimized, spacing)
		
		optimizedBars[i] = optimized
	}
	
	return optimizedBars
}

// optimizeTaskBarSpacing optimizes spacing for a task bar
func (vo *VisualOptimizer) optimizeTaskBarSpacing(bar *OptimizedTaskBar, spacing *OptimizedSpacing) {
	// Get hierarchy spacing
	priority := vo.mapPriorityToString(bar.OriginalBar.Priority)
	_ = vo.getHierarchySpacing(priority, spacing) // Use hierarchy spacing for optimization
	
	// Apply spacing adjustments
	bar.OptimizedHeight = math.Max(bar.OptimizedHeight, spacing.BaseConfig.QualityThresholds.MinTaskBarHeight)
	bar.OptimizedWidth = math.Max(bar.OptimizedWidth, spacing.BaseConfig.QualityThresholds.MinTaskBarWidth)
	
	// Adjust for category spacing
	if categorySpacing, exists := spacing.HierarchySpacing.CategorySpacing[bar.OriginalBar.Category]; exists {
		bar.OptimizedHeight += categorySpacing * 0.1 // Small adjustment
	}
	
	// Apply responsive spacing
	bar.OptimizedHeight *= vo.getResponsiveMultiplier(spacing.Context.TaskDensity)
	bar.OptimizedWidth *= vo.getResponsiveMultiplier(spacing.Context.TaskDensity)
}

// optimizeTaskBarAlignment optimizes alignment for a task bar
func (vo *VisualOptimizer) optimizeTaskBarAlignment(bar *OptimizedTaskBar, spacing *OptimizedSpacing) {
	// Apply horizontal alignment
	switch spacing.Alignment.Horizontal {
	case "left":
		// Keep original X position
	case "center":
		bar.OptimizedX += bar.OptimizedWidth * 0.1 // Slight centering adjustment
	case "right":
		bar.OptimizedX += bar.OptimizedWidth * 0.2
	}
	
	// Apply vertical alignment
	switch spacing.Alignment.Vertical {
	case "top":
		// Keep original Y position
	case "middle":
		bar.OptimizedY += bar.OptimizedHeight * 0.1
	case "bottom":
		bar.OptimizedY += bar.OptimizedHeight * 0.2
	}
}

// enhanceTaskBarVisuals enhances visual appearance of task bars
func (vo *VisualOptimizer) enhanceTaskBarVisuals(bar *OptimizedTaskBar, spacing *OptimizedSpacing) {
	// Calculate visual quality score
	qualityScore := 1.0
	
	// Check minimum size requirements
	if bar.OptimizedHeight >= spacing.BaseConfig.QualityThresholds.MinTaskBarHeight {
		qualityScore *= 1.1
	}
	if bar.OptimizedWidth >= spacing.BaseConfig.QualityThresholds.MinTaskBarWidth {
		qualityScore *= 1.1
	}
	
	// Check spacing requirements
	if bar.OptimizedHeight >= spacing.TaskBarSpacing.MinSpacing {
		qualityScore *= 1.05
	}
	
	// Apply quality-based adjustments
	bar.VisualQuality = math.Min(qualityScore, 1.5) // Cap at 1.5
	
	// Enhance based on priority
	priority := bar.OriginalBar.Priority
	if priority >= 4 {
		bar.OptimizedHeight *= 1.1 // Slightly larger for high priority
		bar.OptimizedWidth *= 1.05
	}
}

// getHierarchySpacing gets hierarchy spacing for a priority level
func (vo *VisualOptimizer) getHierarchySpacing(priority string, spacing *OptimizedSpacing) float64 {
	switch priority {
	case "CRITICAL":
		return spacing.HierarchySpacing.Level1Spacing
	case "HIGH":
		return spacing.HierarchySpacing.Level2Spacing
	case "MEDIUM":
		return spacing.HierarchySpacing.Level3Spacing
	case "LOW":
		return spacing.HierarchySpacing.Level4Spacing
	case "MINIMAL":
		return spacing.HierarchySpacing.Level5Spacing
	default:
		return spacing.HierarchySpacing.Level3Spacing
	}
}

// getResponsiveMultiplier gets responsive multiplier for task density
func (vo *VisualOptimizer) getResponsiveMultiplier(density DensityLevel) float64 {
	switch density {
	case DensityLow:
		return 1.1
	case DensityNormal:
		return 1.0
	case DensityHigh:
		return 0.9
	case DensityVeryHigh:
		return 0.8
	default:
		return 1.0
	}
}

// calculateQualityMetrics calculates visual quality metrics
func (vo *VisualOptimizer) calculateQualityMetrics(bars []*OptimizedTaskBar, spacing *OptimizedSpacing) *VisualQualityMetrics {
	metrics := &VisualQualityMetrics{
		OverallScore:     0.0,
		SpacingScore:     0.0,
		AlignmentScore:   0.0,
		ReadabilityScore: 0.0,
		VisualClarity:    0.0,
		LayoutEfficiency: 0.0,
		VisualNoise:      0.0,
	}
	
	if len(bars) == 0 {
		return metrics
	}
	
	// Calculate spacing score
	spacingScore := 0.0
	for _, bar := range bars {
		if bar.OptimizedHeight >= spacing.BaseConfig.QualityThresholds.MinTaskBarHeight {
			spacingScore += 1.0
		}
		if bar.OptimizedWidth >= spacing.BaseConfig.QualityThresholds.MinTaskBarWidth {
			spacingScore += 1.0
		}
	}
	metrics.SpacingScore = spacingScore / float64(len(bars)*2)
	
	// Calculate alignment score
	alignmentScore := 0.0
	for _, bar := range bars {
		// Check if bar is properly aligned within bounds
		if bar.OptimizedX >= 0 && bar.OptimizedY >= 0 {
			alignmentScore += 1.0
		}
	}
	metrics.AlignmentScore = alignmentScore / float64(len(bars))
	
	// Calculate readability score
	readabilityScore := 0.0
	for _, bar := range bars {
		// Check if text will be readable
		if bar.OptimizedHeight >= 8.0 && bar.OptimizedWidth >= 12.0 {
			readabilityScore += 1.0
		}
	}
	metrics.ReadabilityScore = readabilityScore / float64(len(bars))
	
	// Calculate visual clarity
	metrics.VisualClarity = (metrics.SpacingScore + metrics.AlignmentScore + metrics.ReadabilityScore) / 3.0
	
	// Calculate layout efficiency
	usedSpace := 0.0
	totalSpace := 100.0 * 100.0 // Assuming 100x100 space
	for _, bar := range bars {
		usedSpace += bar.OptimizedWidth * bar.OptimizedHeight
	}
	metrics.LayoutEfficiency = usedSpace / totalSpace
	
	// Calculate visual noise (inverse of clarity)
	metrics.VisualNoise = 1.0 - metrics.VisualClarity
	
	// Calculate overall score
	metrics.OverallScore = (metrics.VisualClarity + metrics.LayoutEfficiency + (1.0 - metrics.VisualNoise)) / 3.0
	
	return metrics
}

// GenerateOptimizedLaTeX generates LaTeX code with optimized spacing
func (vo *VisualOptimizer) GenerateOptimizedLaTeX(optimizedResult *OptimizedLayoutResult) (string, error) {
	var laTeX strings.Builder
	
	// Add visual spacing template
	laTeX.WriteString("\\input{templates/monthly/visual_spacing.tpl}\n")
	
	// Add spacing configuration
	laTeX.WriteString(optimizedResult.OptimalSpacing.GenerateLaTeXSpacingCommands())
	laTeX.WriteString("\n")
	
	// Add view-specific spacing
	laTeX.WriteString(fmt.Sprintf("\\ViewSpecificSpacing{%s}\n", 
		strings.ToLower(string(optimizedResult.ViewType))))
	laTeX.WriteString("\n")
	
	// Generate optimized task bars
	laTeX.WriteString("\\begin{tikzpicture}[overlay, remember picture]\n")
	for _, bar := range optimizedResult.OptimizedBars {
		laTeX.WriteString(vo.generateOptimizedTaskBarLaTeX(bar))
	}
	laTeX.WriteString("\\end{tikzpicture}\n")
	
	return laTeX.String(), nil
}

// generateOptimizedTaskBarLaTeX generates LaTeX for an optimized task bar
func (vo *VisualOptimizer) generateOptimizedTaskBarLaTeX(bar *OptimizedTaskBar) string {
	priority := vo.mapPriorityToString(bar.OriginalBar.Priority)
	
	// Use professional task bar macro
	return fmt.Sprintf(
		"\\node[anchor=north west] at (%.2f,%.2f) {\n"+
		"  \\ProfessionalTaskBar{%s}{%s}{%s}{%s}\n"+
		"};\n",
		bar.OptimizedX,
		bar.OptimizedY,
		priority,
		bar.OriginalBar.Category,
		bar.OriginalBar.TaskName,
		bar.OriginalBar.Description,
	)
}

// OptimizedLayoutResult contains the result of visual optimization
type OptimizedLayoutResult struct {
	OriginalResult    *calendar.IntegratedLayoutResult `json:"original_result"`
	OptimizedBars     []*OptimizedTaskBar              `json:"optimized_bars"`
	OptimalSpacing    *OptimizedSpacing                `json:"optimal_spacing"`
	QualityMetrics    *VisualQualityMetrics            `json:"quality_metrics"`
	ValidationResult  *SpacingValidationResult         `json:"validation_result"`
	ViewType          ViewType                         `json:"view_type"`
	OptimizationTime  time.Time                        `json:"optimization_time"`
}

// OptimizedTaskBar represents an optimized task bar
type OptimizedTaskBar struct {
	OriginalBar    *calendar.IntegratedTaskBar `json:"original_bar"`
	OptimizedX     float64                     `json:"optimized_x"`
	OptimizedY     float64                     `json:"optimized_y"`
	OptimizedWidth float64                     `json:"optimized_width"`
	OptimizedHeight float64                    `json:"optimized_height"`
	VisualQuality  float64                     `json:"visual_quality"`
}

// VisualQualityMetrics contains visual quality measurements
type VisualQualityMetrics struct {
	OverallScore     float64 `json:"overall_score"`
	SpacingScore     float64 `json:"spacing_score"`
	AlignmentScore   float64 `json:"alignment_score"`
	ReadabilityScore float64 `json:"readability_score"`
	VisualClarity    float64 `json:"visual_clarity"`
	LayoutEfficiency float64 `json:"layout_efficiency"`
	VisualNoise      float64 `json:"visual_noise"`
}

// VisualOptimizerLogger provides logging for visual optimizer
type VisualOptimizerLogger struct{}

func (l *VisualOptimizerLogger) Info(msg string, args ...interface{})  { fmt.Printf("[VISUAL-INFO] "+msg+"\n", args...) }
func (l *VisualOptimizerLogger) Error(msg string, args ...interface{}) { fmt.Printf("[VISUAL-ERROR] "+msg+"\n", args...) }
func (l *VisualOptimizerLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[VISUAL-DEBUG] "+msg+"\n", args...) }
