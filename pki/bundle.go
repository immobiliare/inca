package pki

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"
	"strings"
)

func WrapCrt(crt *x509.Certificate, key *Key, ca *x509.Certificate, signer *Key) (*pem.Block, error) {
	bytes, err := x509.CreateCertificate(rand.Reader, crt, ca, key.Public(), signer.Value)
	if err != nil {
		return nil, err
	}

	return &pem.Block{Type: "CERTIFICATE", Bytes: bytes}, nil
}

func WrapKey(key *Key) (*pem.Block, error) {
	bytes, err := x509.MarshalPKCS8PrivateKey(key.Value)
	if err != nil {
		return nil, err
	}

	return &pem.Block{Type: "PRIVATE KEY", Bytes: bytes}, nil
}

func Wrap(crt *x509.Certificate, key *Key, ca *x509.Certificate, signer *Key) (*pem.Block, *pem.Block, error) {
	crtBlock, err := WrapCrt(crt, key, ca, signer)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, err := WrapKey(key)
	if err != nil {
		return nil, nil, err
	}

	return crtBlock, keyBlock, nil
}

func Export(block *pem.Block, path string) (err error) {
	var (
		perms = fs.FileMode(0644)
	)
	if strings.Contains(block.Type, "PRIVATE KEY") {
		perms = fs.FileMode(0600)
	}

	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return err
	}
	defer output.Close()

	if err := pem.Encode(output, block); err != nil {
		return err
	}

	return
}
