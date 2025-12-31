package core

import (
	"fmt"
	"sync"
	"time"
)

// Spinner provides a terminal loading animation
type Spinner struct {
	msg    string
	stop   chan struct{}
	wg     sync.WaitGroup
	active bool
	mu     sync.Mutex
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(msg string) *Spinner {
	return &Spinner{
		msg:  msg,
		stop: make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if !colorEnabled() {
		fmt.Println(s.msg + "...")
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
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				fmt.Printf("\r%s %s", CyanText(chars[i]), s.msg)
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop ends the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return
	}

	close(s.stop)
	s.wg.Wait()
	s.active = false
	fmt.Print("\r\033[K") // Clear line
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(msg string) {
	s.Stop()
	fmt.Println(Success("✅ " + msg))
}

// Fail stops the spinner and prints an error message
func (s *Spinner) Fail(msg string) {
	s.Stop()
	fmt.Println(Error("❌ " + msg))
}
