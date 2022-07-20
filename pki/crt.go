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
	Algo          string
	Duration      time.Duration
}

const DefaultCrtDuration = time.Duration(100 * 365 * 24 * time.Hour)

func Parse(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseBytes(data)
}

func ParseBytes(data []byte) (*x509.Certificate, error) {
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

func NewRequest(options map[string]any) Request {
	req := Request{
		Hosts:    []string{},
		CA:       false,
		Algo:     DefaultCrtAlgo,
		Duration: DefaultCrtDuration,
	}

	if organization, ok := options["organization"]; ok {
		req.Organization = organization.(string)
	}
	if country, ok := options["country"]; ok {
		req.Country = country.(string)
	}
	if province, ok := options["province"]; ok {
		req.Province = province.(string)
	}
	if locality, ok := options["locality"]; ok {
		req.Locality = locality.(string)
	}
	if streetAddress, ok := options["street_address"]; ok {
		req.StreetAddress = streetAddress.(string)
	}
	if postalCode, ok := options["postal_code"]; ok {
		req.PostalCode = postalCode.(string)
	}
	if hosts, ok := options["hosts"]; ok {
		req.Hosts = hosts.([]string)
	}
	if ca, ok := options["ca"]; ok {
		req.CA = ca.(bool)
	}
	if algo, ok := options["algo"]; ok {
		req.Algo = algo.(string)
	}
	if duration, ok := options["duration"]; ok {
		req.Duration = duration.(time.Duration)
	}

	return req
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
