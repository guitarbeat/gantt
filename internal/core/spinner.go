package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a terminal loading animation
type Spinner struct {
	mu      sync.Mutex
	stop    chan bool
	message string
	active  bool
	silent  bool
}

// NewSpinner creates a new spinner with the given message.
// The silent flag overrides any environment settings if true.
// If silent is false, it still checks IsSilent() from logger configuration.
func NewSpinner(message string, silent bool) *Spinner {
	return &Spinner{
		stop:    make(chan bool),
		message: message,
		silent:  silent || IsSilent(),
	}
}

// Start begins the spinner animation in a background goroutine
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

	// Hide cursor
	fmt.Print("\033[?25l")

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				s.mu.Lock()
				msg := s.message
				s.mu.Unlock()
				// \r moves cursor to start of line, \033[K clears the line
				fmt.Printf("\r%s %s", Info(frames[i%len(frames)]), msg)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	if s.silent {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.stop <- true
	s.active = false
	// Show cursor again
	fmt.Print("\033[?25h")
	// Clear the line
	fmt.Print("\r\033[K")
}

// Success stops the spinner and prints a success message with checkmark
func (s *Spinner) Success(msg string) {
	s.Stop()
	if !s.silent {
		fmt.Printf("\r%s %s\n", Success("✅"), msg)
	}
}

// Fail stops the spinner and prints an error message
func (s *Spinner) Fail(msg string) {
	s.Stop()
	if !s.silent {
		fmt.Printf("\r%s %s\n", Error("❌"), msg)
	}
}

// UpdateMessage updates the spinner text while it's running
func (s *Spinner) UpdateMessage(msg string) {
	if s.silent {
		return
	}
	s.mu.Lock()
	s.message = msg
	s.mu.Unlock()
}
