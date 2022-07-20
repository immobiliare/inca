package provider

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"gitlab.rete.farm/sistemi/inca/pki"
)

type Local struct {
	Provider
	crt *x509.Certificate
	key *pki.Key
}

func (p Local) ID() string {
	return "local"
}

func (p *Local) Tune(options ...string) (err error) {
	if len(options) != 2 {
		return fmt.Errorf("invalid number of options for provider %s: %s", p.ID(), options)
	}

	p.crt, p.key, err = pki.ParseKeyPair(options[0], options[1])
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
	reqOptions := make(map[string]any)
	for key, value := range options {
		reqOptions[key] = value
	}

	reqOptions["hosts"] = []string{name}
	if altNames, ok := options["alt"]; ok {
		reqOptions["hosts"] = append(reqOptions["hosts"].([]string), strings.Split(altNames, ",")...)
	}

	if algo, ok := options["algo"]; ok {
		switch algo {
		case "eddsa":
			reqOptions["algo"] = pki.EDDSA
		case "ecdsa":
			reqOptions["algo"] = pki.ECDSA
		case "rsa":
			reqOptions["algo"] = pki.RSA
		}
	}

	req := pki.NewRequest(reqOptions)
	crt, key, err := pki.New(req)
	if err != nil {
		return nil, nil, err
	}

	return pki.Wrap(crt, key, p.crt, p.key)
}

func (p *Local) CA() (*pem.Block, error) {
	return pki.WrapCrt(p.crt, p.key, p.crt, p.key)
}
