package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a terminal spinner for long-running operations
type Spinner struct {
	mu       sync.Mutex
	message  string
	chars    []string
	delay    time.Duration
	active   bool
	stopChan chan struct{}
	silent   bool
	wg       sync.WaitGroup
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string, silent bool) *Spinner {
	return &Spinner{
		message:  message,
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:    100 * time.Millisecond,
		stopChan: make(chan struct{}),
		silent:   silent,
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

	// Hide cursor
	fmt.Print("\033[?25l")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.mu.Lock()
				// Check active state again inside lock to avoid race if stopped
				if !s.active {
					s.mu.Unlock()
					return
				}
				char := s.chars[i%len(s.chars)]
				fmt.Printf("\r%s %s", Info(char), s.message)
				s.mu.Unlock()
				time.Sleep(s.delay)
				i++
			}
		}
	}()
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// Stop ends the spinner animation and prints a completion status
func (s *Spinner) Stop(success bool) {
	if s.silent {
		return
	}

	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}

	// Signal goroutine to stop
	close(s.stopChan)
	s.active = false
	s.mu.Unlock()

	// Wait for goroutine to finish
	s.wg.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Restore cursor
	fmt.Print("\033[?25h")

	// clear line
	fmt.Print("\r\033[K")

	if success {
		fmt.Printf("%s %s\n", Success("✅"), s.message)
	} else {
		fmt.Printf("%s %s\n", Error("❌"), s.message)
	}
}

// StopWithFailure ends the spinner with a failure message
func (s *Spinner) StopWithFailure(errMessage string) {
	if s.silent {
		return
	}

	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}

	// Signal goroutine to stop
	close(s.stopChan)
	s.active = false
	s.mu.Unlock()

	// Wait for goroutine to finish
	s.wg.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Restore cursor
	fmt.Print("\033[?25h")

	// clear line
	fmt.Print("\r\033[K")

	fmt.Printf("%s %s (%s)\n", Error("❌"), s.message, errMessage)
}
