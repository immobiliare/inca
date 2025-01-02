package provider

import (
	"reflect"
	"testing"
)

func TestLetsEncrypt_ID(t *testing.T) {
	t.Parallel()
	p := LetsEncrypt{}
	if got := p.ID(); got != "LetsEncrypt" {
		t.Errorf("ID() = %v, want %v", got, "LetsEncrypt")
	}
}

func TestLetsEncrypt_For(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		targets  []*LetsEncryptTarget
		domain   string
		expected bool
	}{
		{
			name: "matching domain",
			targets: []*LetsEncryptTarget{
				{domain: "example.com"},
			},
			domain:   "test.example.com",
			expected: true,
		},
		{
			name: "non-matching domain",
			targets: []*LetsEncryptTarget{
				{domain: "example.com"},
			},
			domain:   "test.different.com",
			expected: false,
		},
		{
			name:     "empty targets",
			targets:  []*LetsEncryptTarget{},
			domain:   "test.example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &LetsEncrypt{
				targets: tt.targets,
			}
			if got := p.For(tt.domain); got != tt.expected {
				t.Errorf("For() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLetsEncrypt_Config(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		ca      string
		email   string
		targets []*LetsEncryptTarget
		want    map[string]string
	}{
		{
			name:  "basic config",
			ca:    "https://example.com/ca.pem",
			email: "test@example.com",
			targets: []*LetsEncryptTarget{
				{domain: "example.com", provider: "dns01"},
				{domain: "test.com", provider: "http01"},
			},
			want: map[string]string{
				"CA":      "https://example.com/ca.pem",
				"Email":   "test@example.com",
				"Targets": "example.com (dns01), test.com (http01)",
			},
		},
		{
			name:    "empty targets",
			ca:      "https://example.com/ca.pem",
			email:   "test@example.com",
			targets: []*LetsEncryptTarget{},
			want: map[string]string{
				"CA":      "https://example.com/ca.pem",
				"Email":   "test@example.com",
				"Targets": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &LetsEncrypt{
				ca:      tt.ca,
				user:    &LetsEncryptUser{email: tt.email},
				targets: tt.targets,
			}
			if got := p.Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config() = %v, want %v", got, tt.want)
			}
		})
	}
}
