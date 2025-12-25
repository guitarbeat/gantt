package core

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	mu   sync.Mutex
	wg   sync.WaitGroup
	stop chan struct{}
	msg  string
	on   bool
	mute bool
}

func NewSpinner(msg string, silent bool) *Spinner {
	return &Spinner{stop: make(chan struct{}), msg: msg, mute: silent}
}

func (s *Spinner) Start() {
	if s.mute || s.on { return }
	s.mu.Lock(); s.on = true; s.mu.Unlock()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for i := 0; ; i++ {
			select {
			case <-s.stop: return
			case <-t.C:
				s.mu.Lock()
				fmt.Printf("\r\033[36m%s\033[0m %s\033[K", chars[i%10], s.msg)
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Spinner) Update(msg string) {
	if s.mute { return }
	s.mu.Lock(); s.msg = msg; s.mu.Unlock()
}

func (s *Spinner) Stop(msg string) {
	if s.mute { return }
	s.mu.Lock()
	if !s.on { s.mu.Unlock(); return }
	s.on = false
	close(s.stop)
	s.mu.Unlock()
	s.wg.Wait()
	fmt.Printf("\r\033[K%s", msg)
	if msg != "" { fmt.Println() }
}
