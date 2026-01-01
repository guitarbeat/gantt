package core

import (
	"fmt"
	"time"
)

// Spinner is a simple loading indicator
type Spinner struct {
	msg  string
	done chan bool
}

// NewSpinner creates and starts a spinner
func NewSpinner(msg string) *Spinner {
	s := &Spinner{msg: msg, done: make(chan bool)}
	if !colorEnabled() {
		fmt.Print(msg + "... ")
		return s
	}
	go s.spin()
	return s
}

func (s *Spinner) spin() {
	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-s.done:
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("\r%s %s... ", colorize(Cyan, chars[i]), s.msg)
			i = (i + 1) % len(chars)
		}
	}
}

// Stop stops the spinner and prints a final status
func (s *Spinner) Success(mark string) {
	s.stop()
	if !colorEnabled() {
		fmt.Println(mark)
		return
	}
	fmt.Printf("\r\033[K%s %s... %s\n", Success("✔"), s.msg, mark)
}

func (s *Spinner) Fail(mark string) {
	s.stop()
	if !colorEnabled() {
		fmt.Println(mark)
		return
	}
	fmt.Printf("\r\033[K%s %s... %s\n", Error("✘"), s.msg, mark)
}

func (s *Spinner) stop() {
	// If colorEnabled() was false, spin() is not running
	if colorEnabled() {
		s.done <- true
	}
}
