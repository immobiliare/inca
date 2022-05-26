package provider

import (
	"fmt"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Local struct {
	Provider
}

func (p Local) ID() string {
	return "local"
}

func (p Local) Get(commonName string) (*pki.CRT, error) {
	return nil, fmt.Errorf("write me")
}
func (p Local) For(commonName string) error {
	return fmt.Errorf("write me")
}
