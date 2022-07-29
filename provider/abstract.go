package provider

import (
	"errors"
	"strings"
)

type Provider interface {
	ID() string
	Tune(options map[string]interface{}) error
	For(name string) bool
	Get(name string, options map[string]string) ([]byte, []byte, error)
	CA() ([]byte, error)
}

func Tune(id string, options map[string]interface{}) (*Provider, error) {
	for _, provider := range []Provider{
		new(Local),
		new(LetsEncrypt),
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

func GetByID(id string, providers []*Provider) (results []*Provider) {
	for _, p := range providers {
		if strings.EqualFold(id, (*p).ID()) {
			results = append(results, p)
		}
	}
	return
}

func GetByTargetName(name string, options map[string]string, providers []*Provider) *Provider {
	for _, p := range providers {
		if (*p).For(name) {
			return p
		}
	}
	return nil
}
