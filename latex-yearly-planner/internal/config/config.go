package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/kudrykv/latex-yearly-planner/internal/data"
	"github.com/kudrykv/latex-yearly-planner/internal/layout"
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
	MonthsWithTasks []data.MonthYear

	Pages Pages

	Layout Layout
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

type Layout struct {
	Paper Paper

	Numbers Numbers
	Lengths layout.Lengths
	Colors  Colors
}

type Numbers struct {
	ArrayStretch float64
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

func NewConfig(pathConfigs ...string) (Config, error) {
	var (
		bts []byte
		err error
		cfg Config
	)

	for _, filepath := range pathConfigs {
		if bts, err = os.ReadFile(strings.ToLower(filepath)); err != nil {
			return cfg, fmt.Errorf("read file: %w", err)
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
	reader := data.NewReader(cfg.CSVFilePath)
	dateRange, err := reader.GetDateRange()
	if err != nil {
		return fmt.Errorf("failed to get date range: %w", err)
	}

	cfg.StartYear = dateRange.Earliest.Year()
	cfg.EndYear = dateRange.Latest.Year()

	// Get months with tasks
	monthsWithTasks, err := reader.GetMonthsWithTasks()
	if err != nil {
		return fmt.Errorf("failed to get months with tasks: %w", err)
	}
	cfg.MonthsWithTasks = monthsWithTasks

	// If we have months with tasks, limit the year range to only those years
	if len(monthsWithTasks) > 0 {
		// Find the unique years from the months with tasks
		yearSet := make(map[int]bool)
		for _, monthYear := range monthsWithTasks {
			yearSet[monthYear.Year] = true
		}

		// Set the year range to only include years with tasks
		years := make([]int, 0, len(yearSet))
		for year := range yearSet {
			years = append(years, year)
		}

		if len(years) > 0 {
			cfg.StartYear = years[0]
			cfg.EndYear = years[len(years)-1]
		}
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
