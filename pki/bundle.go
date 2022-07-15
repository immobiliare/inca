package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"
	"strings"
)

func Wrap(crt *x509.Certificate, key *Key, ca *x509.Certificate, signer crypto.PrivateKey) (*pem.Block, *pem.Block, error) {
	crtBytes, err := x509.CreateCertificate(rand.Reader, crt, ca, key.Public(), signer)
	if err != nil {
		return nil, nil, err
	}
	crtBlock := pem.Block{Type: "CERTIFICATE", Bytes: crtBytes}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key.Value)
	if err != nil {
		return nil, nil, err
	}
	keyBlock := pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}

	return &crtBlock, &keyBlock, nil
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
