package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a loading indicator
type Spinner struct {
	message string
	stop    chan bool
	wg      sync.WaitGroup
	silent  bool
}

// NewSpinner creates a new spinner instance
func NewSpinner(message string, silent bool) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan bool),
		silent:  silent,
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
			case <-s.stop:
				return
			default:
				fmt.Printf("\r%s %s", Info(chars[i]), s.message)
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop ends the spinner with a success message
func (s *Spinner) Stop() {
	if s.silent {
		return
	}
	s.stop <- true
	s.wg.Wait()
	fmt.Printf("\r%s %s\n", Success("✅"), s.message)
}

// Fail ends the spinner with an error message
func (s *Spinner) Fail(err error) {
	if s.silent {
		return
	}
	s.stop <- true
	s.wg.Wait()
	fmt.Printf("\r%s %s: %v\n", Error("❌"), s.message, err)
}
