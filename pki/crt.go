package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

type Request struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	Hosts         []string
	CA            bool
	Algo          int
	Duration      time.Duration
}

func Parse(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return crt, nil
}

func ParseKeyPair(crtPath, keyPath string) (*x509.Certificate, *Key, error) {
	tls, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, nil, err
	}

	crt, err := x509.ParseCertificate(tls.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	return crt, &Key{Value: tls.PrivateKey}, nil
}

func NewRequest(names ...string) Request {
	return Request{
		Organization:  "Immobiliare.it",
		Country:       "IT",
		Province:      "RM",
		Locality:      "Rome",
		StreetAddress: "Via di Santa Prassede",
		PostalCode:    "00184",
		Hosts:         names,
		CA:            false,
		Algo:          ECDSA,
		Duration:      time.Duration(100 * 365 * 24 * time.Hour),
	}
}

func New(req Request) (*x509.Certificate, *Key, error) {
	var crt = x509.Certificate{}
	crt.Subject = pkix.Name{
		CommonName:    req.Hosts[0],
		Organization:  []string{req.Organization},
		Country:       []string{req.Country},
		Province:      []string{req.Province},
		Locality:      []string{req.Locality},
		StreetAddress: []string{req.StreetAddress},
		PostalCode:    []string{req.PostalCode},
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

func IsValidCN(name string) bool {
	return len(name) > 3
}
