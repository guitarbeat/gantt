package data

import (
	"testing"
	"time"
)

func TestIntegration_ValidateDataIntegrity(t *testing.T) {
	validator := NewDataIntegrityValidator()
	
	tasks := []*Task{
		{
			Name:      "Integration Test Task",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Status:    "Planned",
			Priority:  1,
		},
	}
	
	result := validator.ValidateDataIntegrity(tasks)
	if !result.IsValid {
		t.Errorf("Expected valid result, got invalid with %d errors", result.ErrorCount)
	}
}