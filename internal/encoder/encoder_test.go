package encoder

import (
	"bytes"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/liyue201/goqr"
)

func TestEncodeRoundTrip(t *testing.T) {
	content := `WIFI:T:WPA;S:my\; net;P:s3cret;;`
	out := filepath.Join(t.TempDir(), "qr.png")

	if err := Encode(content, Options{OutputFile: out, QRWidth: 21}); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decoding image: %v", err)
	}
	if format != "png" {
		t.Errorf("output format = %q, want png", format)
	}

	codes, err := goqr.Recognize(img)
	if err != nil {
		t.Fatalf("recognizing QR code: %v", err)
	}
	if len(codes) != 1 {
		t.Fatalf("recognized %d QR codes, want 1", len(codes))
	}
	if got := string(codes[0].Payload); got != content {
		t.Errorf("decoded payload = %q, want %q", got, content)
	}
}

func TestEncodeMissingHalftoneImage(t *testing.T) {
	out := filepath.Join(t.TempDir(), "qr.png")
	err := Encode("https://example.com", Options{
		OutputFile:    out,
		HalftoneImage: filepath.Join(t.TempDir(), "nope.png"),
	})
	if err == nil {
		t.Fatal("Encode() succeeded, want error for missing halftone image")
	}
}
