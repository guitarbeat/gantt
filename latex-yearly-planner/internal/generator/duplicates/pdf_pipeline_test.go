package generator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
)

func TestPDFPipeline(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdf-pipeline-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create PDF pipeline
	pipeline := NewPDFPipeline(tempDir, outputDir)

	// Test configuration validation
	t.Run("ConfigValidation", func(t *testing.T) {
		// Test empty config
		emptyCfg := config.Config{}
		err := pipeline.validateConfig(emptyCfg)
		if err == nil {
			t.Error("Expected error for empty configuration")
		}

		// Test valid config
		validCfg := config.Config{
			MonthsWithTasks: []data.MonthYear{
				{Month: time.January, Year: 2024},
			},
		}
		err = pipeline.validateConfig(validCfg)
		if err != nil {
			t.Errorf("Expected no error for valid configuration, got: %v", err)
		}
	})

	t.Run("TempWorkspace", func(t *testing.T) {
		workspace, err := pipeline.createTempWorkspace()
		if err != nil {
			t.Fatalf("Failed to create temp workspace: %v", err)
		}
		defer os.RemoveAll(workspace)

		// Check if directory exists
		if _, err := os.Stat(workspace); os.IsNotExist(err) {
			t.Error("Temp workspace was not created")
		}

		// Check if directory is writable
		testFile := filepath.Join(workspace, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Errorf("Temp workspace is not writable: %v", err)
		}
	})

	t.Run("PreambleGeneration", func(t *testing.T) {
		cfg := config.Config{}
		options := PDFGenerationOptions{
			ExtraPackages:  []string{"graphicx", "hyperref"},
			CustomPreamble: "% Custom preamble test",
		}

		// Create a buffer to capture output
		var buf []byte
		writer := &bytesBuffer{data: &buf}

		err := pipeline.writePreamble(writer, cfg, options)
		if err != nil {
			t.Fatalf("Failed to write preamble: %v", err)
		}

		preamble := string(buf)
		
		// Check for required packages
		if !simpleContains(preamble, "\\usepackage{graphicx}") {
			t.Error("Extra package 'graphicx' not found in preamble")
		}
		
		if !simpleContains(preamble, "\\usepackage{hyperref}") {
			t.Error("Extra package 'hyperref' not found in preamble")
		}

		// Check for custom preamble
		if !simpleContains(preamble, "% Custom preamble test") {
			t.Error("Custom preamble not found")
		}

		// Check for required document structure
		if !simpleContains(preamble, "\\begin{document}") {
			t.Error("Document begin not found in preamble")
		}
	})

	t.Run("FileOperations", func(t *testing.T) {
		// Test file copying
		srcFile := filepath.Join(tempDir, "source.txt")
		dstFile := filepath.Join(tempDir, "destination.txt")
		
		// Create source file
		testContent := "Test file content for copying"
		if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Test copy
		if err := pipeline.copyFile(srcFile, dstFile); err != nil {
			t.Fatalf("Failed to copy file: %v", err)
		}

		// Verify destination file
		content, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("Failed to read destination file: %v", err)
		}

		if string(content) != testContent {
			t.Errorf("File content mismatch. Expected: %s, Got: %s", testContent, string(content))
		}
	})
}

func TestPDFGenerationOptions(t *testing.T) {
	t.Run("DefaultOptions", func(t *testing.T) {
		options := PDFGenerationOptions{}
		
		// Test defaults
		if options.CompilationEngine != "" {
			t.Error("Expected empty default compilation engine")
		}
		
		if options.MaxRetries != 0 {
			t.Error("Expected zero default max retries")
		}
		
		if options.CleanupTempFiles != false {
			t.Error("Expected false default for cleanup temp files")
		}
	})

	t.Run("CustomOptions", func(t *testing.T) {
		options := PDFGenerationOptions{
			OutputFileName:    "custom.pdf",
			CleanupTempFiles:  true,
			MaxRetries:        5,
			CompilationEngine: "xelatex",
			ExtraPackages:     []string{"tikz", "pgfplots"},
			CustomPreamble:    "% Custom settings",
		}

		if options.OutputFileName != "custom.pdf" {
			t.Error("Custom output filename not set correctly")
		}
		
		if !options.CleanupTempFiles {
			t.Error("Custom cleanup setting not set correctly")
		}
		
		if options.MaxRetries != 5 {
			t.Error("Custom max retries not set correctly")
		}
		
		if options.CompilationEngine != "xelatex" {
			t.Error("Custom compilation engine not set correctly")
		}
		
		if len(options.ExtraPackages) != 2 {
			t.Error("Extra packages not set correctly")
		}
		
		if options.CustomPreamble != "% Custom settings" {
			t.Error("Custom preamble not set correctly")
		}
	})
}

func TestPDFGenerationResult(t *testing.T) {
	result := &PDFGenerationResult{
		Success:           true,
		OutputPath:        "/path/to/output.pdf",
		CompilationTime:   time.Second * 30,
		LaTeXLog:          "Compilation log content",
		Errors:            []string{"Error 1", "Error 2"},
		Warnings:          []string{"Warning 1"},
		PageCount:         12,
		FileSize:          1024000,
		TempFilesCreated:  []string{"/tmp/file1", "/tmp/file2"},
		LayoutStatistics:  "Some statistics",
	}

	// Test result structure
	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.OutputPath != "/path/to/output.pdf" {
		t.Error("Output path not set correctly")
	}

	if result.CompilationTime != time.Second*30 {
		t.Error("Compilation time not set correctly")
	}

	if len(result.Errors) != 2 {
		t.Error("Errors not set correctly")
	}

	if len(result.Warnings) != 1 {
		t.Error("Warnings not set correctly")
	}

	if result.PageCount != 12 {
		t.Error("Page count not set correctly")
	}

	if result.FileSize != 1024000 {
		t.Error("File size not set correctly")
	}

	if len(result.TempFilesCreated) != 2 {
		t.Error("Temp files not set correctly")
	}
}

// Helper types and functions for testing

type bytesBuffer struct {
	data *[]byte
}

func (b *bytesBuffer) Write(p []byte) (int, error) {
	*b.data = append(*b.data, p...)
	return len(p), nil
}

// Simplified string contains check
func simpleContains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	
	return false
}

func TestSimpleContains(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "", true},
		{"a", "", true},
		{"", "a", false},
		{"\\usepackage{tikz}", "\\usepackage{tikz}", true},
		{"\\usepackage{tikz}\\usepackage{pgf}", "\\usepackage{pgf}", true},
	}

	for _, tt := range tests {
		if got := simpleContains(tt.s, tt.substr); got != tt.want {
			t.Errorf("simpleContains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}
