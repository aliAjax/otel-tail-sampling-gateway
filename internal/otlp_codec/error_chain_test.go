package otlp_codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

var errInjectedRead = errors.New("injected read failure")

type failingReadCloser struct {
	prefix []byte
}

func (r *failingReadCloser) Read(p []byte) (int, error) {
	if len(r.prefix) == 0 {
		return 0, errInjectedRead
	}
	n := copy(p, r.prefix)
	r.prefix = r.prefix[n:]
	return n, nil
}

func (r *failingReadCloser) Close() error { return nil }

func TestDecodePreservesBodyReadError(t *testing.T) {
	req := &http.Request{Body: &failingReadCloser{prefix: []byte(`{"trace_id":`)}, Header: make(http.Header), ContentLength: -1}
	var dst map[string]any
	err := Decode(req, 1024, &dst)
	if !errors.Is(err, errInjectedRead) {
		t.Fatalf("body read error left the chain: %v", err)
	}
}

func TestDecodePreservesGzipHeaderError(t *testing.T) {
	req := &http.Request{Body: &failingReadCloser{}, Header: make(http.Header), ContentLength: -1}
	req.Header.Set("Content-Encoding", "gzip")
	var dst map[string]any
	err := Decode(req, 1024, &dst)
	if !errors.Is(err, errInjectedRead) {
		t.Fatalf("gzip reader error left the chain: %v", err)
	}
}

func TestDecodeRejectsTrailingDocument(t *testing.T) {
	req := &http.Request{Body: io.NopCloser(strings.NewReader(`{"trace_id":"a"}{"trace_id":"b"}`)), Header: make(http.Header), ContentLength: -1}
	var dst map[string]any
	if err := Decode(req, 1024, &dst); err == nil {
		t.Fatal("accepted a second JSON document")
	}
}

func TestDecodeFramePreservesHeaderError(t *testing.T) {
	_, err := DecodeFrame(&failingReadCloser{prefix: []byte("OT")}, 1024)
	if !errors.Is(err, errInjectedRead) {
		t.Fatalf("header error left the chain: %v", err)
	}
}

func framePrefix(version byte, payload []byte) []byte {
	b := bytes.NewBuffer(nil)
	b.Write([]byte{'O', 'T', 'L', 'P', version, 0})
	binary.Write(b, binary.BigEndian, uint32(len(payload)))
	b.Write(payload)
	return b.Bytes()
}

func TestDecodeFramePreservesPayloadError(t *testing.T) {
	header := framePrefix(1, []byte("ab"))[:10]
	binary.BigEndian.PutUint32(header[6:], 8)
	r := &failingReadCloser{prefix: append(header, []byte("ab")...)}
	_, err := DecodeFrame(r, 1024)
	if !errors.Is(err, errInjectedRead) {
		t.Fatalf("payload error left the chain: %v", err)
	}
}

func TestDecodeFramePreservesChecksumReadError(t *testing.T) {
	r := &failingReadCloser{prefix: framePrefix(1, []byte("payload"))}
	_, err := DecodeFrame(r, 1024)
	if !errors.Is(err, errInjectedRead) {
		t.Fatalf("checksum error left the chain: %v", err)
	}
}

func TestDecodeFrameRejectsUnsupportedVersion(t *testing.T) {
	frame := EncodeFrame([]byte("payload"))
	frame[4] = 9
	if _, err := DecodeFrame(bytes.NewReader(frame), 1024); err == nil {
		t.Fatal("accepted unsupported frame version")
	}
}
