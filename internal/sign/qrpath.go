package sign

import (
	"fmt"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
)

// matrixWriter captures the QR module matrix from qrcode.Save instead of
// rendering an image.
type matrixWriter struct {
	bitmap [][]bool
}

func (w *matrixWriter) Write(mat qrcode.Matrix) error {
	w.bitmap = mat.Bitmap()
	return nil
}

func (w *matrixWriter) Close() error { return nil }

// qrBitmap encodes content as a QR code and returns its module bitmap
// (true = dark), without any quiet zone.
func qrBitmap(content string) ([][]bool, error) {
	qrc, err := qrcode.New(content)
	if err != nil {
		return nil, fmt.Errorf("build QR code: %w", err)
	}
	var w matrixWriter
	if err := qrc.Save(&w); err != nil {
		return nil, fmt.Errorf("extract QR matrix: %w", err)
	}
	return w.bitmap, nil
}

// pathData converts a QR module bitmap into SVG path data, one unit per
// module, merging horizontal runs of dark modules into single subpaths.
func pathData(bitmap [][]bool) string {
	var b strings.Builder
	for y, row := range bitmap {
		for x := 0; x < len(row); {
			if !row[x] {
				x++
				continue
			}
			run := 0
			for x+run < len(row) && row[x+run] {
				run++
			}
			fmt.Fprintf(&b, "M%d %dh%dv1h-%dz", x, y, run, run)
			x += run
		}
	}
	return b.String()
}
