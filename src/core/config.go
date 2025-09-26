package core

import (
	"fmt"
	"os"
	"strconv"
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

// TaskColors removed - using algorithmic colors

// RGBColor removed - using algorithmic colors

type AlgorithmicColors struct {
	Proposal     string
	Laser        string
	Imaging      string
	Admin        string
	Dissertation string
	Research     string
	Publication  string
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
	LineThicknessDefault  string
	LineThicknessThick    string
	ColSep                string

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
	Arc         string
	Left        string
	Right       string
	Top         string
	Bottom      string
	BoxRule     string
	TaskBoxRule string
	TaskArc     string
	TaskLeft    string
	TaskRight   string
	TaskTop     string
	TaskBottom  string
}

type Typography struct {
	HyphenPenalty          int
	Tolerance              int
	EmergencyStretch       string
	SloppyEmergencyStretch string
}

type Spacing struct {
	ColSep          string
	TableColSep     string
	ColorLegendSep  string
	// Page break hardcoded in templates
	// Template-specific spacing values
	Col              string `yaml:"col"`
	TaskContentVspace string `yaml:"task_content_vspace"`
	TaskOverlayArc   string `yaml:"task_overlay_arc"`
}

type Document struct {
	FontSize  string
	ParIndent string
}

type Constraints struct {
	MaxStackHeight     float64
	MinTaskHeight      float64
	MaxTaskHeight      float64
	MinTaskWidth       float64
	MaxTaskWidth       float64
	// Spacing values are hardcoded in stacking.go
	CollisionThreshold float64
	OverflowThreshold  float64
	ExpansionThreshold float64

	// Task sizing constraints
	MaxTaskWidthDays float64 `yaml:"max_task_width_days"`

	// Visual styling (inherited from layout_engine)

	// Algorithm thresholds (inherited from layout_engine)

	// Task positioning (hardcoded in layout_manager.go)
}

type Layout struct {
	Paper Paper

	Numbers     Numbers
	Lengths     Lengths
	Colors            Colors
	// TaskColors removed - using algorithmic colors
	AlgorithmicColors AlgorithmicColors
	LaTeX       LaTeX      `yaml:"latex"`
	Constraints Constraints
	Calendar    Calendar
	Stacking    Stacking
	LayoutEngine LayoutEngine `yaml:"layout_engine"`
}

type Calendar struct {
	TaskKernSpacing string `yaml:"taskkernspacing"`
	// Other parameters hardcoded in calendar.go
}

type Stacking struct {
	BaseHeight float64 `yaml:"base_height"`
	MinHeight  float64 `yaml:"min_height"`
	MaxHeight  float64 `yaml:"max_height"`

	// Visual quality thresholds
	OverflowVertical   float64 `yaml:"overflow_vertical"`
	CollisionThreshold float64 `yaml:"collision_threshold"`
	// Other thresholds hardcoded in stacking.go
}

type LayoutEngine struct {
	// Task positioning multipliers (relative to day dimensions)
	InitialYPositionMultiplier float64 `yaml:"initial_y_position_multiplier"`
	TaskHeightMultiplier       float64 `yaml:"task_height_multiplier"`
	MaxTaskWidthDays           float64 `yaml:"max_task_width_days"`

	// Visual weight calculation multipliers
	DurationLongMultiplier    float64 `yaml:"duration_long_multiplier"`
	DurationShortMultiplier   float64 `yaml:"duration_short_multiplier"`
	MilestoneWeightMultiplier float64 `yaml:"milestone_weight_multiplier"`
	CategoryWeightMultiplier  float64 `yaml:"category_weight_multiplier"`

	// Urgency multipliers removed - using simplified prominence calculation


	// Quality assessment thresholds (hardcoded as constants)

	// Task rendering configuration
	TaskRendering TaskRendering `yaml:"task_rendering"`

	// Typography settings for task text (inherited from main typography)

	// Grid constraints
	GridConstraints LayoutGridConstraints `yaml:"grid_constraints"`

	// Visual styling
	// Visual styling (hardcoded in layout_manager.go)

	// Calendar layout constants
	CalendarLayout LayoutCalendarLayout `yaml:"calendar_layout"`

	// Task density calculation removed - not used in code
}

// UrgencyMultipliers struct removed - using simplified prominence calculation

type TaskRendering struct {
	DefaultSpacing   string `yaml:"default_spacing"`
	FirstTaskSpacing string `yaml:"first_task_spacing"`
	DefaultHeight    string `yaml:"default_height"`
	FirstTaskHeight  string `yaml:"first_task_height"`
	VerticalSpacing  string `yaml:"vertical_spacing"`
}

// LayoutTypography struct removed - not used in code

type LayoutGridConstraints struct {
	MinTaskSpacing     float64 `yaml:"min_task_spacing"`
	MaxTaskSpacing     float64 `yaml:"max_task_spacing"`
	MinRowHeight       float64 `yaml:"min_row_height"`
	MaxRowHeight       float64 `yaml:"max_row_height"`
	MinColumnWidth     float64 `yaml:"min_column_width"`
	MaxColumnWidth     float64 `yaml:"max_column_width"`
	GridResolution     float64 `yaml:"grid_resolution"`
	AlignmentTolerance float64 `yaml:"alignment_tolerance"`
	CollisionBuffer    float64 `yaml:"collision_buffer"`
	TransitionBuffer   float64 `yaml:"transition_buffer"`
}

// LayoutVisualStyling removed - values are hardcoded in layout_manager.go

type LayoutCalendarLayout struct {
	DayNumberWidth        string `yaml:"day_number_width"`
	DayContentMargin      string `yaml:"day_content_margin"`
	TaskCellMargin        string `yaml:"task_cell_margin"`
	TaskCellSpacing       string `yaml:"task_cell_spacing"`
	DayCellMinipageWidth  string `yaml:"day_cell_minipage_width"`
	HeaderAngleSizeOffset string `yaml:"header_angle_size_offset"`
}

// LayoutDensityCalculation struct removed - not used in code

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

	// * Set defaults for layout engine configuration
	cfg.setLayoutEngineDefaults()

	// * Set algorithmic colors for predefined categories
	cfg.setAlgorithmicColors()

	// * Validate layout engine configuration
	if err := cfg.validateLayoutEngineConfig(); err != nil {
		return cfg, fmt.Errorf("layout engine config validation failed: %w", err)
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

// setAlgorithmicColors sets the algorithmic colors for predefined categories
func (cfg *Config) setAlgorithmicColors() {
	cfg.Layout.AlgorithmicColors = AlgorithmicColors{
		Proposal:     hexToRGB(generateCategoryColor("PROPOSAL")),
		Laser:        hexToRGB(generateCategoryColor("LASER")),
		Imaging:      hexToRGB(generateCategoryColor("IMAGING")),
		Admin:        hexToRGB(generateCategoryColor("ADMIN")),
		Dissertation: hexToRGB(generateCategoryColor("DISSERTATION")),
		Research:     hexToRGB(generateCategoryColor("RESEARCH")),
		Publication:  hexToRGB(generateCategoryColor("PUBLICATION")),
	}
}

// hexToRGB converts hex color to RGB format for LaTeX
func hexToRGB(hex string) string {
	// Remove # prefix if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Convert hex to RGB
	if len(hex) == 6 {
		// Parse hex values
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		return fmt.Sprintf("%d,%d,%d", r, g, b)
	}

	// Fallback for invalid hex
	return "128,128,128"
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

// setLayoutEngineDefaults sets default values for layout engine configuration
func (cfg *Config) setLayoutEngineDefaults() {
	// * Set defaults for layout engine if not already set
	if cfg.Layout.LayoutEngine.InitialYPositionMultiplier == 0 {
		cfg.Layout.LayoutEngine.InitialYPositionMultiplier = 0.1
	}
	if cfg.Layout.LayoutEngine.TaskHeightMultiplier == 0 {
		cfg.Layout.LayoutEngine.TaskHeightMultiplier = 0.6
	}
	if cfg.Layout.LayoutEngine.MaxTaskWidthDays == 0 {
		cfg.Layout.LayoutEngine.MaxTaskWidthDays = 7.0
	}
	if cfg.Layout.LayoutEngine.DurationLongMultiplier == 0 {
		cfg.Layout.LayoutEngine.DurationLongMultiplier = 1.2
	}
	if cfg.Layout.LayoutEngine.DurationShortMultiplier == 0 {
		cfg.Layout.LayoutEngine.DurationShortMultiplier = 0.8
	}
	if cfg.Layout.LayoutEngine.MilestoneWeightMultiplier == 0 {
		cfg.Layout.LayoutEngine.MilestoneWeightMultiplier = 1.5
	}
	if cfg.Layout.LayoutEngine.CategoryWeightMultiplier == 0 {
		cfg.Layout.LayoutEngine.CategoryWeightMultiplier = 1.0
	}

	// Urgency multiplier defaults removed - using simplified prominence calculation

	// Quality threshold defaults removed - using hardcoded constants

	// * Set task rendering defaults
	if cfg.Layout.LayoutEngine.TaskRendering.DefaultSpacing == "" {
		cfg.Layout.LayoutEngine.TaskRendering.DefaultSpacing = "0.8ex"
	}
	if cfg.Layout.LayoutEngine.TaskRendering.FirstTaskSpacing == "" {
		cfg.Layout.LayoutEngine.TaskRendering.FirstTaskSpacing = "0.5ex"
	}
	if cfg.Layout.LayoutEngine.TaskRendering.DefaultHeight == "" {
		cfg.Layout.LayoutEngine.TaskRendering.DefaultHeight = "3.0ex"
	}
	if cfg.Layout.LayoutEngine.TaskRendering.FirstTaskHeight == "" {
		cfg.Layout.LayoutEngine.TaskRendering.FirstTaskHeight = "3.5ex"
	}
	if cfg.Layout.LayoutEngine.TaskRendering.VerticalSpacing == "" {
		cfg.Layout.LayoutEngine.TaskRendering.VerticalSpacing = "0.1ex"
	}

	// * Set typography defaults (using main typography settings)
	if cfg.Layout.LaTeX.Typography.HyphenPenalty == 0 {
		cfg.Layout.LaTeX.Typography.HyphenPenalty = 50
	}
	if cfg.Layout.LaTeX.Typography.Tolerance == 0 {
		cfg.Layout.LaTeX.Typography.Tolerance = 1000
	}
	if cfg.Layout.LaTeX.Typography.EmergencyStretch == "" {
		cfg.Layout.LaTeX.Typography.EmergencyStretch = "2em"
	}
	if cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch == "" {
		cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch = "3em"
	}

	// * Set grid constraints defaults
	if cfg.Layout.LayoutEngine.GridConstraints.MinTaskSpacing == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MinTaskSpacing = 1.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing = 10.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.MinRowHeight == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MinRowHeight = 8.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.MaxRowHeight == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MaxRowHeight = 20.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.MinColumnWidth == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MinColumnWidth = 5.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.MaxColumnWidth == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.MaxColumnWidth = 50.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.GridResolution == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.GridResolution = 1.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.AlignmentTolerance == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.AlignmentTolerance = 0.5
	}
	if cfg.Layout.LayoutEngine.GridConstraints.CollisionBuffer == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.CollisionBuffer = 2.0
	}
	if cfg.Layout.LayoutEngine.GridConstraints.TransitionBuffer == 0 {
		cfg.Layout.LayoutEngine.GridConstraints.TransitionBuffer = 2.0
	}

	// Visual styling defaults removed - values are hardcoded in layout_manager.go

	// * Set calendar layout defaults
	if cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth = "6mm"
	}
	if cfg.Layout.LayoutEngine.CalendarLayout.DayContentMargin == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.DayContentMargin = "8mm"
	}
	if cfg.Layout.LayoutEngine.CalendarLayout.TaskCellMargin == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.TaskCellMargin = "8mm"
	}
	if cfg.Layout.LayoutEngine.CalendarLayout.TaskCellSpacing == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.TaskCellSpacing = "6mm"
	}
	if cfg.Layout.LayoutEngine.CalendarLayout.DayCellMinipageWidth == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.DayCellMinipageWidth = "6mm"
	}
	if cfg.Layout.LayoutEngine.CalendarLayout.HeaderAngleSizeOffset == "" {
		cfg.Layout.LayoutEngine.CalendarLayout.HeaderAngleSizeOffset = "0.86pt"
	}

	// Density calculation defaults removed - not used in code
}

