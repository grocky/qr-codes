package decoder

import (
	"path/filepath"
	"testing"

	"qr-codes/internal/encoder"
)

func TestDecode(t *testing.T) {
	content := "https://example.com/page"
	path := filepath.Join(t.TempDir(), "qr.png")
	if err := encoder.Encode(content, encoder.Options{OutputFile: path, QRWidth: 21}); err != nil {
		t.Fatalf("generating test QR code: %v", err)
	}

	payloads, err := Decode(path)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}
	if len(payloads) != 1 || payloads[0] != content {
		t.Errorf("Decode() = %q, want [%q]", payloads, content)
	}
}

func TestDecodeMissingFile(t *testing.T) {
	if _, err := Decode(filepath.Join(t.TempDir(), "nope.png")); err == nil {
		t.Fatal("Decode() succeeded, want error for missing file")
	}
}
