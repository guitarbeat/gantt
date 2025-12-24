package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a terminal spinner for long-running operations
type Spinner struct {
	msg      string
	silent   bool
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
}

// NewSpinner creates a new spinner instance
func NewSpinner(msg string, silent bool) *Spinner {
	return &Spinner{
		msg:      msg,
		silent:   silent,
		stopChan: make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if s.silent {
		return
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			case <-time.After(100 * time.Millisecond):
				s.mu.Lock()
				fmt.Printf("\r%s %s", Info(chars[i]), s.msg)
				s.mu.Unlock()
				i = (i + 1) % len(chars)
			}
		}
	}()
}

// UpdateMessage updates the message displayed next to the spinner
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msg = msg
}

// Stop ends the spinner animation and prints a success or error symbol
func (s *Spinner) Stop(success bool) {
	if s.silent {
		return
	}
	close(s.stopChan)
	s.wg.Wait()

	// Clear the line
	fmt.Print("\r\033[K")

	if success {
		fmt.Printf("\r%s %s\n", Success("✔"), s.msg)
	} else {
		fmt.Printf("\r%s %s\n", Error("✖"), s.msg)
	}
}
