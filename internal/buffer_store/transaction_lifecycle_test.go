package buffer_store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func TestAppendRollbackOnWALError(t *testing.T) {
	dir := t.TempDir()
	store := New(dir)
	store.dir = filepath.Join(dir, "missing", "child")
	if err := store.Append(telemetry_domain.Span{TraceID: "must-not-commit"}); err == nil {
		t.Fatal("expected WAL open failure")
	}
	if got := store.Len(); got != 0 {
		t.Fatalf("memory committed before WAL: len=%d", got)
	}
}

func TestSegmentWriteRejectsEmptyBatch(t *testing.T) {
	dir := t.TempDir()
	m := NewSegments(dir)
	if _, err := m.Write(nil); err == nil {
		t.Fatal("expected empty segment rejection")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 || len(m.List()) != 0 {
		t.Fatalf("empty segment published: files=%d index=%d", len(entries), len(m.List()))
	}
}

func TestSegmentRecoverRejectsCorruptFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "seg-bad.json"), []byte("{broken\n"), 0600); err != nil {
		t.Fatal(err)
	}
	m := NewSegments(dir)
	m.segments["existing"] = Segment{ID: "existing"}
	if err := m.Recover(); err == nil {
		t.Fatal("expected corrupt segment error")
	}
	list := m.List()
	if len(list) != 1 || list[0].ID != "existing" {
		t.Fatalf("failed recovery changed published index: %#v", list)
	}
}

func TestSegmentRecoverRestoresCount(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "seg-valid.json")
	data := []byte("{\"trace_id\":\"a\"}\n{\"trace_id\":\"b\"}\n")
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatal(err)
	}
	m := NewSegments(dir)
	if err := m.Recover(); err != nil {
		t.Fatal(err)
	}
	list := m.List()
	if len(list) != 1 || list[0].Count != 2 || list[0].Bytes != int64(len(data)) {
		t.Fatalf("recovered metadata is incomplete: %#v", list)
	}
}
