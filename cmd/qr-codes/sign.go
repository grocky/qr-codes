package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"qr-codes/internal/sign"
)

func runSign(args []string) error {
	fs := flag.NewFlagSet("qr-codes sign", flag.ExitOnError)
	out := fs.String("o", "./out/wifi-sign.svg", "output SVG file")
	logoPath := fs.String("logo", "", "path to a PNG/JPEG/SVG logo image")

	p := sign.Params{ShowSuits: true}
	fs.StringVar(&p.SSID, "ssid", "", "network name (required)")
	fs.StringVar(&p.Password, "password", "", "network password")
	fs.StringVar(&p.Auth, "auth", "WPA", "auth type: WPA, WEP, or nopass")
	fs.BoolVar(&p.Hidden, "hidden", false, "network is hidden")
	fs.StringVar(&p.AccentColor, "accent", "", "accent color (#rrggbb)")
	fs.StringVar(&p.BackgroundColor, "background", "", "background color (#rrggbb)")
	fs.BoolVar(&p.ShowSuits, "suits", true, "show the card-suit ornament row")
	fs.StringVar(&p.Headline, "headline", "", "headline text")
	fs.StringVar(&p.Subtitle, "subtitle", "", "subtitle text")
	fs.StringVar(&p.FooterText, "footer", "Point your camera at the code — no typing required.", "footer text (empty hides the line)")
	fs.Parse(args)

	if *logoPath != "" {
		uri, err := logoDataURI(*logoPath)
		if err != nil {
			return err
		}
		p.LogoDataURI = uri
	}

	paramsJSON, err := json.Marshal(p)
	if err != nil {
		return err
	}
	var res sign.Result
	if err := json.Unmarshal(sign.Render(paramsJSON), &res); err != nil {
		return err
	}
	if len(res.Errors) > 0 {
		for _, e := range res.Errors {
			fmt.Fprintf(os.Stderr, "%s: %s\n", e.Field, e.Message)
		}
		return fmt.Errorf("invalid sign parameters")
	}

	if dir := filepath.Dir(*out); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(*out, []byte(res.SVG), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s (QR %dx%d modules)\n", *out, res.QRModules, res.QRModules)
	return nil
}

// logoDataURI reads an image file and encodes it as a base64 data URI.
func logoDataURI(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	mime := http.DetectContentType(data)
	if mime == "text/xml; charset=utf-8" || filepath.Ext(path) == ".svg" {
		mime = "image/svg+xml"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}
