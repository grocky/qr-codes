package sign

import (
	"strings"
	"testing"
)

func fieldsOf(errs []FieldError) []string {
	fields := make([]string, len(errs))
	for i, e := range errs {
		fields[i] = e.Field
	}
	return fields
}

func assertFields(t *testing.T, errs []FieldError, want ...string) {
	t.Helper()
	got := fieldsOf(errs)
	if len(got) != len(want) {
		t.Fatalf("Validate() errors on fields %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Validate() errors on fields %v, want %v", got, want)
		}
	}
}

func validParams() Params {
	p := Params{SSID: "home", Password: "s3cret", Auth: "WPA"}
	p.ApplyDefaults()
	return p
}

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Params)
		wantFields []string
	}{
		{"missing ssid", func(p *Params) { p.SSID = "" }, []string{"ssid"}},
		{"ssid over 32 bytes", func(p *Params) { p.SSID = strings.Repeat("x", 33) }, []string{"ssid"}},
		{"ssid at 32 bytes ok", func(p *Params) { p.SSID = strings.Repeat("x", 32) }, nil},
		{"missing password", func(p *Params) { p.Password = "" }, []string{"password"}},
		{"password over 63 chars", func(p *Params) { p.Password = strings.Repeat("x", 64) }, []string{"password"}},
		{"nopass needs no password", func(p *Params) { p.Auth = "nopass"; p.Password = "" }, nil},
		{"unknown auth", func(p *Params) { p.Auth = "open" }, []string{"auth"}},
		{"wpa2 normalizes", func(p *Params) { p.Auth = "wpa2" }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validParams()
			tt.mutate(&p)
			assertFields(t, p.Validate(), tt.wantFields...)
		})
	}
}

func TestValidateColors(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Params)
		wantFields []string
	}{
		{"accent not hex", func(p *Params) { p.AccentColor = "cream" }, []string{"accentColor"}},
		{"accent shorthand rejected", func(p *Params) { p.AccentColor = "#fff" }, []string{"accentColor"}},
		{"background not hex", func(p *Params) { p.BackgroundColor = "#12345g" }, []string{"backgroundColor"}},
		{"uppercase hex ok", func(p *Params) { p.AccentColor = "#E8E2D4" }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validParams()
			tt.mutate(&p)
			assertFields(t, p.Validate(), tt.wantFields...)
		})
	}
}

func TestValidateLogo(t *testing.T) {
	smallPNG := "data:image/png;base64," + strings.Repeat("A", 100)
	tests := []struct {
		name       string
		logo       string
		wantFields []string
	}{
		{"no logo ok", "", nil},
		{"png ok", smallPNG, nil},
		{"jpeg ok", "data:image/jpeg;base64,AAAA", nil},
		{"svg ok", "data:image/svg+xml;base64,AAAA", nil},
		{"gif rejected", "data:image/gif;base64,AAAA", []string{"logoDataUri"}},
		{"not a data uri", "https://example.com/logo.png", []string{"logoDataUri"}},
		{"non-base64 encoding rejected", "data:image/png,rawbytes", []string{"logoDataUri"}},
		{"over 2MiB rejected", "data:image/png;base64," + strings.Repeat("A", (2<<20)*4/3+8), []string{"logoDataUri"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validParams()
			p.LogoDataURI = tt.logo
			assertFields(t, p.Validate(), tt.wantFields...)
		})
	}
}

func TestValidateTextLimits(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Params)
		wantFields []string
	}{
		{"headline over 24 runes", func(p *Params) { p.Headline = strings.Repeat("é", 25) }, []string{"headline"}},
		{"headline at 24 runes ok", func(p *Params) { p.Headline = strings.Repeat("é", 24) }, nil},
		{"subtitle over 40 runes", func(p *Params) { p.Subtitle = strings.Repeat("x", 41) }, []string{"subtitle"}},
		{"footer over 120 runes", func(p *Params) { p.FooterText = strings.Repeat("x", 121) }, []string{"footerText"}},
		{"empty footer ok", func(p *Params) { p.FooterText = "" }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validParams()
			tt.mutate(&p)
			assertFields(t, p.Validate(), tt.wantFields...)
		})
	}
}

func TestValidateAggregatesAllErrors(t *testing.T) {
	p := Params{} // no ssid, no password, no auth
	p.ApplyDefaults()
	errs := p.Validate()
	if len(errs) < 2 {
		t.Fatalf("Validate() = %v, want multiple errors reported at once", errs)
	}
}

func TestValidateValidParams(t *testing.T) {
	p := Params{
		SSID:     "home",
		Password: "s3cret",
		Auth:     "WPA",
	}
	p.ApplyDefaults()

	if errs := p.Validate(); len(errs) != 0 {
		t.Fatalf("Validate() = %v, want no errors", errs)
	}
	if p.AccentColor != "#e8e2d4" {
		t.Errorf("default AccentColor = %q, want #e8e2d4", p.AccentColor)
	}
	if p.BackgroundColor != "#1b6f4c" {
		t.Errorf("default BackgroundColor = %q, want #1b6f4c", p.BackgroundColor)
	}
	if p.Headline != "WI-FI" {
		t.Errorf("default Headline = %q, want WI-FI", p.Headline)
	}
	if p.Subtitle != "SCAN TO CONNECT" {
		t.Errorf("default Subtitle = %q, want SCAN TO CONNECT", p.Subtitle)
	}
}
