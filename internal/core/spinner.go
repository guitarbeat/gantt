package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a CLI loading spinner
type Spinner struct {
	message string
	stop    chan bool
	wg      sync.WaitGroup
	active  bool
}

// NewSpinner creates a new spinner instance
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if IsSilent() {
		return
	}

	// If not interactive, just print the message once without animation
	if !IsInteractive() {
		fmt.Printf("%s...\n", s.message)
		return
	}

	s.active = true
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
				fmt.Print(ClearLine())
				fmt.Printf("%s %s", CyanText(chars[i]), Info(s.message))
				time.Sleep(100 * time.Millisecond)
				i = (i + 1) % len(chars)
			}
		}
	}()
}

// Stop ends the spinner animation
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.stop <- true
	s.wg.Wait()
	// Leave the clearing/final status message to the caller
}
