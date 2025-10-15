package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"phd-dissertation-planner/src/core"
)

// BenchmarkTemplateRendering benchmarks template rendering performance
func BenchmarkTemplateRendering(b *testing.B) {
	cfg := core.DefaultConfig()
	cfg.Year = 2024

	tpl := NewTpl()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := tpl.Document(&buf, cfg)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes() // Prevent optimization
	}
}

// BenchmarkConfigurationLoading benchmarks config loading performance
func BenchmarkConfigurationLoading(b *testing.B) {
	// Create temporary config file
	tempDir := b.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	configContent := `
year: 2024
output_dir: generated
csv_file_path: input_data/research_timeline_v5.1_comprehensive.csv

pages:
  - name: monthly
    render_blocks:
      - calendar
      - toc

layout:
  calendar_layout:
    day_number_width: 0.8
    day_content_margin: 0.1
    task_cell_margin: 0.05
    task_cell_spacing: 0.02
    day_cell_minipage_width: 0.95
    header_angle_size_offset: 0.0

typography:
  font_size: 10pt
  line_spread: 1.1

document:
  paper_size: a4paper
  orientation: portrait
  margin: 1.0in
  header_margin: 0.5in
  footer_margin: 0.5in

latex:
  hyphen_penalty: 10000
  tolerance: 200
  tab_col_sep: 0.8em
  header_side_months_width: 2.5cm
  monthly_cell_height: 1.2cm
  header_resize_box: 2.5cm
  line_thickness_default: 0.4pt
  line_thickness_thick: 0.8pt
  col_sep: 1.0em
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.NewConfig(configFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMemoryUsage benchmarks memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	cfg := core.DefaultConfig()
	cfg.Year = 2024

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tpl := NewTpl()
		var buf bytes.Buffer
		err := tpl.Document(&buf, cfg)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes() // Prevent optimization
	}
}

// BenchmarkConcurrentProcessing benchmarks concurrent template processing
func BenchmarkConcurrentProcessing(b *testing.B) {
	cfg := core.DefaultConfig()
	cfg.Year = 2024

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		tpl := NewTpl()
		for pb.Next() {
			var buf bytes.Buffer
			err := tpl.Document(&buf, cfg)
			if err != nil {
				b.Fatal(err)
			}
			_ = buf.Bytes() // Prevent optimization
		}
	})
}
