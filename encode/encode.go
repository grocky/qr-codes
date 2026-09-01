// Command encode generates QR code images for different payload types.
//
//	encode url [flags] <url>
//	encode wifi -ssid <ssid> [flags]
package main

import (
	"flag"
	"fmt"
	"os"

	"qr-codes/internal/encoder"
	"qr-codes/internal/payload"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "url":
		err = runURL(os.Args[2:])
	case "wifi":
		err = runWiFi(os.Args[2:])
	case "help", "-h", "-help", "--help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown type %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: encode <type> [flags]

types:
  url   encode a URL:                      encode url [flags] <url>
  wifi  encode WiFi network credentials:   encode wifi -ssid <ssid> [flags]

run "encode <type> -h" for the flags of each type.
`)
}

func runURL(args []string) error {
	fs := flag.NewFlagSet("encode url", flag.ExitOnError)
	opts, width := imageFlags(fs, "./qr-url.png")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: encode url [flags] <url>")
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

func runWiFi(args []string) error {
	fs := flag.NewFlagSet("encode wifi", flag.ExitOnError)
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
