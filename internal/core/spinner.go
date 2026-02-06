// Package core - Spinner provides a simple CLI spinner for long-running operations.
package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a CLI loading spinner
type Spinner struct {
	mu      sync.Mutex
	msg     string
	frames  []string
	active  bool
	stop    chan bool
	silent  bool
	wg      sync.WaitGroup
}

// NewSpinner creates a new spinner with the given message.
// If silent is true, the spinner will not output anything.
func NewSpinner(msg string, silent bool) *Spinner {
	return &Spinner{
		msg:    msg,
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:   make(chan bool),
		silent: silent,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if s.silent {
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				frame := s.frames[i%len(s.frames)]
				i++
				// Print current frame and message, then carriage return
				// We use \033[K to clear the rest of the line
				fmt.Printf("\r%s %s %s\033[K", Info(frame), s.msg, DimText("..."))
				s.mu.Unlock()
			}
		}
	}()
}

// Stop ends the spinner animation
func (s *Spinner) Stop(success bool) {
	if s.silent {
		return
	}

	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	// Signal the goroutine to stop
	close(s.stop)
	s.wg.Wait() // Wait for goroutine to finish

	// Clear the line and print final status
	fmt.Print("\r\033[K") // Clear line
	if success {
		fmt.Printf("%s %s\n", Success("✅"), s.msg)
	} else {
		fmt.Printf("%s %s\n", Error("❌"), s.msg)
	}
}

// UpdateMessage changes the spinner message
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msg = msg
}
