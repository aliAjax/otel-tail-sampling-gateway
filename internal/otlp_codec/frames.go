package otlp_codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"time"
)

type Frame struct {
	Version  byte
	Flags    byte
	Payload  []byte
	Checksum uint32
}

func EncodeFrame(payload []byte) []byte {
	b := bytes.NewBuffer(nil)
	b.Write([]byte{'O', 'T', 'L', 'P', 1, 0})
	binary.Write(b, binary.BigEndian, uint32(len(payload)))
	b.Write(payload)
	binary.Write(b, binary.BigEndian, crc32.ChecksumIEEE(payload))
	return b.Bytes()
}
func DecodeFrame(r io.Reader, max uint32) (Frame, error) {
	var h [10]byte
	if _, e := io.ReadFull(r, h[:]); e != nil {
		return Frame{}, fmt.Errorf("frame_header: %v", e)
	}
	if string(h[:4]) != "OTLP" {
		return Frame{}, fmt.Errorf("frame_magic")
	}
	n := binary.BigEndian.Uint32(h[6:])
	if n > max {
		return Frame{}, fmt.Errorf("frame_too_large")
	}
	p := make([]byte, n)
	if _, e := io.ReadFull(r, p); e != nil {
		return Frame{}, fmt.Errorf("frame_partial: %v", e)
	}
	var c uint32
	if e := binary.Read(r, binary.BigEndian, &c); e != nil {
		return Frame{}, fmt.Errorf("frame_checksum: %v", e)
	}
	if crc32.ChecksumIEEE(p) != c {
		return Frame{}, fmt.Errorf("frame_checksum_mismatch")
	}
	return Frame{Version: h[4], Flags: h[5], Payload: p, Checksum: c}, nil
}
func ValidateVersion(v byte) error {
	if v == 0 || v > 2 {
		return fmt.Errorf("unsupported_version")
	}
	return nil
}
func SplitBatch(data []byte, size int) [][]byte {
	if size < 1 {
		size = 1
	}
	var out [][]byte
	for len(data) > 0 {
		n := size
		if n > len(data) {
			n = len(data)
		}
		out = append(out, append([]byte(nil), data[:n]...))
		data = data[n:]
	}
	return out
}
func CompressionRatio(raw, compressed int) float64 {
	if raw <= 0 {
		return 0
	}
	return float64(compressed) / float64(raw)
}
func Age(ts time.Time) time.Duration {
	if ts.IsZero() {
		return 0
	}
	return time.Since(ts)
}
