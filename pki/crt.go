package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type CRT struct {
	x509.Certificate
}

type Request struct {
	Organization string
	Hosts        []string
	CA           bool
	Algo         int
	Duration     time.Duration
}

func Parse(path string) (*CRT, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &CRT{*crt}, nil
}

func NewRequest() Request {
	return Request{
		Organization: "",
		Hosts:        []string{},
		CA:           true,
		Algo:         EDDSA,
		Duration:     time.Duration(100 * 365 * 24 * time.Hour),
	}
}

func New(req Request) (*CRT, *Key, error) {
	var crt = CRT{}
	crt.Subject = pkix.Name{
		Organization: []string{req.Organization},
		// Country:       []string{"US"},
		// Province:      []string{""},
		// Locality:      []string{"San Francisco"},
		// StreetAddress: []string{"Golden Gate Bridge"},
		// PostalCode:    []string{"94016"},
	}
	crt.BasicConstraintsValid = true
	crt.NotBefore = time.Now()
	crt.NotAfter = crt.NotBefore.Add(req.Duration)

	crt.KeyUsage = x509.KeyUsageDigitalSignature
	crt.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	if req.Algo == RSA {
		crt.KeyUsage |= x509.KeyUsageKeyEncipherment
	}

	if serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)); err != nil {
		return nil, nil, err
	} else {
		crt.SerialNumber = serialNumber
	}

	for _, host := range req.Hosts {
		if ip := net.ParseIP(host); ip != nil {
			crt.IPAddresses = append(crt.IPAddresses, ip)
		} else {
			crt.DNSNames = append(crt.DNSNames, host)
		}
	}

	if req.CA {
		crt.IsCA = true
		crt.KeyUsage |= x509.KeyUsageCertSign
	}

	key, err := newKey(req.Algo)
	return &crt, key, err
}

func Export(crt *CRT, key *Key, path string) error {
	derBytes, err := x509.CreateCertificate(rand.Reader, &crt.Certificate, &crt.Certificate, key.Public(), key.Value)
	if err != nil {
		return err
	}

	certOut, err := os.Create(filepath.Join(path, "crt.pem"))
	if err != nil {
		return err
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	if err := certOut.Close(); err != nil {
		return err
	}

	keyOut, err := os.OpenFile(filepath.Join(path, "key.pem"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(key.Value)
	if err != nil {
		return err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	if err := keyOut.Close(); err != nil {
		return err
	}

	return nil
}
