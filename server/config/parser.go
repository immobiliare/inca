package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/storage"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Storage   *storage.Storage
	Providers []*provider.Provider
}

type Wrapper struct {
	Storage   map[string]interface{}   `yaml:"storage"`
	Providers []map[string]interface{} `yaml:"providers"`
}

func Parse(path string) (*Config, error) {
	cfg := Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	wrapper := Wrapper{}
	if err := yaml.Unmarshal(content, &wrapper); err != nil {
		return nil, err
	}

	id, ok := wrapper.Storage["type"]
	if !ok {
		return nil, errors.New("storage type not defined")
	}
	storage, err := storage.Get(id.(string), wrapper.Storage)
	if err != nil {
		return nil, err
	}
	cfg.Storage = storage

	providers := []*provider.Provider{}
	for _, providerConfig := range wrapper.Providers {
		id, ok := providerConfig["type"]
		if !ok {
			return nil, fmt.Errorf("provider type not defined")
		}

		provider, err := provider.Get(id.(string), providerConfig)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	cfg.Providers = providers

	return &cfg, nil
}
