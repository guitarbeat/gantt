package templates

import "embed"

// FS contains the embedded LaTeX template files for the planner.
//go:embed monthly/*.tpl
var FS embed.FS