// validateLayoutEngineConfig validates the layout engine configuration
func (cfg *Config) validateLayoutEngineConfig() error {
	// * Validate multiplier ranges (0.0 to 10.0)
	multipliers := []struct {
		name  string
		value float64
	}{
		{"initial_y_position_multiplier", cfg.Layout.LayoutEngine.InitialYPositionMultiplier},
		{"task_height_multiplier", cfg.Layout.LayoutEngine.TaskHeightMultiplier},
		{"max_task_width_days", cfg.Layout.LayoutEngine.MaxTaskWidthDays},
		{"duration_long_multiplier", cfg.Layout.LayoutEngine.DurationLongMultiplier},
		{"duration_short_multiplier", cfg.Layout.LayoutEngine.DurationShortMultiplier},
		{"milestone_weight_multiplier", cfg.Layout.LayoutEngine.MilestoneWeightMultiplier},
		{"category_weight_multiplier", cfg.Layout.LayoutEngine.CategoryWeightMultiplier},
	}

	for _, m := range multipliers {
		if m.value < 0.0 || m.value > 10.0 {
			return fmt.Errorf("invalid %s: %f (must be between 0.0 and 10.0)", m.name, m.value)
		}
	}

	// Quality threshold validation removed - using hardcoded constants

	// Urgency multiplier validation removed - using simplified prominence calculation

	// * Validate grid constraints
	if cfg.Layout.LayoutEngine.GridConstraints.MinTaskSpacing > cfg.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing {
		return fmt.Errorf("min_task_spacing (%f) cannot be greater than max_task_spacing (%f)",
			cfg.Layout.LayoutEngine.GridConstraints.MinTaskSpacing,
			cfg.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing)
	}

	if cfg.Layout.LayoutEngine.GridConstraints.MinRowHeight > cfg.Layout.LayoutEngine.GridConstraints.MaxRowHeight {
		return fmt.Errorf("min_row_height (%f) cannot be greater than max_row_height (%f)",
			cfg.Layout.LayoutEngine.GridConstraints.MinRowHeight,
			cfg.Layout.LayoutEngine.GridConstraints.MaxRowHeight)
	}

	if cfg.Layout.LayoutEngine.GridConstraints.MinColumnWidth > cfg.Layout.LayoutEngine.GridConstraints.MaxColumnWidth {
		return fmt.Errorf("min_column_width (%f) cannot be greater than max_column_width (%f)",
			cfg.Layout.LayoutEngine.GridConstraints.MinColumnWidth,
			cfg.Layout.LayoutEngine.GridConstraints.MaxColumnWidth)
	}

	return nil
}

