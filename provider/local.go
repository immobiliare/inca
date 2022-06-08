package provider

import (
	"errors"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Local struct {
	Provider
}

func (p Local) ID() string {
	return "local"
}

func (p Local) Get(commonName string) (*pki.CRT, error) {
	return nil, errors.New("write me")
}
func (p Local) For(commonName string) error {
	return errors.New("write me")
}
