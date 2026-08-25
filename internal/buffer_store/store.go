package buffer_store

import (
	"bufio"
	"encoding/json"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"os"
	"sync"
)

type Store struct {
	mu    sync.Mutex
	dir   string
	items []telemetry_domain.Span
	drops map[string]uint64
}

func New(dir string) *Store {
	os.MkdirAll(dir, 0750)
	return &Store{dir: dir, drops: map[string]uint64{}}
}
func (s *Store) Append(v telemetry_domain.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, v)
	f, e := os.OpenFile(s.dir+"/wal.jsonl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e == nil {
		defer f.Close()
		json.NewEncoder(f).Encode(v)
	}
	return e
}
func (s *Store) Drain(n int) []telemetry_domain.Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n <= 0 || n > len(s.items) {
		n = len(s.items)
	}
	o := append([]telemetry_domain.Span(nil), s.items[:n]...)
	s.items = s.items[n:]
	return o
}
func (s *Store) Len() int           { s.mu.Lock(); defer s.mu.Unlock(); return len(s.items) }
func (s *Store) Drop(reason string) { s.mu.Lock(); defer s.mu.Unlock(); s.drops[reason]++ }
func (s *Store) Drops() map[string]uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	o := map[string]uint64{}
	for k, v := range s.drops {
		o[k] = v
	}
	return o
}
func Recover(path string) ([]telemetry_domain.Span, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	var out []telemetry_domain.Span
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var s telemetry_domain.Span
		if json.Unmarshal(sc.Bytes(), &s) == nil {
			out = append(out, s)
		}
	}
	return out, sc.Err()
}
