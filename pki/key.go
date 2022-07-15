package pki

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

type Key struct {
	Value any
	Algo  int
}

const (
	UnsupportedAlgorithm = iota
	EDDSA
	ECDSA
	RSA
)

func ParseKey(path string) (*Key, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	keyData, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key := Key{keyData, UnsupportedAlgorithm}

	switch keyData.(type) {
	case *rsa.PrivateKey:
		key.Algo = RSA
	case *ecdsa.PrivateKey:
		key.Algo = ECDSA
	case ed25519.PrivateKey:
		key.Algo = EDDSA
	default:
		key.Algo = UnsupportedAlgorithm
	}

	return &key, nil
}

func newKey(algo int) (*Key, error) {
	var (
		key any
		err error
	)
	switch algo {
	case EDDSA:
		_, key, err = ed25519.GenerateKey(rand.Reader)
	case ECDSA:
		key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case RSA:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
	default:
		return nil, errors.New("requested algorithm not supported")
	}

	if err != nil {
		return nil, err
	}
	return &Key{key, algo}, nil
}

func (key *Key) Public() any {
	switch k := key.Value.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}
