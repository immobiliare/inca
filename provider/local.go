package provider

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Local struct {
	Provider
	crt *x509.Certificate
	tls *tls.Certificate
}

func (p Local) ID() string {
	return "local"
}

func (p *Local) Tune(options ...string) (err error) {
	if len(options) != 2 {
		return fmt.Errorf("invalid number of options for provider %s: %s", p.ID(), options)
	}

	p.crt, p.tls, err = pki.ParseKeyPair(options[0], options[1])
	return
}

func (p *Local) For(name string) bool {
	for _, dns := range p.crt.DNSNames {
		if strings.HasSuffix(name, dns) {
			return true
		}
	}
	return false
}

func (p *Local) Get(name string, options map[string]string) (*pem.Block, *pem.Block, error) {
	crt, key, err := pki.New(pki.NewRequest(name))
	if err != nil {
		return nil, nil, err
	}

	return pki.Wrap(crt, key, p.crt, p.tls.PrivateKey)
}
