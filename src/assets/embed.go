package assets

import "embed"

// FS contains embedded small runtime assets for the planner.
// This includes icons, fonts, and other small files needed at runtime.
// Large files like PDFs should remain in static_assets/ to avoid bloating the binary.
//
//go:embed .gitkeep
var FS embed.FS
