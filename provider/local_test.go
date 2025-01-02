package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/immobiliare/inca/pki"
)

func TestLocal_ID(t *testing.T) {
	t.Parallel()
	p := Local{}
	if got := p.ID(); got != "Local" {
		t.Errorf("Local.ID() = %v, want %v", got, "Local")
	}
}

func TestLocal_Tune(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "cert-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}), 0600); err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		options map[string]interface{}
		wantErr bool
	}{
		{
			name:    "missing crt",
			options: map[string]interface{}{"key": keyPath},
			wantErr: true,
		},
		{
			name:    "missing key",
			options: map[string]interface{}{"crt": certPath},
			wantErr: true,
		},
		{
			name: "valid options",
			options: map[string]interface{}{
				"crt": certPath,
				"key": keyPath,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Local{}
			err := p.Tune(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.Tune() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_For(t *testing.T) {
	t.Parallel()

	p := &Local{
		crt: &x509.Certificate{
			DNSNames: []string{"example.com", "test.org"},
			Subject: pkix.Name{
				CommonName: "root.ca",
			},
		},
	}

	tests := []struct {
		name string
		dns  string
		want bool
	}{
		{
			name: "matching domain",
			dns:  "sub.example.com",
			want: true,
		},
		{
			name: "matching root domain",
			dns:  "example.com",
			want: true,
		},
		{
			name: "matching common name",
			dns:  "root.ca",
			want: true,
		},
		{
			name: "non-matching domain",
			dns:  "other.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := p.For(tt.dns); got != tt.want {
				t.Errorf("Local.For() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocal_Del(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "empty input",
			data:    []byte{},
			wantErr: false,
		},
		{
			name:    "with data",
			data:    []byte("test data"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Local{}
			if err := p.Del("test.com", tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Local.Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_CA(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "cert-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}), 0600); err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}), 0600); err != nil {
		t.Fatal(err)
	}

	p := &Local{}
	err = p.Tune(map[string]interface{}{
		"crt": certPath,
		"key": keyPath,
	})
	if err != nil {
		t.Fatal(err)
	}

	cert, err := p.CA()
	if err != nil {
		t.Errorf("Local.CA() error = %v", err)
		return
	}

	if len(cert) == 0 {
		t.Error("Local.CA() returned empty certificate")
	}

	block, _ := pem.Decode(cert)
	if block == nil {
		t.Error("Local.CA() returned invalid PEM data")
	}
}

func TestLocal_Config(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	p := &Local{
		crt: &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "test.example.com",
			},
			DNSNames:  []string{"test1.example.com", "test2.example.com"},
			NotBefore: testTime,
			NotAfter:  testTime.Add(24 * time.Hour),
		},
	}

	want := map[string]string{
		"Subject":           "test.example.com",
		"Alternative Names": "test1.example.com, test2.example.com",
		"Not Before":        "01/01/2023",
		"Not After":         "02/01/2023",
	}

	got := p.Config()

	for k, v := range want {
		if got[k] != v {
			t.Errorf("Local.Config()[%s] = %v, want %v", k, got[k], v)
		}
	}
}

func TestLocal_Get(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "cert-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Test CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	caCertBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}

	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatal(err)
	}

	p := &Local{
		crt: caCert,
		key: &pki.Key{Value: caKey, Algo: "rsa"},
	}

	tests := []struct {
		name    string
		domain  string
		options map[string]string
		wantErr bool
	}{
		{
			name:    "basic certificate",
			domain:  "test.example.com",
			options: map[string]string{},
			wantErr: false,
		},
		{
			name:   "with RSA algorithm",
			domain: "test.example.com",
			options: map[string]string{
				"algo": "rsa",
			},
			wantErr: false,
		},
		{
			name:   "with ECDSA algorithm",
			domain: "test.example.com",
			options: map[string]string{
				"algo": "ecdsa",
			},
			wantErr: false,
		},
		{
			name:   "with EDDSA algorithm",
			domain: "test.example.com",
			options: map[string]string{
				"algo": "eddsa",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certBytes, keyBytes, err := p.Get(tt.domain, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				certBlock, _ := pem.Decode(certBytes)
				if certBlock == nil {
					t.Error("Local.Get() returned invalid certificate PEM data")
					return
				}

				keyBlock, _ := pem.Decode(keyBytes)
				if keyBlock == nil {
					t.Error("Local.Get() returned invalid key PEM data")
					return
				}

				cert, err := x509.ParseCertificate(certBlock.Bytes)
				if err != nil {
					t.Errorf("Failed to parse generated certificate: %v", err)
					return
				}

				if cert.Subject.CommonName != tt.domain {
					t.Errorf("Certificate CommonName = %v, want %v", cert.Subject.CommonName, tt.domain)
				}

				if err := cert.CheckSignatureFrom(caCert); err != nil {
					t.Errorf("Certificate not properly signed by CA: %v", err)
				}
			}
		})
	}
}
