package main

import (
	"fmt"
	"io/ioutil"
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
	Lengths Lengths
	Colors  Colors
}

type Numbers struct {
	ArrayStretch        float64
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
		if bts, err = ioutil.ReadFile(strings.ToLower(filepath)); err != nil {
			return cfg, fmt.Errorf("ioutil read file: %w", err)
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
	reader := NewReader(cfg.CSVFilePath)
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
