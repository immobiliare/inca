package config

import (
	"io/ioutil"
	"os"
	"strings"

	"gitlab.rete.farm/sistemi/inca/provider"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Providers []*provider.Provider
	Origins   []struct {
		Type    string `yaml:"type"`
		Options string `yaml:"options"`
	} `yaml:"origins"`
}

func Parse(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Config{Providers: []*provider.Provider{}}
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	for _, provPuppet := range cfg.Origins {
		prov, err := provider.Get(provPuppet.Type, strings.Split(provPuppet.Options, " ")...)
		if err != nil {
			return nil, err
		}
		cfg.Providers = append(cfg.Providers, prov)
	}

	return &cfg, nil
}
