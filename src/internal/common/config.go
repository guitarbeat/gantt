package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug Debug

	Year                int `env:"PLANNER_YEAR"`
	WeekStart           time.Weekday
	Dotted              bool
	CalAfterSchedule    bool
	ClearTopRightCorner bool
	AMPMTime            bool
	AddLastHalfHour     bool

	// Data source configuration
	CSVFilePath string `env:"PLANNER_CSV_FILE"`
	StartYear   int    `env:"PLANNER_START_YEAR"`
	EndYear     int    `env:"PLANNER_END_YEAR"`

	// Months with tasks (populated from CSV)
	MonthsWithTasks []MonthYear

	Pages Pages

	Layout Layout

	// OutputDir is the directory where generated .tex and .pdf files will be written
	// Defaults to "build" when not provided via environment or config
	OutputDir string `env:"PLANNER_OUTPUT_DIR"`

}


type Debug struct {
	ShowFrame bool
	ShowLinks bool
}

type Pages []Page
type Page struct {
	Name         string
	RenderBlocks RenderBlocks
}

type RenderBlocks []RenderBlock

type Modules []Module
type Module struct {
	Cfg  Config
	Tpl  string
	Body interface{}
}

type RenderBlock struct {
	FuncName string
	Tpls     []string
}

type Colors struct {
	Gray      string
	LightGray string
}

type TaskColors struct {
	Proposal     RGBColor
	Laser        RGBColor
	Imaging      RGBColor
	Admin        RGBColor
	Dissertation RGBColor
	Research     RGBColor
	Publication  RGBColor
}

type RGBColor struct {
	R int
	G int
	B int
}

type LaTeX struct {
	TabColSep             string
	HeaderSideMonthsWidth string
	TaskBorderWidth       string
	TaskPaddingH          string
	TaskPaddingV          string
	TaskVerticalOffset    string
	ArrayStretch          float64
	TaskFontSize          string
	TaskBarHeight         string
	MonthlyCellHeight     string
	HeaderResizeBox       string
	LineThicknessDefault   string
	LineThicknessThick     string

	// Task styling parameters
	TaskBackgroundOpacity int    `yaml:"task_background_opacity"`
	TaskBorderOpacity     int    `yaml:"task_border_opacity"`
	TaskContentSpacing    string `yaml:"task_content_spacing"`

	// TColorBox styling
	TColorBox TColorBox `yaml:"tcolorbox"`

	// Typography settings
	Typography Typography `yaml:"typography"`

	// Spacing and layout
	Spacing Spacing `yaml:"spacing"`

	// Document settings
	Document Document `yaml:"document"`
}

type TColorBox struct {
	Arc       string
	Left      string
	Right     string
	Top       string
	Bottom    string
	BoxRule   string
	TaskBoxRule string
	TaskArc   string
	TaskLeft  string
	TaskRight string
	TaskTop   string
	TaskBottom string
}

type Typography struct {
	HyphenPenalty       int
	Tolerance           int
	EmergencyStretch    string
	SloppyEmergencyStretch string
}

type Spacing struct {
	TwoColSep      string
	TriColSep      string
	FiveColSep     string
	TableColSep    string
	ColorLegendSep string
	PageBreak      string
}

type Document struct {
	FontSize  string
	ParIndent string
	FBoxSep   string
}

type Constraints struct {
	MaxStackHeight     float64
	MinTaskHeight      float64
	MaxTaskHeight      float64
	MinTaskWidth       float64
	MaxTaskWidth       float64
	VerticalSpacing    float64
	HorizontalSpacing  float64
	CollisionThreshold float64
	OverflowThreshold  float64
	ExpansionThreshold float64

	// Task sizing multipliers
	MinTaskHeightMultiplier float64 `yaml:"min_task_height_multiplier"`
	MaxTaskHeightMultiplier float64 `yaml:"max_task_height_multiplier"`
	MinTaskWidthMultiplier  float64 `yaml:"min_task_width_multiplier"`
	MaxTaskWidthDays        float64 `yaml:"max_task_width_days"`

	// Visual styling
	TaskBarOpacity   float64 `yaml:"task_bar_opacity"`
	BorderWidth      float64 `yaml:"border_width"`
	OpacityThreshold float64 `yaml:"opacity_threshold"`

	// Algorithm thresholds
	CompressionThreshold     float64 `yaml:"compression_threshold"`
	QualityThreshold         float64 `yaml:"quality_threshold"`
	BalanceThreshold         float64 `yaml:"balance_threshold"`
	EfficiencyThreshold      float64 `yaml:"efficiency_threshold"`
	EfficiencyGoodThreshold  float64 `yaml:"efficiency_good_threshold"`

	// Task positioning
	TaskVerticalOffset           float64 `yaml:"task_vertical_offset"`
	TaskHorizontalOffset         float64 `yaml:"task_horizontal_offset"`
	ExpandedTaskVerticalOffset   float64 `yaml:"expanded_task_vertical_offset"`
	ExpandedTaskHorizontalOffset float64 `yaml:"expanded_task_horizontal_offset"`
	ExpandedTaskHeightMultiplier float64 `yaml:"expanded_task_height_multiplier"`
	ExpandedTaskWidthMultiplier  float64 `yaml:"expanded_task_width_multiplier"`
	CollapsedTaskHeightMultiplier float64 `yaml:"collapsed_task_height_multiplier"`
	CollapsedTaskWidthMultiplier float64 `yaml:"collapsed_task_width_multiplier"`
}

