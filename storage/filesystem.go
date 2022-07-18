package storage

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
)

const (
	fsCrtName = "crt.pem"
	fsKeyName = "key.pem"
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

func (s *FileSystem) Get(name string) ([]byte, []byte, error) {
	crtData, err := ioutil.ReadFile(filepath.Join(s.path, name, fsCrtName))
	if err != nil {
		return nil, nil, err
	}

	keyData, err := ioutil.ReadFile(filepath.Join(s.path, name, fsKeyName))
	if err != nil {
		return nil, nil, err
	}

	return crtData, keyData, nil
}

func (s *FileSystem) Put(name string, crtData *pem.Block, keyData *pem.Block) error {
	var (
		dirPath = filepath.Join(s.path, name)
		crtPath = filepath.Join(dirPath, fsCrtName)
		keyPath = filepath.Join(dirPath, fsKeyName)
	)
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	if err := pki.Export(crtData, crtPath); err != nil {
		return err
	}

	if err := pki.Export(keyData, keyPath); err != nil {
		return err
	}

	return nil
}

func (s *FileSystem) Del(name string) error {
	// needed as os.RemoveAll does not return an error
	// when the directory does not exist
	if _, err := os.Stat(filepath.Join(s.path, name)); errors.Is(err, os.ErrNotExist) {
		return err
	}

	return os.RemoveAll(filepath.Join(s.path, name))
}
