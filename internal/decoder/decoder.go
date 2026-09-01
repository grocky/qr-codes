package decoder

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/liyue201/goqr"
)

// Decode reads an image file and returns the payloads of every QR code found in it.
func Decode(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image %s: %w", path, err)
	}
	codes, err := goqr.Recognize(img)
	if err != nil {
		return nil, fmt.Errorf("recognize QR codes in %s: %w", path, err)
	}
	payloads := make([]string, len(codes))
	for i, code := range codes {
		payloads[i] = string(code.Payload)
	}
	return payloads, nil
}