type Layout struct {
	Paper Paper

	Numbers     Numbers
	Lengths     Lengths
	Colors      Colors
	TaskColors  TaskColors `yaml:"task_colors"`
	LaTeX       LaTeX      `yaml:"latex"`
	Constraints Constraints
	Calendar    Calendar
	Stacking    Stacking
}

type Calendar struct {
	DayNumberWidth    string
	DayContentMargin  string
	TaskKernSpacing   string
	CollapseThreshold int
	TaskBarOpacity    float64
	BorderWidth       float64

	// Typography settings for calendar content
	EmergencyStretch string `yaml:"emergencystretch"`
	InnerSep         string `yaml:"inner_sep"`

	// Cell rendering parameters
	DayCellMinipageWidth string `yaml:"day_cell_minipage_width"`
	TaskCellMargin       string `yaml:"task_cell_margin"`
	TaskCellSpacing      string `yaml:"task_cell_spacing"`

	// Task rendering parameters
	MaxTaskChars            int    `yaml:"max_task_chars"`
	MaxTaskCharsCompact     int    `yaml:"max_task_chars_compact"`
	MaxTaskCharsVeryCompact int    `yaml:"max_task_chars_very_compact"`
	DefaultTaskSpacing      string `yaml:"default_task_spacing"`
	FirstTaskSpacing        string `yaml:"first_task_spacing"`
	DefaultTaskHeight       string `yaml:"default_task_height"`
	FirstTaskHeight         string `yaml:"first_task_height"`
}

type Stacking struct {
	BaseHeight float64 `yaml:"base_height"`
	MinHeight  float64 `yaml:"min_height"`
	MaxHeight  float64 `yaml:"max_height"`

	// Prominence multipliers
	ProminenceCritical float64 `yaml:"prominence_critical"`
	ProminenceHigh     float64 `yaml:"prominence_high"`
	ProminenceMedium   float64 `yaml:"prominence_medium"`
	ProminenceLow      float64 `yaml:"prominence_low"`
	ProminenceMinimal  float64 `yaml:"prominence_minimal"`

	// Duration-based weighting
	DurationShortWeight  float64 `yaml:"duration_short_weight"`
	DurationMediumWeight float64 `yaml:"duration_medium_weight"`
	DurationLongWeight   float64 `yaml:"duration_long_weight"`

	// Complexity-based weighting
	ComplexityMinimalWeight float64 `yaml:"complexity_minimal_weight"`
	ComplexityNormalWeight  float64 `yaml:"complexity_normal_weight"`
	ComplexityComplexWeight float64 `yaml:"complexity_complex_weight"`

	// Task positioning parameters
	VerticalSpacing   float64 `yaml:"vertical_spacing"`
	HorizontalSpacing float64 `yaml:"horizontal_spacing"`

	// Visual quality thresholds
	VisibilityThreshold float64 `yaml:"visibility_threshold"`
	OverflowVertical    float64 `yaml:"overflow_vertical"`
	CollisionThreshold  float64 `yaml:"collision_threshold"`
	BoundingBoxBuffer   float64 `yaml:"bounding_box_buffer"`

	// Quality assessment thresholds
	QualityHighThreshold      float64 `yaml:"quality_high_threshold"`
	BalanceHighThreshold      float64 `yaml:"balance_high_threshold"`
	CompressionHighThreshold  float64 `yaml:"compression_high_threshold"`

	// Weight calculation parameters
	DurationWeightFactor float64 `yaml:"duration_weight_factor"`
	CategoryHighWeight   float64 `yaml:"category_high_weight"`
	CategoryMediumWeight float64 `yaml:"category_medium_weight"`
	CategoryNormalWeight float64 `yaml:"category_normal_weight"`
	CategoryHighBonus    float64 `yaml:"category_high_bonus"`
}

type Numbers struct {
	ArrayStretch float64
}

type Lengths struct {
	TabColSep             string
	LineThicknessDefault  string
	LineThicknessThick    string
	LineHeightButLine     string
	TwoColSep             string
	TriColSep             string
	FiveColSep            string
	MonthlyCellHeight     string
	HeaderResizeBox       string
	HeaderSideMonthsWidth string
	MonthlySpring         string
}

type Paper struct {
	Width  string `env:"PLANNER_LAYOUT_PAPER_WIDTH"`
	Height string `env:"PLANNER_LAYOUT_PAPER_HEIGHT"`

	Margin Margin

	ReverseMargins bool
	MarginParWidth string
	MarginParSep   string
}

