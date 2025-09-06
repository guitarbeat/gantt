package generator

import (
	"fmt"
	"math"
)

// VisualSpacingConfig defines comprehensive spacing and alignment settings
type VisualSpacingConfig struct {
	// Calendar grid spacing
	CalendarGridSpacing    SpacingConfig `json:"calendar_grid_spacing"`
	TaskBarSpacing         SpacingConfig `json:"task_bar_spacing"`
	TextSpacing            SpacingConfig `json:"text_spacing"`
	
	// Alignment settings
	Alignment              AlignmentConfig `json:"alignment"`
	
	// Visual hierarchy spacing
	HierarchySpacing       HierarchyConfig `json:"hierarchy_spacing"`
	
	// Responsive spacing (adjusts based on content density)
	ResponsiveSpacing      ResponsiveConfig `json:"responsive_spacing"`
	
	// Quality thresholds
	QualityThresholds      QualityThresholds `json:"quality_thresholds"`
}

// SpacingConfig defines spacing parameters for different elements
type SpacingConfig struct {
	// Basic spacing
	Padding        float64 `json:"padding"`        // Internal padding
	Margin         float64 `json:"margin"`         // External margin
	Gap            float64 `json:"gap"`            // Gap between elements
	
	// Advanced spacing
	MinSpacing     float64 `json:"min_spacing"`     // Minimum spacing
	MaxSpacing     float64 `json:"max_spacing"`     // Maximum spacing
	PreferredSpacing float64 `json:"preferred_spacing"` // Preferred spacing
	
	// Spacing units
	Unit           string  `json:"unit"`           // Unit (pt, mm, cm, etc.)
}

// AlignmentConfig defines alignment settings
type AlignmentConfig struct {
	// Horizontal alignment
	Horizontal     string  `json:"horizontal"`     // left, center, right, justify
	
	// Vertical alignment
	Vertical       string  `json:"vertical"`       // top, middle, bottom, baseline
	
	// Text alignment
	TextAlign      string  `json:"text_align"`     // left, center, right, justify
	
	// Task bar alignment
	TaskBarAlign   string  `json:"task_bar_align"` // left, center, right, stretch
}

// HierarchyConfig defines spacing for visual hierarchy
type HierarchyConfig struct {
	// Level-based spacing
	Level1Spacing  float64 `json:"level1_spacing"`  // Critical tasks
	Level2Spacing  float64 `json:"level2_spacing"`  // High priority tasks
	Level3Spacing  float64 `json:"level3_spacing"`  // Medium priority tasks
	Level4Spacing  float64 `json:"level4_spacing"`  // Low priority tasks
	Level5Spacing  float64 `json:"level5_spacing"`  // Minimal priority tasks
	
	// Category-based spacing
	CategorySpacing map[string]float64 `json:"category_spacing"`
	
	// Element-based spacing
	TitleSpacing   float64 `json:"title_spacing"`   // Task title spacing
	DescSpacing    float64 `json:"desc_spacing"`    // Task description spacing
	OverflowSpacing float64 `json:"overflow_spacing"` // Overflow indicator spacing
}

// ResponsiveConfig defines responsive spacing adjustments
type ResponsiveConfig struct {
	// Density-based adjustments
	LowDensity     SpacingMultiplier `json:"low_density"`     // When few tasks
	NormalDensity  SpacingMultiplier `json:"normal_density"`  // Normal task count
	HighDensity    SpacingMultiplier `json:"high_density"`    // Many tasks
	
	// Content-based adjustments
	ShortContent   SpacingMultiplier `json:"short_content"`   // Short task names
	LongContent    SpacingMultiplier `json:"long_content"`    // Long task names
	
	// View-based adjustments
	MonthlyView    SpacingMultiplier `json:"monthly_view"`    // Monthly calendar
	WeeklyView     SpacingMultiplier `json:"weekly_view"`     // Weekly calendar
	DailyView      SpacingMultiplier `json:"daily_view"`      // Daily calendar
}

// SpacingMultiplier defines how spacing should be multiplied
type SpacingMultiplier struct {
	Padding    float64 `json:"padding"`    // Padding multiplier
	Margin     float64 `json:"margin"`     // Margin multiplier
	Gap        float64 `json:"gap"`        // Gap multiplier
	FontSize   float64 `json:"font_size"`  // Font size multiplier
}

