package app

import (
	"os"

	"phd-dissertation-planner/src/core"

	"github.com/urfave/cli/v2"
)

const (
	fConfig       = "config"
	pConfig       = "preview"
	fOutDir       = "outdir"
	fTestCoverage = "test-coverage"
)

func New() *cli.App {
	// Initialize the composer map
	core.ComposerMap["monthly"] = Monthly

	return &cli.App{
		Name: "plannergen",

		Writer:    os.Stdout,
		ErrWriter: os.Stderr,

		Flags: []cli.Flag{
			&cli.PathFlag{Name: fConfig, Required: false, Value: "src/core/base.yaml", Usage: "config file(s), comma-separated"},
			&cli.BoolFlag{Name: pConfig, Required: false, Usage: "render only one page per unique module"},
			&cli.PathFlag{Name: fOutDir, Required: false, Value: "", Usage: "output directory for generated files (overrides config)"},
			&cli.BoolFlag{Name: "test-coverage", Required: false, Usage: "run tests with coverage analysis"},
			&cli.BoolFlag{Name: "validate", Required: false, Usage: "validate CSV file without generating PDF"},
			&cli.StringFlag{Name: "preset", Required: false, Usage: "Configuration preset: academic, compact, presentation", EnvVars: []string{"PLANNER_PRESET"}},
		},

		Action: action,
	}
}
