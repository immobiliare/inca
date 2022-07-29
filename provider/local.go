package provider

import (
	"crypto/x509"
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

func (p *Local) Tune(options map[string]interface{}) (err error) {
	crtPath, ok := options["crt"]
	if !ok {
		return fmt.Errorf("storage %s: crt not defined", p.ID())
	}

	keyPath, ok := options["key"]
	if !ok {
		return fmt.Errorf("storage %s: key not defined", p.ID())
	}

	p.crt, p.key, err = pki.ParseKeyPair(crtPath.(string), keyPath.(string))
	return
}

func (p *Local) For(name string) bool {
	for _, dns := range append(p.crt.DNSNames, p.crt.Subject.CommonName) {
		if strings.HasSuffix(name, dns) {
			return true
		}
	}
	return false
}

func (p *Local) Get(name string, options map[string]string) ([]byte, []byte, error) {
	reqOptions := make(map[string]any)
	reqOptions["cn"] = name
	for key, value := range options {
		reqOptions[key] = value
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

	if crt, key, err := pki.Wrap(crt, key, p.crt, p.key); err != nil {
		return nil, nil, err
	} else {
		return pki.ExportBytes(crt), pki.ExportBytes(key), nil
	}
}

func (p *Local) CA() ([]byte, error) {
	crt, err := pki.WrapCrt(p.crt, p.key, p.crt, p.key)
	if err != nil {
		return nil, err
	}

	return pki.ExportBytes(crt), nil
}