// QualityThresholds defines quality thresholds for visual validation
type QualityThresholds struct {
	// Spacing thresholds
	MinTaskBarHeight    float64 `json:"min_task_bar_height"`    // Minimum task bar height
	MinTaskBarWidth     float64 `json:"min_task_bar_width"`     // Minimum task bar width
	MinTextSpacing      float64 `json:"min_text_spacing"`       // Minimum text spacing
	MaxOverlapRatio     float64 `json:"max_overlap_ratio"`      // Maximum allowed overlap
	
	// Readability thresholds
	MinFontSize         float64 `json:"min_font_size"`          // Minimum readable font size
	MinContrastRatio    float64 `json:"min_contrast_ratio"`     // Minimum contrast ratio
	MaxLineLength       float64 `json:"max_line_length"`        // Maximum line length for readability
	
	// Visual quality thresholds
	MinVisualClarity    float64 `json:"min_visual_clarity"`     // Minimum visual clarity score
	MinLayoutEfficiency float64 `json:"min_layout_efficiency"`  // Minimum layout efficiency
	MaxVisualNoise      float64 `json:"max_visual_noise"`       // Maximum visual noise level
}

// GetDefaultVisualSpacingConfig returns the default visual spacing configuration
func GetDefaultVisualSpacingConfig() *VisualSpacingConfig {
	return &VisualSpacingConfig{
		CalendarGridSpacing: SpacingConfig{
			Padding:           2.0,
			Margin:           1.0,
			Gap:              1.5,
			MinSpacing:       0.5,
			MaxSpacing:       4.0,
			PreferredSpacing: 2.0,
			Unit:             "pt",
		},
		TaskBarSpacing: SpacingConfig{
			Padding:           1.5,
			Margin:           0.5,
			Gap:              0.8,
			MinSpacing:       0.3,
			MaxSpacing:       2.5,
			PreferredSpacing: 1.2,
			Unit:             "pt",
		},
		TextSpacing: SpacingConfig{
			Padding:           0.5,
			Margin:           0.3,
			Gap:              0.4,
			MinSpacing:       0.2,
			MaxSpacing:       1.0,
			PreferredSpacing: 0.6,
			Unit:             "pt",
		},
		Alignment: AlignmentConfig{
			Horizontal:   "center",
			Vertical:     "middle",
			TextAlign:    "left",
			TaskBarAlign: "stretch",
		},
		HierarchySpacing: HierarchyConfig{
			Level1Spacing: 2.5, // Critical
			Level2Spacing: 2.0, // High
			Level3Spacing: 1.5, // Medium
			Level4Spacing: 1.0, // Low
			Level5Spacing: 0.8, // Minimal
			CategorySpacing: map[string]float64{
				"PROPOSAL":     2.2,
				"LASER":        2.0,
				"IMAGING":      1.8,
				"ADMIN":        1.5,
				"DISSERTATION": 2.5,
				"RESEARCH":     2.0,
				"PUBLICATION":  2.2,
			},
			TitleSpacing:   1.2,
			DescSpacing:    0.8,
			OverflowSpacing: 0.6,
		},
		ResponsiveSpacing: ResponsiveConfig{
			LowDensity: SpacingMultiplier{
				Padding:  1.2,
				Margin:   1.1,
				Gap:      1.3,
				FontSize: 1.1,
			},
			NormalDensity: SpacingMultiplier{
				Padding:  1.0,
				Margin:   1.0,
				Gap:      1.0,
				FontSize: 1.0,
			},
			HighDensity: SpacingMultiplier{
				Padding:  0.8,
				Margin:   0.9,
				Gap:      0.7,
				FontSize: 0.9,
			},
			ShortContent: SpacingMultiplier{
				Padding:  0.9,
				Margin:   0.9,
				Gap:      0.8,
				FontSize: 1.0,
			},
			LongContent: SpacingMultiplier{
				Padding:  1.1,
				Margin:   1.1,
				Gap:      1.2,
				FontSize: 0.95,
			},
			MonthlyView: SpacingMultiplier{
				Padding:  1.0,
				Margin:   1.0,
				Gap:      1.0,
				FontSize: 1.0,
			},
			WeeklyView: SpacingMultiplier{
				Padding:  1.2,
				Margin:   1.1,
				Gap:      1.3,
				FontSize: 1.1,
			},
			DailyView: SpacingMultiplier{
				Padding:  1.5,
				Margin:   1.3,
				Gap:      1.6,
				FontSize: 1.2,
			},
		},
		QualityThresholds: QualityThresholds{
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
		},
	}
}

// CalculateOptimalSpacing calculates optimal spacing based on content and context
func (vsc *VisualSpacingConfig) CalculateOptimalSpacing(content AnalysisContext) *OptimizedSpacing {
	spacing := &OptimizedSpacing{
		BaseConfig: vsc,
		Context:    content,
	}
	
	// Calculate density-based adjustments
	densityMultiplier := vsc.getDensityMultiplier(content.TaskDensity)
	contentMultiplier := vsc.getContentMultiplier(content.AvgTaskNameLength)
	viewMultiplier := vsc.getViewMultiplier(content.ViewType)
	
	// Apply multipliers
	spacing.CalendarGridSpacing = vsc.applyMultipliers(vsc.CalendarGridSpacing, densityMultiplier, contentMultiplier, viewMultiplier)
	spacing.TaskBarSpacing = vsc.applyMultipliers(vsc.TaskBarSpacing, densityMultiplier, contentMultiplier, viewMultiplier)
	spacing.TextSpacing = vsc.applyMultipliers(vsc.TextSpacing, densityMultiplier, contentMultiplier, viewMultiplier)
	
	// Calculate hierarchy-based spacing
	spacing.HierarchySpacing = vsc.calculateHierarchySpacing(content)
	
	// Calculate alignment adjustments
	spacing.Alignment = vsc.calculateAlignment(content)
	
	return spacing
}

