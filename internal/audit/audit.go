package audit

import (
	"sync"
	"time"
)

type Event struct {
	RequestID, Tenant, Action, Result string
	At                                time.Time
	Details                           map[string]string
	Labels                            []string
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
	l.events = append(l.events, storeEvent(e))
}

func storeEvent(event Event) Event {
	return event
}

func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.events)
}

func (l *Log) List() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	for i, event := range l.events {
		out[i] = snapshotEvent(event)
	}
	return out
}

func snapshotEvent(event Event) Event {
	return event
}
