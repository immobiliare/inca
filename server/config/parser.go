package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/immobiliare/inca/provider"
	"github.com/immobiliare/inca/storage"
	"gopkg.in/yaml.v2"
)

const (
	defaultEnvironment   = "development"
	defaultTemplatesPath = "./server/views"
)

type Config struct {
	Sentry        string
	Environment   string
	TemplatesPath string
	Storage       *storage.Storage
	Providers     []*provider.Provider
	ACL           map[string][]string
}

type Wrapper struct {
	Sentry        string                   `yaml:"sentry"`
	TemplatesPath string                   `yaml:"templates_path"`
	Environment   string                   `yaml:"environment"`
	Storage       map[string]interface{}   `yaml:"storage"`
	Providers     []map[string]interface{} `yaml:"providers"`
	ACL           map[string][]string      `yaml:"acl"`
}

func Parse(path string) (*Config, error) {
	cfg := Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	content, err := os.ReadFile(path)
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

		provider, err := provider.Tune(id.(string), providerConfig)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	cfg.Providers = providers

	if len(wrapper.Sentry) > 0 {
		cfg.Sentry = wrapper.Sentry
	}

	cfg.Environment = defaultEnvironment
	if len(wrapper.Environment) > 0 {
		cfg.Environment = wrapper.Environment
	}

	cfg.TemplatesPath = defaultTemplatesPath
	if len(wrapper.TemplatesPath) > 0 {
		cfg.TemplatesPath = wrapper.TemplatesPath
	}

	if len(wrapper.ACL) > 0 {
		cfg.ACL = wrapper.ACL
	}

	return &cfg, nil
}