// AnalysisContext provides context for spacing calculations
type AnalysisContext struct {
	TaskDensity        DensityLevel `json:"task_density"`
	AvgTaskNameLength  float64      `json:"avg_task_name_length"`
	ViewType           ViewType     `json:"view_type"`
	TaskCount          int          `json:"task_count"`
	CategoryDistribution map[string]int `json:"category_distribution"`
	PriorityDistribution map[string]int `json:"priority_distribution"`
	AvailableSpace     float64      `json:"available_space"`
	ContentComplexity  float64      `json:"content_complexity"`
}

// DensityLevel represents task density levels
type DensityLevel int

const (
	DensityLow DensityLevel = iota
	DensityNormal
	DensityHigh
	DensityVeryHigh
)

// OptimizedSpacing contains calculated optimal spacing values
type OptimizedSpacing struct {
	BaseConfig        *VisualSpacingConfig `json:"base_config"`
	Context           AnalysisContext      `json:"context"`
	CalendarGridSpacing SpacingConfig      `json:"calendar_grid_spacing"`
	TaskBarSpacing    SpacingConfig        `json:"task_bar_spacing"`
	TextSpacing       SpacingConfig        `json:"text_spacing"`
	HierarchySpacing  HierarchyConfig      `json:"hierarchy_spacing"`
	Alignment         AlignmentConfig      `json:"alignment"`
	QualityScore      float64              `json:"quality_score"`
}

// getDensityMultiplier returns spacing multiplier based on task density
func (vsc *VisualSpacingConfig) getDensityMultiplier(density DensityLevel) SpacingMultiplier {
	switch density {
	case DensityLow:
		return vsc.ResponsiveSpacing.LowDensity
	case DensityNormal:
		return vsc.ResponsiveSpacing.NormalDensity
	case DensityHigh:
		return vsc.ResponsiveSpacing.HighDensity
	case DensityVeryHigh:
		return SpacingMultiplier{
			Padding:  0.7,
			Margin:   0.8,
			Gap:      0.6,
			FontSize: 0.85,
		}
	default:
		return vsc.ResponsiveSpacing.NormalDensity
	}
}

// getContentMultiplier returns spacing multiplier based on content length
func (vsc *VisualSpacingConfig) getContentMultiplier(avgLength float64) SpacingMultiplier {
	if avgLength < 15 {
		return vsc.ResponsiveSpacing.ShortContent
	} else if avgLength > 30 {
		return vsc.ResponsiveSpacing.LongContent
	}
	return vsc.ResponsiveSpacing.NormalDensity
}

// getViewMultiplier returns spacing multiplier based on view type
func (vsc *VisualSpacingConfig) getViewMultiplier(viewType ViewType) SpacingMultiplier {
	switch viewType {
	case ViewTypeMonthly:
		return vsc.ResponsiveSpacing.MonthlyView
	case ViewTypeWeekly:
		return vsc.ResponsiveSpacing.WeeklyView
	case ViewTypeDaily:
		return vsc.ResponsiveSpacing.DailyView
	default:
		return vsc.ResponsiveSpacing.MonthlyView
	}
}

// applyMultipliers applies multiple spacing multipliers to a spacing config
func (vsc *VisualSpacingConfig) applyMultipliers(base SpacingConfig, multipliers ...SpacingMultiplier) SpacingConfig {
	result := base
	
	for _, mult := range multipliers {
		result.Padding *= mult.Padding
		result.Margin *= mult.Margin
		result.Gap *= mult.Gap
	}
	
	// Ensure values stay within bounds
	result.Padding = math.Max(result.MinSpacing, math.Min(result.MaxSpacing, result.Padding))
	result.Margin = math.Max(result.MinSpacing, math.Min(result.MaxSpacing, result.Margin))
	result.Gap = math.Max(result.MinSpacing, math.Min(result.MaxSpacing, result.Gap))
	
	return result
}

