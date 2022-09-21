package server

import (
	"crypto/x509"
	"strings"

	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/storage"
)

const (
	foreignProviderID = "Foreign"
)

type JSONProvider struct {
	ID     string            `json:"id"`
	Config map[string]string `json:"config"`
}

type JSONStorage JSONProvider

type JSONCrt struct {
	CN            string       `json:"name"`
	AltNames      []string     `json:"alt"`
	NotBefore     string       `json:"not_before"`
	NotAfter      string       `json:"not_after"`
	Organization  string       `json:"organization"`
	Country       string       `json:"country"`
	Province      string       `json:"province"`
	Locality      string       `json:"locality"`
	StreetAddress string       `json:"street_address"`
	PostalCode    string       `json:"postal_code"`
	Provider      JSONProvider `json:"provider"`
}

func EncodeCrt(crt *x509.Certificate, provider *provider.Provider) JSONCrt {
	dnsNames, ipAddresses := pki.AltNames(crt)
	return JSONCrt{
		crt.Subject.CommonName,
		append(dnsNames, ipAddresses...),
		crt.NotBefore.Format("02/01/2006"),
		crt.NotAfter.Format("02/01/2006"),
		strings.Join(crt.Subject.Organization, ","),
		strings.Join(crt.Subject.Country, ","),
		strings.Join(crt.Subject.Province, ","),
		strings.Join(crt.Subject.Locality, ","),
		strings.Join(crt.Subject.StreetAddress, ","),
		strings.Join(crt.Subject.PostalCode, ","),
		EncodeProvider(provider),
	}
}

func EncodeProvider(provider *provider.Provider) JSONProvider {
	if provider == nil {
		return JSONProvider{foreignProviderID, map[string]string{}}
	}
	return JSONProvider{(*provider).ID(), (*provider).Config()}
}

func EncodeProviders(providers []*provider.Provider) (encoded []JSONProvider) {
	for _, provider := range providers {
		encoded = append(encoded, EncodeProvider(provider))
	}
	return
}

func EncodeStorage(storage *storage.Storage) JSONStorage {
	return JSONStorage{(*storage).ID(), (*storage).Config()}
}
