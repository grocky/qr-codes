package encoder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// Options control how the QR code image is rendered.
type Options struct {
	OutputFile    string
	HalftoneImage string // optional image blended into the QR modules
	Transparent   bool
	QRWidth       uint8 // module width in pixels; 0 uses the writer default
}

// Encode renders content as a QR code image according to opts.
func Encode(content string, opts Options) error {
	qrc, err := qrcode.New(content)
	if err != nil {
		return fmt.Errorf("build QR code: %w", err)
	}

	// The writer encodes JPEG by default regardless of the output filename;
	// pick the encoder that matches the extension.
	format := standard.JPEG_FORMAT
	if strings.EqualFold(filepath.Ext(opts.OutputFile), ".png") {
		format = standard.PNG_FORMAT
	}
	imageOptions := []standard.ImageOption{standard.WithBuiltinImageEncoder(format)}
	if opts.HalftoneImage != "" {
		if _, err := os.Stat(opts.HalftoneImage); err != nil {
			return fmt.Errorf("halftone image not found: %s", opts.HalftoneImage)
		}
		imageOptions = append(imageOptions, standard.WithHalftone(opts.HalftoneImage))
	}
	if opts.QRWidth > 0 {
		imageOptions = append(imageOptions, standard.WithQRWidth(opts.QRWidth))
	}
	if opts.Transparent {
		imageOptions = append(imageOptions, standard.WithBgTransparent())
	}

	if dir := filepath.Dir(opts.OutputFile); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	writer, err := standard.New(opts.OutputFile, imageOptions...)
	if err != nil {
		return fmt.Errorf("create writer: %w", err)
	}
	defer writer.Close()

	if err := qrc.Save(writer); err != nil {
		return fmt.Errorf("save QR code: %w", err)
	}
	return nil
}
