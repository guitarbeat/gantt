package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a terminal spinner for loading states
type Spinner struct {
	mu       sync.Mutex
	msg      string
	frames   []string
	active   bool
	stopChan chan struct{}
	silent   bool
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(msg string, silent bool) *Spinner {
	return &Spinner{
		msg:      msg,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		active:   false,
		stopChan: make(chan struct{}),
		silent:   silent,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if s.silent || !colorEnabled() {
		if !s.silent {
			fmt.Printf("%s %s\n", Info("•"), s.msg)
		}
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				frame := s.frames[i%len(s.frames)]
				i++
				// Clear line and print frame + message
				fmt.Printf("\r%s %s", Info(frame), s.msg)
				s.mu.Unlock()
			}
		}
	}()
}

// Stop stops the spinner animation and clears the line
func (s *Spinner) Stop() {
	if s.silent || !colorEnabled() {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	close(s.stopChan)
	s.active = false
	// Clear the line
	fmt.Print("\r\033[K")
}

// Update updates the spinner message
func (s *Spinner) Update(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msg = msg
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(msg ...string) {
	s.Stop()
	if s.silent {
		return
	}

	finalMsg := s.msg
	if len(msg) > 0 {
		finalMsg = msg[0]
	}

	if colorEnabled() {
		fmt.Printf("\r%s %s\n", Success("✔"), finalMsg)
	} else {
		fmt.Printf("✔ %s\n", finalMsg)
	}
}

// Fail stops the spinner and prints a failure message
func (s *Spinner) Fail(msg ...string) {
	s.Stop()
	if s.silent {
		return
	}

	finalMsg := s.msg
	if len(msg) > 0 {
		finalMsg = msg[0]
	}

	if colorEnabled() {
		fmt.Printf("\r%s %s\n", Error("✖"), finalMsg)
	} else {
		fmt.Printf("✖ %s\n", finalMsg)
	}
}
