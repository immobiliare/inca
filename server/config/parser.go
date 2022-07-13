package config

import (
	"io/ioutil"
	"os"
	"strings"

	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/storage"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Providers []*provider.Provider
	Origins   []struct {
		Type    string `yaml:"type"`
		Options string `yaml:"options"`
	} `yaml:"origins"`

	Storage *storage.Storage
	Data    struct {
		Type    string `yaml:"type"`
		Options string `yaml:"options"`
	} `yaml:"data"`
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

	storage, err := storage.Get(cfg.Data.Type, strings.Split(cfg.Data.Options, " ")...)
	if err != nil {
		return nil, err
	}
	cfg.Storage = storage

	return &cfg, nil
}
