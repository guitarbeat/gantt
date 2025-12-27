package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner provides a visual indicator for long-running operations.
// It uses a goroutine to animate a spinner character until Stop is called.
type Spinner struct {
	msg    string
	silent bool
	stop   chan struct{}
	wg     sync.WaitGroup
}

// NewSpinner creates a new spinner instance with the given message.
// If silent is true, the spinner will not produce any output.
func NewSpinner(msg string, silent bool) *Spinner {
	return &Spinner{
		msg:    msg,
		silent: silent,
		stop:   make(chan struct{}),
	}
}

// Start begins the spinner animation in a background goroutine.
func (s *Spinner) Start() {
	if s.silent {
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// Braille patterns for smooth spinning
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				// Print frame and message, using \r to overwrite the line
				fmt.Printf("\r%s %s... ", frames[i%len(frames)], Info(s.msg))
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()
}

// Stop ends the spinner animation and prints a success checkmark.
func (s *Spinner) Stop() {
	if s.silent {
		return
	}

	close(s.stop)
	s.wg.Wait()

	// Overwrite the spinner line with the success state
	fmt.Printf("\r%s %s... \n", Success("✅"), Info(s.msg))
}

// Fail ends the spinner animation and prints a failure cross.
func (s *Spinner) Fail() {
	if s.silent {
		return
	}

	close(s.stop)
	s.wg.Wait()

	// Overwrite the spinner line with the failure state
	fmt.Printf("\r%s %s... \n", Error("❌"), Info(s.msg))
}
