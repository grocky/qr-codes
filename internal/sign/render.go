package sign

import (
	_ "embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"qr-codes/internal/payload"
)

//go:embed template.svg.tmpl
var templateSVG string

// text/template, not html/template: the contextual escaper of html/template
// mangles SVG attributes. All user strings are escaped with escapeXML before
// they reach the template.
var signTemplate = template.Must(template.New("sign").Parse(templateSVG))

// Result is what Render returns, JSON-encoded: either SVG or Errors is set.
type Result struct {
	SVG       string       `json:"svg,omitempty"`
	QRModules int          `json:"qrModules,omitempty"`
	Errors    []FieldError `json:"errors,omitempty"`
}

// Render turns params JSON into a Result JSON. It never returns an error:
// every failure is reported inside the Result so the caller (the web form)
// can present it.
func Render(paramsJSON []byte) []byte {
	res := buildResult(paramsJSON)
	out, err := json.Marshal(res)
	if err != nil {
		// Result contains only strings and ints; this cannot happen.
		return []byte(`{"errors":[{"field":"","message":"internal: encoding result failed"}]}`)
	}
	return out
}

func buildResult(paramsJSON []byte) Result {
	var p Params
	if err := json.Unmarshal(paramsJSON, &p); err != nil {
		return Result{Errors: []FieldError{{Field: "", Message: "invalid parameters: " + err.Error()}}}
	}
	p.ApplyDefaults()
	if errs := p.Validate(); len(errs) > 0 {
		return Result{Errors: errs}
	}

	auth, _ := payload.ParseAuth(p.Auth) // validated above
	content, err := payload.WiFi{
		SSID:     p.SSID,
		Password: p.Password,
		Auth:     auth,
		Hidden:   p.Hidden,
	}.Build()
	if err != nil {
		return Result{Errors: []FieldError{{Field: "ssid", Message: err.Error()}}}
	}

	bitmap, err := qrBitmap(content)
	if err != nil {
		return Result{Errors: []FieldError{{Field: "", Message: err.Error()}}}
	}
	n := len(bitmap)

	headline := fitText(p.Headline, 94, 4, 0.72, 700, 40)
	subtitle := fitText(p.Subtitle, 26, 9, 0.62, 700, 12)
	footer := fitText(p.FooterText, 20, 0.5, 0.62, 780, 12)
	tagline := fitText(p.Tagline, 24, 2, 0.62, 700, 12)

	data := struct {
		Accent, FeltInner, FeltOuter                        string
		LogoDataURI                                         string
		Ornament                                            string
		Headline, Subtitle, FooterText, TaglineText         string
		HeadlineSize, SubtitleSize, FooterSize, TaglineSize string
		HeadlineFit, SubtitleFit, FooterFit, TaglineFit     string
		QRScale, QRPath                                     string
	}{
		Accent:       p.AccentColor,
		FeltInner:    p.BackgroundColor,
		FeltOuter:    darken(p.BackgroundColor),
		LogoDataURI:  escapeXML(p.LogoDataURI),
		Ornament:     p.Ornament,
		Headline:     escapeXML(p.Headline),
		Subtitle:     escapeXML(p.Subtitle),
		FooterText:   escapeXML(p.FooterText),
		TaglineText:  escapeXML(p.Tagline),
		HeadlineSize: headline.size,
		SubtitleSize: subtitle.size,
		FooterSize:   footer.size,
		TaglineSize:  tagline.size,
		HeadlineFit:  headline.fit,
		SubtitleFit:  subtitle.fit,
		FooterFit:    footer.fit,
		TaglineFit:   tagline.fit,
		QRScale:      formatNum(qrBoxSize / float64(n)),
		QRPath:       pathData(bitmap),
	}

	var b strings.Builder
	if err := signTemplate.Execute(&b, data); err != nil {
		return Result{Errors: []FieldError{{Field: "", Message: "internal: template failed: " + err.Error()}}}
	}
	return Result{SVG: b.String(), QRModules: n}
}

// qrBoxSize is the QR area in template units (344×344 at 253,610); the white
// card around it provides the quiet zone.
const qrBoxSize = 344

type fittedText struct {
	size string // font-size to emit
	fit  string // extra attributes when the shrink floor was hit, else ""
}

// fitText shrinks a line's font-size so its estimated width fits maxWidth.
// Width is estimated as runes*size*k + (runes-1)*letterSpacing. If even the
// floor size overflows, a textLength clamp is emitted as a hard stop.
func fitText(s string, baseSize, letterSpacing, k, maxWidth, floor float64) fittedText {
	runes := float64(len([]rune(s)))
	if runes == 0 {
		return fittedText{size: formatNum(baseSize)}
	}
	size := baseSize
	width := func(sz float64) float64 { return runes*sz*k + (runes-1)*letterSpacing }
	if width(size) > maxWidth {
		size = (maxWidth - (runes-1)*letterSpacing) / (runes * k)
	}
	if size < floor {
		return fittedText{
			size: formatNum(floor),
			fit:  fmt.Sprintf(` textLength="%s" lengthAdjust="spacingAndGlyphs"`, formatNum(maxWidth)),
		}
	}
	return fittedText{size: formatNum(size)}
}

// darken halves each RGB channel, matching the original felt gradient's
// #1b6f4c → #0a3826 inner/outer ratio. hex is a validated #rrggbb string.
func darken(hex string) string {
	v, _ := strconv.ParseUint(hex[1:], 16, 32)
	r, g, b := (v>>16)&0xff/2, (v>>8)&0xff/2, v&0xff/2
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func formatNum(f float64) string {
	return strconv.FormatFloat(f, 'g', 6, 64)
}

func escapeXML(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
