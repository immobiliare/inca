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
