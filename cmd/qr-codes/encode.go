package main

import (
	"flag"
	"fmt"
	"os"

	"qr-codes/internal/encoder"
	"qr-codes/internal/payload"
)

func runEncode(args []string) error {
	if len(args) < 1 {
		encodeUsage()
		os.Exit(2)
	}

	switch args[0] {
	case "url":
		return runEncodeURL(args[1:])
	case "wifi":
		return runEncodeWiFi(args[1:])
	case "help", "-h", "-help", "--help":
		encodeUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "unknown type %q\n\n", args[0])
		encodeUsage()
		os.Exit(2)
		return nil
	}
}

func encodeUsage() {
	fmt.Fprint(os.Stderr, `usage: qr-codes encode <type> [flags]

types:
  url   encode a URL:                      qr-codes encode url [flags] <url>
  wifi  encode WiFi network credentials:   qr-codes encode wifi -ssid <ssid> [flags]

run "qr-codes encode <type> -h" for the flags of each type.
`)
}

func runEncodeURL(args []string) error {
	fs := flag.NewFlagSet("qr-codes encode url", flag.ExitOnError)
	opts, width := imageFlags(fs, "./qr-url.png")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: qr-codes encode url [flags] <url>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}

	content, err := payload.URL{URL: fs.Arg(0)}.Build()
	if err != nil {
		return err
	}
	return encode(content, opts, *width)
}

func runEncodeWiFi(args []string) error {
	fs := flag.NewFlagSet("qr-codes encode wifi", flag.ExitOnError)
	opts, width := imageFlags(fs, "./qr-wifi.png")
	ssid := fs.String("ssid", "", "network name (required)")
	password := fs.String("password", "", "network password")
	auth := fs.String("auth", "WPA", "auth type: WPA, WEP, or nopass")
	hidden := fs.Bool("hidden", false, "network is hidden")
	fs.Parse(args)

	authType, err := payload.ParseAuth(*auth)
	if err != nil {
		return err
	}
	content, err := payload.WiFi{
		SSID:     *ssid,
		Password: *password,
		Auth:     authType,
		Hidden:   *hidden,
	}.Build()
	if err != nil {
		return err
	}
	return encode(content, opts, *width)
}

// imageFlags registers the rendering flags shared by every payload type.
func imageFlags(fs *flag.FlagSet, defaultOutput string) (*encoder.Options, *uint) {
	opts := &encoder.Options{}
	fs.StringVar(&opts.OutputFile, "o", defaultOutput, "output image file")
	fs.StringVar(&opts.HalftoneImage, "halftone", "", "path to an image to blend into the QR modules")
	fs.BoolVar(&opts.Transparent, "transparent", false, "set background to transparent")
	width := fs.Uint("width", 21, "QR module width in pixels")
	return opts, width
}

func encode(content string, opts *encoder.Options, width uint) error {
	if width > 255 {
		return fmt.Errorf("width %d out of range (max 255)", width)
	}
	opts.QRWidth = uint8(width)

	if err := encoder.Encode(content, *opts); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", opts.OutputFile)
	return nil
}
