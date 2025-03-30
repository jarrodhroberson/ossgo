package sendgrid

import (
	"testing"
)

func TestIsValidDomain(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid IP Address",
			args:    args{domain: "127.0.0.1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid Hostname",
			args:    args{domain: "example.com"},
			want:    true, //  This test relies on your network connection to resolve "example.com"
			wantErr: false,
		},
		{
			name:    "Hostname with trailing dot",
			args:    args{domain: "example.com."},
			want:    true, //  This test also relies on network resolution
			wantErr: false,
		},
		{
			name:    "Invalid IP Address Format",
			args:    args{domain: "127.0.0.256"},
			want:    false,
			wantErr: true, //  net.ParseIP will fail on invalid IP.
		},
		{
			name:    "Invalid Hostname",
			args:    args{domain: "invalid.domain.com"},
			want:    false,
			wantErr: true, //  Will fail if the domain can't be resolved.
		},
		{
			name:    "Empty Domain",
			args:    args{domain: ""},
			want:    false,
			wantErr: true, // net.LookupHost will likely fail, but potentially not always.
		},
		{
			name:    "IP Address with leading/trailing spaces",
			args:    args{domain: "  127.0.0.1  "}, // should still validate (less reliable).
			want:    true,                          // It may still validate depending on your network.
			wantErr: false,
		},
		{
			name:    "Valid IPv6 Address",
			args:    args{domain: "2001:db8::1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid IPv6 with brackets",
			args:    args{domain: "[2001:db8::1]"}, // net.ParseIP doesn't handle brackets for IP address directly.
			want:    false,                         // this depends on your needs
			wantErr: true,
		},
		{
			name:    "Invalid IPv6 Format",
			args:    args{domain: "2001:db8:x:x:x:x:x:x"}, // invalid hex values.
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isValidDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isValidDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmailAddressFormatValid(t *testing.T) {
	type args struct {
		emailAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid email",
			args:    args{emailAddress: "test@example.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid email with subdomain",
			args:    args{emailAddress: "user.name@sub.example.co.uk"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid email with plus sign",
			args:    args{emailAddress: "user+alias@example.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid email with numbers",
			args:    args{emailAddress: "user123@example.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Email with leading/trailing spaces",
			args:    args{emailAddress: "  test@example.com  "},
			want:    true, // This may depend on your `isValidRFC822Email` implementation
			wantErr: false,
		},
		{
			name:    "Missing @ symbol",
			args:    args{emailAddress: "testexample.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Missing domain",
			args:    args{emailAddress: "test@"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Missing local part",
			args:    args{emailAddress: "@example.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid character in local part",
			args:    args{emailAddress: "test!@example.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid character in domain",
			args:    args{emailAddress: "test@exam!ple.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid TLD",
			args:    args{emailAddress: "test@example.c"}, // TLD less than 2 chars (should be caught by regex)
			want:    false,
			wantErr: true,
		},
		{
			name:    "Empty string",
			args:    args{emailAddress: ""},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Long email address (over 254 chars)",
			args:    args{emailAddress: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Local part too long (over 64)",
			args:    args{emailAddress: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Domain part too long (over 255)",
			args:    args{emailAddress: "test@" + string(make([]rune, 256)) + ".com"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsEmailAddressFormatValid(tt.args.emailAddress) // Corrected function call
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidRFC822Email() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isValidRFC822Email() got = %v, want %v", got, tt.want)
			}
		})
	}
}
