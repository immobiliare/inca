package provider

import (
	"errors"
)

type Provider interface {
	ID() string
	Tune(options ...string) error
	For(name string) bool
	Get(name string, options map[string]string) ([]byte, []byte, error)
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

func GetFor(name string, options map[string]string, providers []*Provider) *Provider {
	for _, p := range providers {
		if (*p).For(name) {
			return p
		}
	}
	return nil
}
