package audit

import (
	"sync"
	"time"
)

type Event struct {
	RequestID, Tenant, Action, Result string
	At                                time.Time
}
type Log struct {
	mu     sync.Mutex
	events []Event
}

func (l *Log) Add(e Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e.At.IsZero() {
		e.At = time.Now()
	}
	l.events = append(l.events, e)
}
func (l *Log) List() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]Event(nil), l.events...)
}
