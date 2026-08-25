package buffer_store

import (
	"encoding/json"
	"fmt"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Segment struct {
	ID      string
	Created time.Time
	Count   int
	Bytes   int64
	Path    string
}
type SegmentManager struct {
	mu       sync.Mutex
	dir      string
	segments map[string]Segment
}

func NewSegments(dir string) *SegmentManager {
	os.MkdirAll(dir, 0750)
	return &SegmentManager{dir: dir, segments: map[string]Segment{}}
}
func (m *SegmentManager) Write(batch []telemetry_domain.Span) (Segment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("seg-%d", time.Now().UnixNano())
	p := filepath.Join(m.dir, id+".json")
	f, e := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if e != nil {
		return Segment{}, e
	}
	enc := json.NewEncoder(f)
	for _, s := range batch {
		if e = enc.Encode(s); e != nil {
			f.Close()
			return Segment{}, e
		}
	}
	st, _ := f.Stat()
	f.Close()
	seg := Segment{ID: id, Created: time.Now(), Count: len(batch), Bytes: st.Size(), Path: p}
	m.segments[id] = seg
	return seg, nil
}
func (m *SegmentManager) List() []Segment {
	m.mu.Lock()
	defer m.mu.Unlock()
	o := make([]Segment, 0, len(m.segments))
	for _, s := range m.segments {
		o = append(o, s)
	}
	return o
}
func (m *SegmentManager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.segments[id]
	if !ok {
		return os.ErrNotExist
	}
	if e := os.Remove(s.Path); e != nil {
		return e
	}
	delete(m.segments, id)
	return nil
}
func (m *SegmentManager) Recover() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	es, e := os.ReadDir(m.dir)
	if e != nil {
		return e
	}
	for _, en := range es {
		if filepath.Ext(en.Name()) != ".json" {
			continue
		}
		st, e := en.Info()
		if e != nil {
			continue
		}
		id := en.Name()[:len(en.Name())-5]
		m.segments[id] = Segment{ID: id, Path: filepath.Join(m.dir, en.Name()), Bytes: st.Size()}
	}
	return nil
}