type Margin struct {
	Top    string `env:"PLANNER_LAYOUT_PAPER_MARGIN_TOP"`
	Bottom string `env:"PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM"`
	Left   string `env:"PLANNER_LAYOUT_PAPER_MARGIN_LEFT"`
	Right  string `env:"PLANNER_LAYOUT_PAPER_MARGIN_RIGHT"`
}

// Composer is a function type for generating modules
type Composer func(cfg Config, tpls []string) (Modules, error)

// ComposerMap maps function names to their implementations
// This will be populated by the app package to avoid circular imports
var ComposerMap = map[string]Composer{}

// FilterUniqueModules removes duplicate modules based on template name
func FilterUniqueModules(array []Module) []Module {
	filtered := make([]Module, 0)
	found := map[string]bool{}

	for _, val := range array {
		if _, present := found[val.Tpl]; !present {
			filtered = append(filtered, val)
			found[val.Tpl] = true
		}
	}

	return filtered
}

// NewConfig creates a new configuration from config files and environment variables
func NewConfig(pathConfigs ...string) (Config, error) {
	var (
		bts []byte
		err error
		cfg Config
	)

	for _, filepath := range pathConfigs {
		// * Skip missing files instead of failing
		if bts, err = os.ReadFile(strings.ToLower(filepath)); err != nil {
			if os.IsNotExist(err) {
				// * File doesn't exist, skip it
				continue
			}
			return cfg, fmt.Errorf("read file: %w", err)
		}

		// * Skip empty files
		if len(strings.TrimSpace(string(bts))) == 0 {
			continue
		}

		if err = yaml.Unmarshal(bts, &cfg); err != nil {
			return cfg, fmt.Errorf("yaml unmarshal: %w", err)
		}
	}

	if err = env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("env parse: %w", err)
	}

	if cfg.Year == 0 {
		cfg.Year = time.Now().Year()
	}

	// Default output dir
	if strings.TrimSpace(cfg.OutputDir) == "" {
		cfg.OutputDir = "build"
	}

	// If CSV file is provided, determine date range dynamically
	if cfg.CSVFilePath != "" {
		if err := cfg.setDateRangeFromCSV(); err != nil {
			return cfg, fmt.Errorf("failed to set date range from CSV: %w", err)
		}
	}

	return cfg, nil
}

// setDateRangeFromCSV reads the CSV file and sets the start and end years
func (cfg *Config) setDateRangeFromCSV() error {
	reader := NewReader(cfg.CSVFilePath)

	// Read tasks once and derive both date range and months from the same slice
	tasks, err := reader.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks found in CSV file")
	}

	// Compute earliest/latest dates and collect months with tasks
	earliest := tasks[0].StartDate
	latest := tasks[0].EndDate
	monthsSet := make(map[MonthYear]bool)

	for _, task := range tasks {
		if task.StartDate.Before(earliest) {
			earliest = task.StartDate
		}
		if task.EndDate.After(latest) {
			latest = task.EndDate
		}

		// Walk months from task start to end and add to set
		current := time.Date(task.StartDate.Year(), task.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(task.EndDate.Year(), task.EndDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		for !current.After(end) {
			monthsSet[MonthYear{Year: current.Year(), Month: current.Month()}] = true
			current = current.AddDate(0, 1, 0)
		}
	}

	// Assign year range
	cfg.StartYear = earliest.Year()
	cfg.EndYear = latest.Year()

	// Convert months set to a sorted slice for deterministic ordering
	months := make([]MonthYear, 0, len(monthsSet))
	for m := range monthsSet {
		months = append(months, m)
	}

	// Sort by year, then by month
	for i := 0; i < len(months)-1; i++ {
		for j := 0; j < len(months)-i-1; j++ {
			if months[j].Year > months[j+1].Year ||
				(months[j].Year == months[j+1].Year && months[j].Month > months[j+1].Month) {
				months[j], months[j+1] = months[j+1], months[j]
			}
		}
	}
	cfg.MonthsWithTasks = months

	// Limit year range to only years present in months (keeps behavior consistent)
	if len(months) > 0 {
		yearSet := make(map[int]bool)
		for _, my := range months {
			yearSet[my.Year] = true
		}
		years := make([]int, 0, len(yearSet))
		for y := range yearSet {
			years = append(years, y)
		}
		// Sort years
		for i := 0; i < len(years)-1; i++ {
			for j := 0; j < len(years)-i-1; j++ {
				if years[j] > years[j+1] {
					years[j], years[j+1] = years[j+1], years[j]
				}
			}
		}
		cfg.StartYear = years[0]
		cfg.EndYear = years[len(years)-1]
	}

	// Update the main Year field to the start year if not explicitly set
	if cfg.Year == time.Now().Year() {
		cfg.Year = cfg.StartYear
	}

	return nil
}

// GetYears returns a slice of years to generate planners for
func (cfg *Config) GetYears() []int {
	if cfg.StartYear == 0 || cfg.EndYear == 0 {
		// Fallback to single year
		return []int{cfg.Year}
	}

	years := make([]int, 0, cfg.EndYear-cfg.StartYear+1)
	for year := cfg.StartYear; year <= cfg.EndYear; year++ {
		years = append(years, year)
	}

	return years
}
