package generator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	cal "latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/config"
	tmplfs "latex-yearly-planner/templates"
)

var tpl = func() *template.Template {
	t := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}

				dict[key] = values[i+1]
			}

			return dict, nil
		},

		"incr": func(i int) int {
			return i + 1
		},

		"dec": func(i int) int {
			return i - 1
		},

		"is": func(i interface{}) bool {
			if value, ok := i.(bool); ok {
				return value
			}

			return i != nil
		},

		// Layout integration functions
		"hasLayoutData": func(data interface{}) bool {
			if data == nil {
				return false
			}
			// Check if data has layout-related fields
			if m, ok := data.(map[string]interface{}); ok {
				_, hasLayout := m["LayoutResult"]
				_, hasTaskBars := m["TaskBars"]
				return hasLayout || hasTaskBars
			}
			return false
		},

		"getTaskBars": func(data interface{}) []*cal.IntegratedTaskBar {
			if m, ok := data.(map[string]interface{}); ok {
				if bars, ok := m["TaskBars"].([]*cal.IntegratedTaskBar); ok {
					return bars
				}
			}
			return nil
		},

		"getLayoutStats": func(data interface{}) *cal.IntegratedLayoutStatistics {
			if m, ok := data.(map[string]interface{}); ok {
				if stats, ok := m["LayoutStats"].(*cal.IntegratedLayoutStatistics); ok {
					return stats
				}
			}
			return nil
		},

		"formatTaskBar": func(bar *cal.IntegratedTaskBar) string {
			if bar == nil {
				return ""
			}
			// Convert priority to prominence level
			var prominence string
			switch {
			case bar.Priority >= 4:
				prominence = "CRITICAL"
			case bar.Priority >= 3:
				prominence = "HIGH"
			case bar.Priority >= 2:
				prominence = "MEDIUM"
			case bar.Priority >= 1:
				prominence = "LOW"
			default:
				prominence = "MINIMAL"
			}
			
			// Generate LaTeX for individual task bar using the visual design system
			return fmt.Sprintf("\\TaskOverlayBoxP{%s}{%s}{%s}{%s}",
				prominence,     // prominence level
				bar.Color,      // category color
				bar.TaskName,   // task name
				bar.Description, // description
			)
		},
	})

	// Choose source of templates: embedded by default, filesystem when DEV_TEMPLATES is set
	var (
		err   error
		useFS fs.FS
	)

	if os.Getenv("DEV_TEMPLATES") != "" {
		// Use on-disk templates for development override
		useFS = os.DirFS(filepath.Join("templates", "monthly"))
	} else {
		// Use embedded templates from templates.FS
		// Narrow to monthly/ subdir
		var sub fs.FS
		sub, err = fs.Sub(tmplfs.FS, "monthly")
		if err != nil {
			panic(fmt.Sprintf("failed to sub FS for monthly templates: %v", err))
		}
		useFS = sub
	}

	// Parse all *.tpl templates from the selected FS
	t, err = t.ParseFS(useFS, "*.tpl")
	if err != nil {
		panic(fmt.Sprintf("failed to parse monthly templates: %v", err))
	}

	return t
}()

type Tpl struct {
	tpl *template.Template
}

func NewTpl() Tpl {
	return Tpl{
		tpl: tpl,
	}
}

func (t Tpl) Document(wr io.Writer, cfg config.Config) error {
	type pack struct {
		Cfg   config.Config
		Pages []config.Page
	}

	data := pack{Cfg: cfg, Pages: cfg.Pages}
	if err := t.tpl.ExecuteTemplate(wr, "main_document.tpl", data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func (t Tpl) Execute(wr io.Writer, name string, data interface{}) error {
	if err := t.tpl.ExecuteTemplate(wr, name, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
