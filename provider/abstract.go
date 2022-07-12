package provider

import (
	"errors"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Provider interface {
	ID() string
	Tune(options ...string) error
	Get(commonName string) (*pki.CRT, error)
}

func Get(id string, options ...string) (*Provider, error) {
	for _, provider := range []Provider{
		new(Local),
	} {
		if id != provider.ID() {
			continue
		}
		if err := provider.Tune(options...); err != nil {
			return nil, err
		}
		return &provider, nil
	}

	return nil, errors.New("provider not found")
}
