package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"
)

func Wrap(crt *x509.Certificate, key *Key, ca *x509.Certificate, signer crypto.PrivateKey) ([]byte, []byte, error) {
	crtBytes, err := x509.CreateCertificate(rand.Reader, crt, ca, key.Public(), signer)
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key.Value)
	if err != nil {
		return nil, nil, err
	}

	return crtBytes, keyBytes, nil
}

func Export(payload []byte, path string, isKey bool) (err error) {
	var (
		perms  = fs.FileMode(0644)
		header = "CERTIFICATE"
	)
	if isKey {
		perms = fs.FileMode(0600)
		header = "PRIVATE KEY"
	}

	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return err
	}
	defer output.Close()

	if err := pem.Encode(output, &pem.Block{Type: header, Bytes: payload}); err != nil {
		return err
	}

	return
}
