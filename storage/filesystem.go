package storage

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
)

type FileSystem struct {
	Storage
	path string
}

func (s FileSystem) ID() string {
	return "filesystem"
}

func (s *FileSystem) Tune(options ...string) error {
	if len(options) != 1 {
		return fmt.Errorf("invalid number of options for provider %s: %s", s.ID(), options)
	}

	s.path = options[0]
	if !util.IsDir(s.path) {
		return fmt.Errorf("%s: no such directory", s.path)
	}

	return nil
}

func (s *FileSystem) Get(name string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(s.path, name))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *FileSystem) Put(name string, data *pem.Block) error {
	if err := pki.Export(data, filepath.Join(s.path, name)); err != nil {
		return err
	}
	return nil
}
