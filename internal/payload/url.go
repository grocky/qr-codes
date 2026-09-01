package payload

import (
	"errors"
	"fmt"
	neturl "net/url"
)

// URL builds a plain URL payload, validating that the value is an absolute URL.
type URL struct {
	URL string
}

func (u URL) Build() (string, error) {
	if u.URL == "" {
		return "", errors.New("url: value is required")
	}
	parsed, err := neturl.ParseRequestURI(u.URL)
	if err != nil {
		return "", fmt.Errorf("url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("url: %q must be absolute (include scheme and host)", u.URL)
	}
	return u.URL, nil
}
