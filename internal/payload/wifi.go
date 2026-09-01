package payload

import (
	"errors"
	"fmt"
	"strings"
)

// Auth is the authentication type of a WiFi network.
type Auth string

const (
	AuthWPA  Auth = "WPA"
	AuthWEP  Auth = "WEP"
	AuthNone Auth = "nopass"
)

// ParseAuth normalizes a user-supplied auth type. WPA2/WPA3 map to WPA, which
// is the only value scanners recognize for the WPA family.
func ParseAuth(s string) (Auth, error) {
	switch strings.ToUpper(s) {
	case "WPA", "WPA2", "WPA3":
		return AuthWPA, nil
	case "WEP":
		return AuthWEP, nil
	case "NOPASS", "NONE":
		return AuthNone, nil
	}
	return "", fmt.Errorf("unknown auth type %q (want WPA, WEP, or nopass)", s)
}

// WiFi builds a WIFI: payload per the de-facto ZXing convention:
//
//	WIFI:T:WPA;S:mynetwork;P:mypass;H:true;;
type WiFi struct {
	SSID     string
	Password string
	Auth     Auth
	Hidden   bool
}

func (w WiFi) Build() (string, error) {
	if w.SSID == "" {
		return "", errors.New("wifi: ssid is required")
	}
	switch w.Auth {
	case AuthWPA, AuthWEP:
		if w.Password == "" {
			return "", fmt.Errorf("wifi: password is required for auth type %s", w.Auth)
		}
	case AuthNone:
		// no password expected
	default:
		return "", fmt.Errorf("wifi: invalid auth type %q", w.Auth)
	}

	var b strings.Builder
	b.WriteString("WIFI:")
	fmt.Fprintf(&b, "T:%s;", w.Auth)
	fmt.Fprintf(&b, "S:%s;", escapeWiFi(w.SSID))
	if w.Auth != AuthNone {
		fmt.Fprintf(&b, "P:%s;", escapeWiFi(w.Password))
	}
	if w.Hidden {
		b.WriteString("H:true;")
	}
	b.WriteString(";")
	return b.String(), nil
}

// The convention backslash-escapes these characters in SSID and password.
var wifiEscaper = strings.NewReplacer(
	`\`, `\\`,
	`;`, `\;`,
	`,`, `\,`,
	`:`, `\:`,
	`"`, `\"`,
)

func escapeWiFi(s string) string {
	return wifiEscaper.Replace(s)
}
