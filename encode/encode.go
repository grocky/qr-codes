package main

import (
	"flag"

	"github.com/yeqown/go-qrcode"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var (
	transparent = flag.Bool("transparent", false, "set background to transparent")
)

func main() {

	flag.Parse()

	qrc, err := qrcode.New("https://www.meetup.com/northern-virginia-poker-club/events/307974862/")
	if err != nil {
		panic(err)
	}

	options := []standard.ImageOption{
		standard.WithHalftone("../images/norther-virginia-poker-club-logo_no-bg.png"),
		standard.WithQRWidth(21),
	}
	filename := "./halftone-qr.png"

	if *transparent {
		options = append(
			options,
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			standard.WithBgTransparent(),
		)
		filename = "./halftone-qr-transparent.png"
	}

	w0, err := standard.New(filename, options...)
	handleErr(err)
	err = qrc.Save(w0)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
