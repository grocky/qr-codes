package sign

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type renderResult struct {
	SVG       string       `json:"svg"`
	QRModules int          `json:"qrModules"`
	Errors    []FieldError `json:"errors"`
}

func render(t *testing.T, params string) renderResult {
	t.Helper()
	var res renderResult
	if err := json.Unmarshal(Render([]byte(params)), &res); err != nil {
		t.Fatalf("Render() returned invalid JSON: %v", err)
	}
	return res
}

func TestRenderValidParams(t *testing.T) {
	res := render(t, `{"ssid":"home","password":"s3cret","auth":"WPA"}`)
	if len(res.Errors) != 0 {
		t.Fatalf("Render() errors = %v, want none", res.Errors)
	}
	if res.QRModules < 21 {
		t.Errorf("qrModules = %d, want >= 21", res.QRModules)
	}
	for _, want := range []string{
		"<svg", "</svg>",
		">WI-FI</text>",
		">SCAN TO CONNECT</text>",
		`shape-rendering="crispEdges"`, // vector QR group
		"#e8e2d4",                      // default accent
	} {
		if !strings.Contains(res.SVG, want) {
			t.Errorf("SVG missing %q", want)
		}
	}
	if strings.Contains(res.SVG, "s3cret") {
		t.Error("SVG must not contain the password as text")
	}
}

func TestRenderInvalidParams(t *testing.T) {
	res := render(t, `{"ssid":"","auth":"WPA"}`)
	if res.SVG != "" {
		t.Fatal("Render() produced SVG for invalid params")
	}
	if len(res.Errors) < 2 { // ssid + password
		t.Fatalf("Render() errors = %v, want ssid and password errors", res.Errors)
	}
}

func TestRenderEscapesUserText(t *testing.T) {
	res := render(t, `{"ssid":"x\"><script>alert(1)</script>","password":"s3cret","auth":"WPA","headline":"<b>&\"","footerText":"a < b & c"}`)
	if len(res.Errors) != 0 {
		t.Fatalf("Render() errors = %v", res.Errors)
	}
	for _, forbidden := range []string{"<script>", "<b>"} {
		if strings.Contains(res.SVG, forbidden) {
			t.Errorf("SVG contains unescaped %q", forbidden)
		}
	}
	if !strings.Contains(res.SVG, "&lt;b&gt;&amp;") {
		t.Error("headline was not XML-escaped into the SVG")
	}
	assertWellFormedXML(t, res.SVG)
}

func assertWellFormedXML(t *testing.T, svg string) {
	t.Helper()
	dec := xml.NewDecoder(strings.NewReader(svg))
	for {
		if _, err := dec.Token(); err != nil {
			if err == io.EOF {
				return
			}
			t.Fatalf("SVG is not well-formed XML: %v", err)
		}
	}
}

func TestRenderCustomization(t *testing.T) {
	base := `"ssid":"home","password":"s3cret","auth":"WPA"`
	tests := []struct {
		name        string
		params      string
		contains    []string
		notContains []string
	}{
		{
			name:        "no logo omits image element",
			params:      `{` + base + `}`,
			notContains: []string{"<image"},
		},
		{
			name:     "logo embedded",
			params:   `{` + base + `,"logoDataUri":"data:image/png;base64,AAAA"}`,
			contains: []string{`<image id="logo"`, "data:image/png;base64,AAAA"},
		},
		{
			name:        "suits hidden",
			params:      `{` + base + `,"showSuits":false,"footerText":"hi"}`,
			notContains: []string{"suit-row"},
		},
		{
			name:     "suits shown",
			params:   `{` + base + `,"showSuits":true}`,
			contains: []string{`id="suit-row"`, `xlink:href="#suit-row"`},
		},
		{
			name:        "empty footer hides text line",
			params:      `{` + base + `}`,
			notContains: []string{`y="1046"`},
		},
		{
			name:        "custom colors flow to accent and derived felt",
			params:      `{` + base + `,"accentColor":"#ff0000","backgroundColor":"#204060"}`,
			contains:    []string{`stroke="#ff0000"`, `stop-color="#204060"`, `stop-color="#102030"`},
			notContains: []string{"#e8e2d4", "#1b6f4c"},
		},
		{
			name:        "long headline shrinks font size",
			params:      `{` + base + `,"headline":"VERY LONG HEADLINE OK"}`,
			notContains: []string{`font-size="94"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := render(t, tt.params)
			if len(res.Errors) != 0 {
				t.Fatalf("Render() errors = %v", res.Errors)
			}
			for _, s := range tt.contains {
				if !strings.Contains(res.SVG, s) {
					t.Errorf("SVG missing %q", s)
				}
			}
			for _, s := range tt.notContains {
				if strings.Contains(res.SVG, s) {
					t.Errorf("SVG unexpectedly contains %q", s)
				}
			}
			assertWellFormedXML(t, res.SVG)
		})
	}
}

var update = flag.Bool("update", false, "rewrite golden files")

func TestRenderGolden(t *testing.T) {
	base := `"ssid":"home","password":"s3cret","auth":"WPA"`
	tests := []struct {
		name   string
		params string
	}{
		{"defaults", `{` + base + `,"showSuits":true,"footerText":"Point your camera at the code — no typing required."}`},
		{"customized", `{` + base + `,"accentColor":"#d4af37","backgroundColor":"#1a1a2e","logoDataUri":"data:image/png;base64,iVBORw0KGgo=","showSuits":false,"headline":"GUEST WI-FI","subtitle":"SCAN TO JOIN","footerText":"Welcome!"}`},
		{"long-text-shrink", `{` + base + `,"headline":"THE EXTREMELY LONG NAME","footerText":"` + strings.Repeat("word ", 23) + `end"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := render(t, tt.params)
			if len(res.Errors) != 0 {
				t.Fatalf("Render() errors = %v", res.Errors)
			}
			assertWellFormedXML(t, res.SVG)

			path := filepath.Join("testdata", "golden", tt.name+".svg")
			if *update {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(res.SVG), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading golden (run with -update to create): %v", err)
			}
			if res.SVG != string(want) {
				t.Errorf("SVG differs from golden %s (run with -update after intentional changes)", path)
			}
		})
	}
}

func TestRenderMalformedJSON(t *testing.T) {
	res := render(t, `{not json`)
	if len(res.Errors) == 0 {
		t.Fatal("Render() must report an error for malformed JSON")
	}
}