// calculateHierarchySpacing calculates spacing based on task hierarchy
func (vsc *VisualSpacingConfig) calculateHierarchySpacing(context AnalysisContext) HierarchyConfig {
	hierarchy := vsc.HierarchySpacing
	
	// Adjust based on category distribution
	for category, count := range context.CategoryDistribution {
		if spacing, exists := hierarchy.CategorySpacing[category]; exists {
			// Adjust spacing based on frequency
			frequency := float64(count) / float64(context.TaskCount)
			if frequency > 0.3 { // High frequency
				hierarchy.CategorySpacing[category] = spacing * 0.9
			} else if frequency < 0.1 { // Low frequency
				hierarchy.CategorySpacing[category] = spacing * 1.1
			}
		}
	}
	
	return hierarchy
}

// calculateAlignment calculates optimal alignment based on context
func (vsc *VisualSpacingConfig) calculateAlignment(context AnalysisContext) AlignmentConfig {
	alignment := vsc.Alignment
	
	// Adjust alignment based on content complexity
	if context.ContentComplexity > 0.7 {
		// High complexity: prefer left alignment for readability
		alignment.TextAlign = "left"
		alignment.TaskBarAlign = "left"
	} else if context.ContentComplexity < 0.3 {
		// Low complexity: prefer center alignment for aesthetics
		alignment.TextAlign = "center"
		alignment.TaskBarAlign = "center"
	}
	
	// Adjust based on available space
	if context.AvailableSpace < 100 {
		// Limited space: prefer compact alignment
		alignment.Horizontal = "left"
		alignment.Vertical = "top"
	}
	
	return alignment
}

// ValidateSpacing validates spacing against quality thresholds
func (os *OptimizedSpacing) ValidateSpacing() *SpacingValidationResult {
	result := &SpacingValidationResult{
		IsValid: true,
		Issues:  make([]string, 0),
		Score:   1.0,
	}
	
	thresholds := os.BaseConfig.QualityThresholds
	
	// Validate task bar spacing
	if os.TaskBarSpacing.Padding < thresholds.MinTextSpacing {
		result.Issues = append(result.Issues, "Task bar padding too small for readability")
		result.IsValid = false
		result.Score *= 0.8
	}
	
	// Validate text spacing
	if os.TextSpacing.Gap < thresholds.MinTextSpacing {
		result.Issues = append(result.Issues, "Text spacing too small for readability")
		result.IsValid = false
		result.Score *= 0.9
	}
	
	// Validate hierarchy spacing
	for level, spacing := range map[string]float64{
		"Level1": os.HierarchySpacing.Level1Spacing,
		"Level2": os.HierarchySpacing.Level2Spacing,
		"Level3": os.HierarchySpacing.Level3Spacing,
		"Level4": os.HierarchySpacing.Level4Spacing,
		"Level5": os.HierarchySpacing.Level5Spacing,
	} {
		if spacing < thresholds.MinTextSpacing {
			result.Issues = append(result.Issues, fmt.Sprintf("%s spacing too small", level))
			result.IsValid = false
			result.Score *= 0.95
		}
	}
	
	// Calculate overall quality score
	os.QualityScore = result.Score
	
	return result
}

// SpacingValidationResult contains validation results for spacing
type SpacingValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Issues  []string `json:"issues"`
	Score   float64  `json:"score"`
}

// GenerateLaTeXSpacingCommands generates LaTeX commands for the optimized spacing
func (os *OptimizedSpacing) GenerateLaTeXSpacingCommands() string {
	commands := ""
	
	// Calendar grid spacing
	commands += fmt.Sprintf("\\setlength{\\myLenTabColSep}{%.1f%s}\n", 
		os.CalendarGridSpacing.Gap, os.CalendarGridSpacing.Unit)
	commands += fmt.Sprintf("\\setlength{\\myLenMonthlyCellHeight}{%.1f%s}\n", 
		os.CalendarGridSpacing.Padding*2, os.CalendarGridSpacing.Unit)
	
	// Task bar spacing
	commands += fmt.Sprintf("\\setlength{\\TaskPaddingH}{%.1f%s}\n", 
		os.TaskBarSpacing.Padding, os.TaskBarSpacing.Unit)
	commands += fmt.Sprintf("\\setlength{\\TaskPaddingV}{%.1f%s}\n", 
		os.TaskBarSpacing.Margin, os.TaskBarSpacing.Unit)
	
	// Text spacing
	commands += fmt.Sprintf("\\setlength{\\myLenLineHeightButLine}{%.1f%s}\n", 
		os.TextSpacing.Gap, os.TextSpacing.Unit)
	
	// Hierarchy spacing
	commands += fmt.Sprintf("\\setlength{\\TaskBarCornerRadius}{%.1f%s}\n", 
		os.HierarchySpacing.Level3Spacing*0.3, os.TaskBarSpacing.Unit)
	commands += fmt.Sprintf("\\setlength{\\TaskBorderWidth}{%.1f%s}\n", 
		os.HierarchySpacing.Level3Spacing*0.2, os.TaskBarSpacing.Unit)
	
	return commands
}
