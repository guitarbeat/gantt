package common_test

import (
	"phd-dissertation-planner/internal/common"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config, err := common.NewConfig("testdata/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if config.Year == 0 {
		t.Error("Expected year to be set")
	}
	
	if config.OutputDir == "" {
		t.Error("Expected output directory to be set")
	}
}

func TestNewConfigNotFound(t *testing.T) {
	_, err := common.NewConfig("nonexistent.yaml")
	if err != nil {
		t.Error("Expected no error for non-existent config file (should skip)")
	}
}

func TestConfigGetYears(t *testing.T) {
	config := &common.Config{
		Year:      2024,
		StartYear: 2024,
		EndYear:   2025,
	}
	
	years := config.GetYears()
	if len(years) != 2 {
		t.Errorf("Expected 2 years, got %d", len(years))
	}
	
	if years[0] != 2024 {
		t.Errorf("Expected first year to be 2024, got %d", years[0])
	}
	
	if years[1] != 2025 {
		t.Errorf("Expected second year to be 2025, got %d", years[1])
	}
}

func TestConfigGetYearsFallback(t *testing.T) {
	config := &common.Config{
		Year: 2024,
	}
	
	years := config.GetYears()
	if len(years) != 1 {
		t.Errorf("Expected 1 year, got %d", len(years))
	}
	
	if years[0] != 2024 {
		t.Errorf("Expected year to be 2024, got %d", years[0])
	}
}