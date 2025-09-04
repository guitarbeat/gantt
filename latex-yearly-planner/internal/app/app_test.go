package app

import "testing"

func TestRootFilename(t *testing.T) {
	if got := RootFilename("configs/planner_config.yaml"); got != "planner_config.tex" {
		t.Fatalf("RootFilename wrong: %s", got)
	}
	if got := RootFilename("planner_config.yml"); got != "planner_config.tex" {
		t.Fatalf("RootFilename wrong: %s", got)
	}
	if got := RootFilename("/path/to/other.yaml"); got != "other.tex" {
		t.Fatalf("RootFilename wrong: %s", got)
	}
}
