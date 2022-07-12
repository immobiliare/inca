package provider

import (
	"errors"
	"fmt"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Local struct {
	Provider
	crt *pki.CRT
	key *pki.Key
}

func (p Local) ID() string {
	return "local"
}

func (p *Local) Tune(options ...string) error {
	if len(options) != 2 {
		return fmt.Errorf("invalid number of options for provider %s: %s", p.ID(), options)
	}

	crt, err := pki.Parse(options[0])
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", options[0], err)
	}
	p.crt = crt

	key, err := pki.ParseKey(options[1])
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", options[1], err)
	}
	p.key = key

	return nil
}

func (p *Local) Get(commonName string) (*pki.CRT, error) {
	return nil, errors.New("write me")
}
