package payload

import "testing"

func TestURLBuild(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "https", in: "https://example.com/page?a=1"},
		{name: "http", in: "http://example.com"},
		{name: "empty", in: "", wantErr: true},
		{name: "no scheme", in: "example.com/page", wantErr: true},
		{name: "relative path", in: "/just/a/path", wantErr: true},
		{name: "garbage", in: "not a url", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URL{URL: tt.in}.Build()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Build() = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Build() error: %v", err)
			}
			if got != tt.in {
				t.Errorf("Build() = %q, want %q", got, tt.in)
			}
		})
	}
}
