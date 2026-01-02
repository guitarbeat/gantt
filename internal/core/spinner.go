package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a simple terminal spinner
type Spinner struct {
	msg      string
	stopChan chan bool
	wg       sync.WaitGroup
	active   bool
	mu       sync.Mutex
	silent   bool
}

// NewSpinner creates and starts a new spinner with the given message.
// silent: if true, suppresses all output
func NewSpinner(msg string, silent bool) *Spinner {
	s := &Spinner{
		msg:      msg,
		stopChan: make(chan bool),
		silent:   silent,
	}

	// Don't animate if silent or no color support
	if silent || !colorEnabled() {
		if !silent {
			fmt.Print(Info(msg + "... "))
		}
		return s
	}

	s.active = true
	s.wg.Add(1)
	go s.animate()
	return s
}

func (s *Spinner) animate() {
	defer s.wg.Done()

	// Braille patterns for smooth spinning
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.mu.Lock()
			msg := s.msg
			s.mu.Unlock()

			fmt.Printf("\r%s %s", colorize(Cyan, frames[i]), msg)
			i = (i + 1) % len(frames)
			time.Sleep(80 * time.Millisecond)
		}
	}
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		s.stopChan <- true
		s.wg.Wait()
		s.active = false
		// Clear the line
		fmt.Print("\r\033[K")
	}
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(msg string) {
	s.Stop()
	if !s.silent {
		if msg != "" {
			fmt.Printf("\r%s %s %s\n", Success("✔"), s.msg, msg)
		} else {
			fmt.Printf("\r%s %s\n", Success("✔"), s.msg)
		}
	}
}

// Fail stops the spinner and prints a failure message
func (s *Spinner) Fail(msg string) {
	s.Stop()
	if !s.silent {
		if msg != "" {
			fmt.Printf("\r%s %s %s\n", Error("✖"), s.msg, msg)
		} else {
			fmt.Printf("\r%s %s\n", Error("✖"), s.msg)
		}
	}
}

// Update updates the spinner message
func (s *Spinner) Update(msg string) {
	s.mu.Lock()
	s.msg = msg
	s.mu.Unlock()
}
