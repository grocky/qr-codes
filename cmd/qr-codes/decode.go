package main

import (
	"flag"
	"fmt"
	"os"

	"qr-codes/internal/decoder"
)

func runDecode(args []string) error {
	fs := flag.NewFlagSet("qr-codes decode", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: qr-codes decode <image>...")
	}
	fs.Parse(args)
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(2)
	}

	for _, path := range fs.Args() {
		payloads, err := decoder.Decode(path)
		if err != nil {
			return err
		}
		if len(payloads) == 0 {
			fmt.Printf("%s: no QR codes found\n", path)
			continue
		}
		for _, p := range payloads {
			fmt.Printf("%s: %s\n", path, p)
		}
	}
	return nil
}
