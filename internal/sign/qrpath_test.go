package sign

import (
	"image"
	"image/color"
	"testing"

	"github.com/liyue201/goqr"

	"qr-codes/internal/payload"
)

// TestQRBitmapRoundTrip proves the generated module bitmap encodes the exact
// payload by painting it into an image and decoding it with an independent
// QR reader.
func TestQRBitmapRoundTrip(t *testing.T) {
	content, err := payload.WiFi{SSID: "my home; net", Password: "p@ss:word,1", Auth: payload.AuthWPA}.Build()
	if err != nil {
		t.Fatal(err)
	}

	bitmap, err := qrBitmap(content)
	if err != nil {
		t.Fatalf("qrBitmap() error: %v", err)
	}
	n := len(bitmap)
	if n < 21 || n%2 == 0 {
		t.Fatalf("bitmap is %dx%d, want odd n >= 21", n, n)
	}

	// Paint at 4px per module with a 4-module quiet zone.
	const scale, quiet = 4, 4
	size := (n + 2*quiet) * scale
	img := image.NewGray(image.Rect(0, 0, size, size))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	for y, row := range bitmap {
		for x, dark := range row {
			if !dark {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetGray((x+quiet)*scale+dx, (y+quiet)*scale+dy, color.Gray{Y: 0})
				}
			}
		}
	}

	codes, err := goqr.Recognize(img)
	if err != nil {
		t.Fatalf("recognize: %v", err)
	}
	if len(codes) != 1 {
		t.Fatalf("recognized %d codes, want 1", len(codes))
	}
	if got := string(codes[0].Payload); got != content {
		t.Errorf("decoded %q, want %q", got, content)
	}
}

func TestPathData(t *testing.T) {
	tests := []struct {
		name   string
		bitmap [][]bool
		want   string
	}{
		{
			name:   "single module",
			bitmap: [][]bool{{true}},
			want:   "M0 0h1v1h-1z",
		},
		{
			name: "horizontal run merges",
			bitmap: [][]bool{
				{true, true, true},
				{false, false, false},
				{true, false, true},
			},
			want: "M0 0h3v1h-3zM0 2h1v1h-1zM2 2h1v1h-1z",
		},
		{
			name:   "empty bitmap",
			bitmap: [][]bool{{false, false}, {false, false}},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathData(tt.bitmap); got != tt.want {
				t.Errorf("pathData() = %q, want %q", got, tt.want)
			}
		})
	}
}
