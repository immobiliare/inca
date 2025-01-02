package provider

import (
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
