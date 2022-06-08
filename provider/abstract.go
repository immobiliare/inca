package provider

import (
	"errors"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Provider interface {
	ID() string
	Get(commonName string) (*pki.CRT, error)
	For(commonName string) error
}

// All return the array of usable providers
func All() []Provider {
	return []Provider{
		new(Local),
	}
}

// For returns a provider for a given URL
func For(URL string) (Provider, error) {
	for _, p := range All() {
		if err := p.For(URL); err == nil {
			return p, nil
		}
	}

	return nil, errors.New("no provider found")
}
