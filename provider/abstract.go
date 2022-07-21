package provider

import (
	"encoding/pem"
	"errors"
	"strings"
)

type Provider interface {
	ID() string
	Tune(options map[string]interface{}) error
	For(name string) bool
	Get(name string, options map[string]string) (*pem.Block, *pem.Block, error)
	CA() (*pem.Block, error)
}

func Get(id string, options map[string]interface{}) (*Provider, error) {
	for _, provider := range []Provider{
		new(Local),
	} {
		if !strings.EqualFold(id, provider.ID()) {
			continue
		}

		if err := provider.Tune(options); err != nil {
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

func GetFrom(id string, providers []*Provider) *Provider {
	for _, p := range providers {
		if (*p).ID() == id {
			return p
		}
	}
	return nil
}
