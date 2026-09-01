package payload

import "testing"

func TestWiFiBuild(t *testing.T) {
	tests := []struct {
		name    string
		wifi    WiFi
		want    string
		wantErr bool
	}{
		{
			name: "wpa",
			wifi: WiFi{SSID: "home", Password: "s3cret", Auth: AuthWPA},
			want: "WIFI:T:WPA;S:home;P:s3cret;;",
		},
		{
			name: "wep",
			wifi: WiFi{SSID: "home", Password: "s3cret", Auth: AuthWEP},
			want: "WIFI:T:WEP;S:home;P:s3cret;;",
		},
		{
			name: "nopass omits password",
			wifi: WiFi{SSID: "cafe", Auth: AuthNone},
			want: "WIFI:T:nopass;S:cafe;;",
		},
		{
			name: "hidden",
			wifi: WiFi{SSID: "home", Password: "s3cret", Auth: AuthWPA, Hidden: true},
			want: "WIFI:T:WPA;S:home;P:s3cret;H:true;;",
		},
		{
			name: "special characters escaped",
			wifi: WiFi{SSID: `my; "net"`, Password: `p:a,s\s`, Auth: AuthWPA},
			want: `WIFI:T:WPA;S:my\; \"net\";P:p\:a\,s\\s;;`,
		},
		{
			name:    "missing ssid",
			wifi:    WiFi{Password: "s3cret", Auth: AuthWPA},
			wantErr: true,
		},
		{
			name:    "missing password with wpa",
			wifi:    WiFi{SSID: "home", Auth: AuthWPA},
			wantErr: true,
		},
		{
			name:    "invalid auth",
			wifi:    WiFi{SSID: "home", Password: "s3cret", Auth: "bogus"},
			wantErr: true,
		},
		{
			name:    "empty auth",
			wifi:    WiFi{SSID: "home", Password: "s3cret"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.wifi.Build()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Build() = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Build() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Build() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseAuth(t *testing.T) {
	tests := []struct {
		in      string
		want    Auth
		wantErr bool
	}{
		{in: "WPA", want: AuthWPA},
		{in: "wpa2", want: AuthWPA},
		{in: "WPA3", want: AuthWPA},
		{in: "wep", want: AuthWEP},
		{in: "nopass", want: AuthNone},
		{in: "none", want: AuthNone},
		{in: "", wantErr: true},
		{in: "open", wantErr: true},
	}
	for _, tt := range tests {
		got, err := ParseAuth(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseAuth(%q) = %q, want error", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseAuth(%q) error: %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseAuth(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
