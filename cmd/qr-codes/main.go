// Command qr-codes generates and reads QR code images.
//
//	qr-codes encode url [flags] <url>
//	qr-codes encode wifi -ssid <ssid> [flags]
//	qr-codes decode <image>...
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "encode":
		err = runEncode(os.Args[2:])
	case "decode":
		err = runDecode(os.Args[2:])
	case "sign":
		err = runSign(os.Args[2:])
	case "help", "-h", "-help", "--help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: qr-codes <command> [arguments]

commands:
  encode  generate a QR code image:      qr-codes encode <type> [flags]
  decode  read QR codes from images:     qr-codes decode <image>...
  sign    generate a Wi-Fi sign SVG:     qr-codes sign -ssid <ssid> [flags]

run "qr-codes <command> -h" for details.
`)
}
