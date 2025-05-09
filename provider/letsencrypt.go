package provider

import (
	"crypto"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/providers/http/webroot"
	"github.com/go-acme/lego/v4/registration"
	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/util"
)

type LetsEncrypt struct {
	Provider
	user    *LetsEncryptUser
	client  *lego.Client
	ca      string
	targets []*LetsEncryptTarget
}

type LetsEncryptTarget struct {
	domain      string
	provider    string
	environment map[string]string
}

type LetsEncryptUser struct {
	email        string
	key          crypto.PrivateKey
	registration *registration.Resource
}

func init() {
	log.Logger = util.NewZStdLogger()
}

func (u LetsEncryptUser) GetEmail() string {
	return u.email
}
func (u LetsEncryptUser) GetRegistration() *registration.Resource {
	return u.registration
}
func (u LetsEncryptUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (p LetsEncrypt) ID() string {
	return "LetsEncrypt"
}

func (p *LetsEncrypt) Tune(options map[string]interface{}) (err error) {
	keyPath, ok := options["key"]
	if !ok {
		return fmt.Errorf("provider %s: key not defined", p.ID())
	}

	email, ok := options["email"]
	if !ok {
		return fmt.Errorf("provider %s: email not defined", p.ID())
	}

	configTargets, ok := options["targets"]
	if !ok {
		configTargets = make([]interface{}, 0)
	}
	targets := make([]*LetsEncryptTarget, 0, len(configTargets.([]interface{})))
	for _, configTarget := range configTargets.([]interface{}) {
		var (
			target          = LetsEncryptTarget{environment: make(map[string]string)}
			configTargetMap = configTarget.(map[string]interface{})
		)
		domain, ok := configTargetMap["domain"]
		if !ok {
			return fmt.Errorf("provider %s: target domain not defined", p.ID())
		}
		target.domain = domain.(string)

		challenge, ok := configTargetMap["challenge"]
		if !ok {
			return fmt.Errorf("provider %s: target %s challenge not defined", p.ID(), target.domain)
		}
		challengeMap := challenge.(map[string]interface{})

		challengeProvider, ok := challengeMap["id"]
		if !ok {
			return fmt.Errorf("provider %s: target %s challenge type not defined", p.ID(), target.domain)
		}
		target.provider = challengeProvider.(string)

		for key, value := range challenge.(map[string]interface{}) {
			if strings.EqualFold(key, "id") {
				continue
			}
			target.environment[strings.ToUpper(key)] = value.(string)
		}

		targets = append(targets, &target)
	}

	key, err := pki.ParseKey(keyPath.(string))
	if err != nil {
		return err
	}

	user := &LetsEncryptUser{email.(string), key.Value, nil}
	config := lego.NewConfig(user)
	switch pki.DefaultCrtAlgo {
	case pki.ECDSA:
		config.Certificate.KeyType = certcrypto.EC256
		if caURL, ok := options["ca"]; ok {
			p.ca = caURL.(string)
		} else {
			p.ca = "https://letsencrypt.org/certs/isrg-root-x2.pem"
		}
	case pki.RSA:
		config.Certificate.KeyType = certcrypto.RSA2048
		if caURL, ok := options["ca"]; ok {
			p.ca = caURL.(string)
		} else {
			p.ca = "https://letsencrypt.org/certs/isrgrootx1.pem"
		}
	default:
		return errors.New("unsupported key algorithm")
	}

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	registration, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}

	user.registration = registration
	p.user = user
	p.client = client
	p.targets = targets
	return
}

func (p *LetsEncrypt) For(name string) bool {
	_, err := p.getChallengeProvider(name)
	return err == nil
}

func (p *LetsEncrypt) getChallengeProvider(domain string) (*LetsEncryptTarget, error) {
	for _, target := range p.targets {
		if strings.HasSuffix(domain, target.domain) {
			return target, nil
		}
	}
	return nil, errors.New("challenge provider not found")
}

func (p *LetsEncrypt) Get(name string, options map[string]string) ([]byte, []byte, error) {
	targetProvider, err := p.getChallengeProvider(name)
	if err != nil {
		return nil, nil, err
	}

	for envKey, envValue := range targetProvider.environment {
		if err := os.Setenv(envKey, envValue); err != nil {
			return nil, nil, err
		}
	}

	if err := p.SetChallengeProvider(targetProvider.provider); err != nil {
		return nil, nil, err
	}

	names := []string{name}
	if alt, ok := options["alt"]; ok {
		alt = strings.ReplaceAll(alt, " ", "")
		names = append(names, strings.Split(alt, ",")...)
	}

	request := certificate.ObtainRequest{
		Domains: names,
		Bundle:  true,
	}

	certificates, err := p.client.Certificate.Obtain(request)
	if err != nil {
		return nil, nil, err
	}

	return certificates.Certificate, certificates.PrivateKey, nil
}

func (p *LetsEncrypt) SetChallengeProvider(providerId string) error {
	if providerId == "webroot" {
		provider, err := webroot.NewHTTPProvider("/tmp/server/webroot")
		if err != nil {
			return err
		}

		return p.client.Challenge.SetHTTP01Provider(provider)
	}

	provider, err := dns.NewDNSChallengeProviderByName(providerId)
	if err != nil {
		return err
	}

	return p.client.Challenge.SetDNS01Provider(provider)
}

func (p *LetsEncrypt) Del(name string, data []byte) error {
	targetProvider, err := p.getChallengeProvider(name)
	if err != nil {
		return err
	}

	for envKey, envValue := range targetProvider.environment {
		if err := os.Setenv(envKey, envValue); err != nil {
			return err
		}
	}

	return p.client.Certificate.Revoke(data)
}

func (p *LetsEncrypt) CA() (b []byte, e error) {
	response, err := http.Get(p.ca)
	if err != nil {
		return nil, err
	}
	b, e = io.ReadAll(response.Body)
	e = errors.Join(e, response.Body.Close())
	return
}

func (p *LetsEncrypt) Config() map[string]string {
	targets := []string{}
	for _, target := range p.targets {
		targets = append(targets, fmt.Sprintf("%s (%s)", target.domain, target.provider))
	}
	return map[string]string{
		"CA":      p.ca,
		"Email":   p.user.email,
		"Targets": strings.Join(targets, ", "),
	}
}
