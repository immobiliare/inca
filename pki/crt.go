package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog/log"
)

type Request struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	CN            string
	DNSNames      []string
	IPAddresses   []string
	CA            bool
	Algo          string
	Duration      time.Duration
}

const DefaultCrtDuration = time.Duration(397 * 24 * time.Hour)

var DomainRegex = regexp.MustCompile(`^[*\.]{0,2}(?:(?:[\*a-z][a-z0-9-]+)\.)+[a-z]{2,}$`)

func Parse(path string) (*x509.Certificate, error) {
	data, err := os.ReadFile(path)
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

func ParseKeyPairBytes(crt, key []byte) (*x509.Certificate, *Key, error) {
	tls, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, nil, err
	}

	if crt, err := x509.ParseCertificate(tls.Certificate[0]); err != nil {
		return nil, nil, err
	} else {
		return crt, &Key{Value: tls.PrivateKey}, nil
	}
}

func ParseKeyPair(crtPath, keyPath string) (*x509.Certificate, *Key, error) {
	tls, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, nil, err
	}

	if crt, err := x509.ParseCertificate(tls.Certificate[0]); err != nil {
		return nil, nil, err
	} else {
		return crt, &Key{Value: tls.PrivateKey}, nil
	}
}

func NewRequest(options map[string]any) Request {
	req := Request{
		CA:          false,
		Algo:        DefaultCrtAlgo,
		Duration:    DefaultCrtDuration,
		DNSNames:    []string{},
		IPAddresses: []string{},
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
	if cn, ok := options["cn"]; ok {
		req.CN = cn.(string)
	}
	if alt, ok := options["alt"]; ok {
		req.DNSNames, req.IPAddresses = ParseAltNames(
			append([]string{req.CN}, strings.Split(alt.(string), ",")...))
	}
	if ca, ok := options["ca"]; ok {
		req.CA = ca.(bool)
	}
	if algo, ok := options["algo"]; ok {
		req.Algo = algo.(string)
	}
	if durationString, ok := options["duration"]; ok {
		if duration, err := util.ParseDuration(durationString.(string)); err == nil {
			req.Duration = duration
		} else {
			log.Error().Err(err).Msg("cannot parse duration")
		}
	}

	return req
}

func New(req Request) (*x509.Certificate, *Key, error) {
	var crt = x509.Certificate{}
	crt.Subject = pkix.Name{
		CommonName:    req.CN,
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
	crt.DNSNames = req.DNSNames
	for _, address := range req.IPAddresses {
		if ip := net.ParseIP(address); ip != nil {
			crt.IPAddresses = append(crt.IPAddresses, ip)
		}
	}

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

	if req.CA {
		crt.IsCA = true
		crt.KeyUsage |= x509.KeyUsageCertSign
	}

	key, err := newKey(req.Algo)
	return &crt, key, err
}

func IsValidCN(name string) bool {
	return DomainRegex.MatchString(name)
}

func AltNames(crt *x509.Certificate) ([]string, []string) {
	var (
		dnsNames    = []string{}
		ipAddresses = []string{}
	)
	for _, name := range append([]string{crt.Subject.CommonName}, crt.DNSNames...) {
		name = strings.TrimSpace(name)
		if len(name) > 0 {
			dnsNames = append(dnsNames, name)
		}
	}
	for _, ip := range crt.IPAddresses {
		ipString := strings.TrimSpace(ip.String())
		if len(ipString) > 0 {
			ipAddresses = append(ipAddresses, ip.String())
		}
	}
	dnsNames = util.StringSliceDistinct(dnsNames)
	ipAddresses = util.StringSliceDistinct(ipAddresses)
	sort.Strings(dnsNames)
	sort.Strings(ipAddresses)
	return dnsNames, ipAddresses
}

func ParseAltNames(altNames []string) (dnsNames, ipAddresses []string) {
	for _, name := range altNames {
		name = strings.TrimSpace(name)
		if ip := net.ParseIP(name); ip == nil && len(name) > 0 {
			dnsNames = append(dnsNames, name)
		} else if len(name) > 0 {
			ipAddresses = append(ipAddresses, name)
		}
	}
	dnsNames = util.StringSliceDistinct(dnsNames)
	ipAddresses = util.StringSliceDistinct(ipAddresses)
	sort.Strings(dnsNames)
	sort.Strings(ipAddresses)
	return
}
