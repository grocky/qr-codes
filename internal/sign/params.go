// Package sign renders printable Wi-Fi sign SVGs from user parameters.
package sign

import (
	"encoding/base64"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"qr-codes/internal/payload"
)

// Params are the user-controllable inputs for a sign.
type Params struct {
	// Network credentials, encoded into the QR payload only (never printed).
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
	Hidden   bool   `json:"hidden"`

	// Appearance
	AccentColor     string `json:"accentColor"`
	BackgroundColor string `json:"backgroundColor"`
	LogoDataURI     string `json:"logoDataUri"` // "" = no logo

	// Ornament is what fills the line between the QR card and the footer:
	// an ornament set name ("suits"), "tagline" (renders Tagline), or "none".
	Ornament string `json:"ornament"`
	Tagline  string `json:"tagline"`

	// Text; FooterText == "" hides the footer line.
	Headline   string `json:"headline"`
	Subtitle   string `json:"subtitle"`
	FooterText string `json:"footerText"`
}

// FieldError attributes a validation failure to a specific input field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

const (
	defaultAccentColor     = "#e8e2d4"
	defaultBackgroundColor = "#1b6f4c"
	defaultHeadline        = "WI-FI"
	defaultSubtitle        = "SCAN TO CONNECT"
)

const (
	maxLogoBytes     = 2 << 20 // 2 MiB decoded
	maxHeadlineRunes = 24
	maxSubtitleRunes = 40
	maxFooterRunes   = 120
	maxTaglineRunes  = 60
)

// ornamentSets are the available ornament designs; more sets will be added.
var ornamentSets = []string{"suits"}

// ornamentChoices are all valid Ornament values.
var ornamentChoices = append([]string{"tagline", "none"}, ornamentSets...)

var (
	hexColorRe    = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	logoDataURIRe = regexp.MustCompile(`^data:image/(png|jpeg|svg\+xml);base64,`)
)

// ApplyDefaults fills empty appearance and text fields with the template's
// original values. Credentials are never defaulted.
func (p *Params) ApplyDefaults() {
	if p.AccentColor == "" {
		p.AccentColor = defaultAccentColor
	}
	if p.BackgroundColor == "" {
		p.BackgroundColor = defaultBackgroundColor
	}
	if p.Headline == "" {
		p.Headline = defaultHeadline
	}
	if p.Subtitle == "" {
		p.Subtitle = defaultSubtitle
	}
	if p.Ornament == "" {
		p.Ornament = "suits"
	}
}

// Validate returns every problem with the params at once so a form can
// highlight all offending fields.
func (p *Params) Validate() []FieldError {
	var errs []FieldError
	fail := func(field, message string) {
		errs = append(errs, FieldError{Field: field, Message: message})
	}

	if p.SSID == "" {
		fail("ssid", "SSID is required")
	} else if len(p.SSID) > 32 {
		fail("ssid", "SSID must be at most 32 bytes")
	}

	auth, err := payload.ParseAuth(p.Auth)
	if err != nil {
		fail("auth", err.Error())
	} else if auth != payload.AuthNone {
		if p.Password == "" {
			fail("password", "password is required for "+string(auth))
		} else if len(p.Password) > 63 {
			fail("password", "password must be at most 63 characters")
		}
	}

	if !hexColorRe.MatchString(p.AccentColor) {
		fail("accentColor", "accent color must be a 6-digit hex color like #e8e2d4")
	}
	if !hexColorRe.MatchString(p.BackgroundColor) {
		fail("backgroundColor", "background color must be a 6-digit hex color like #1b6f4c")
	}

	if p.LogoDataURI != "" {
		if m := logoDataURIRe.FindStringSubmatch(p.LogoDataURI); m == nil {
			fail("logoDataUri", "logo must be a base64 data URI of type PNG, JPEG, or SVG")
		} else if decoded := base64.StdEncoding.DecodedLen(len(p.LogoDataURI) - len(m[0])); decoded > maxLogoBytes {
			fail("logoDataUri", "logo must be at most 2 MiB")
		}
	}

	switch {
	case p.Ornament == "tagline":
		if p.Tagline == "" {
			fail("tagline", "tagline text is required when the ornament is a tagline")
		} else if utf8.RuneCountInString(p.Tagline) > maxTaglineRunes {
			fail("tagline", "tagline must be at most 60 characters")
		}
	case !slices.Contains(ornamentChoices, p.Ornament):
		fail("ornament", "ornament must be one of: "+strings.Join(ornamentChoices, ", "))
	}

	if utf8.RuneCountInString(p.Headline) > maxHeadlineRunes {
		fail("headline", "headline must be at most 24 characters")
	}
	if utf8.RuneCountInString(p.Subtitle) > maxSubtitleRunes {
		fail("subtitle", "subtitle must be at most 40 characters")
	}
	if utf8.RuneCountInString(p.FooterText) > maxFooterRunes {
		fail("footerText", "footer text must be at most 120 characters")
	}

	return errs
}
