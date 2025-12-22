package core

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	stop   chan struct{}
	msg    string
	active bool
	wg     sync.WaitGroup
	mu     sync.Mutex
}

func NewSpinner(msg string, silent bool) *Spinner {
	if silent { return &Spinner{active: false} }
	s := &Spinner{stop: make(chan struct{}), msg: msg, active: true}
	fmt.Print("\033[?25l")
	s.wg.Add(1)
	go s.animate()
	return s
}

func (s *Spinner) animate() {
	defer s.wg.Done()
	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-s.stop:
			return
		default:
			fmt.Printf("\r%s %s ", Info(s.msg), chars[i])
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(chars)
		}
	}
}

func (s *Spinner) Stop(success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active { return }

	close(s.stop)
	s.wg.Wait()
	s.active = false
	fmt.Print("\033[?25h")
	mark := Error("❌")
	if success { mark = Success("✅") }
	fmt.Printf("\r%s %s\n", Info(s.msg), mark)
}
