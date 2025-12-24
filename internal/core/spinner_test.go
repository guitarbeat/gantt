package core

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to synchronize output reading
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Test Spinner
	s := NewSpinner("Test spinner", false)
	s.Start()
	time.Sleep(150 * time.Millisecond) // Let it spin a bit
	s.UpdateMessage("Updated message")
	time.Sleep(150 * time.Millisecond)
	s.Stop(true)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	output := <-outC

	// Verify output contains expected strings
	// Note: We can't easily check for the animation characters since they are overwritten with \r
	// But we can check for the final success message
	if !strings.Contains(output, "Updated message") {
		t.Errorf("Expected output to contain 'Updated message', got: %q", output)
	}
	if !strings.Contains(output, "✔") {
		t.Errorf("Expected output to contain checkmark, got: %q", output)
	}
}

func TestSpinnerSilent(t *testing.T) {
	// Redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	s := NewSpinner("Silent spinner", true)
	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop(true)

	w.Close()
	os.Stdout = oldStdout
	output := <-outC

	if output != "" {
		t.Errorf("Expected no output for silent spinner, got: %q", output)
	}
}
