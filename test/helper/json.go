package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

// EncodeJSON encodes a struct into JSON and returns a bytes.Buffer.
// If encoding fails, it fails the test.
func EncodeJSON(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		t.Fatalf("Failed to encode JSON: %v", err)
	}
	return buf
}

// DecodeJSON decodes JSON from a reader into the given struct.
// If decoding fails, it fails the test.
func DecodeJSON(t *testing.T, r io.ReadCloser, v any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
}
