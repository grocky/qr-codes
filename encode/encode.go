package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var (
	transparent = flag.Bool("transparent", false, "set background to transparent")
)

func main() {

	flag.Parse()

	qrc, err := qrcode.New("https://www.meetup.com/northern-virginia-poker-club/events/307974862/")
	handleErr(err)

	halftonePath := "./images/norther-virginia-poker-club-logo_no-bg.png"
	if _, err := os.Stat(halftonePath); os.IsNotExist(err) {
		fmt.Printf("halftone image not found: %s\n", halftonePath)
		return
	}

	options := []standard.ImageOption{
		standard.WithHalftone(halftonePath),
		standard.WithQRWidth(21),
	}

	outputFilename := "./halftone-qr.png"

	if *transparent {
		options = append(
			options,
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			standard.WithBgTransparent(),
		)
		outputFilename = "./halftone-qr-transparent.png"
	}

	writer, err := standard.New(outputFilename, options...)
	handleErr(err)

	defer writer.Close()
	if err := qrc.Save(writer); err != nil {
		fmt.Printf("failed to save QR code: %v\n", err)
	}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
