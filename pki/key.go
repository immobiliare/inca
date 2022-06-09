package pki

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

type Key struct {
	Value any
}

const (
	EDDSA int = iota
	RSA
)

func newKey(algo int) (*Key, error) {
	var (
		key any
		err error
	)
	switch algo {
	case EDDSA:
		_, key, err = ed25519.GenerateKey(rand.Reader)
	case RSA:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
	default:
		return nil, errors.New("requested algorithm not supported")
	}

	if err != nil {
		return nil, err
	}
	return &Key{key}, nil
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
