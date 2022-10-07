package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/util"
)

const (
	fsCrtName = "crt.pem"
	fsKeyName = "key.pem"
)

type FS struct {
	Storage
	path string
}

func (s FS) ID() string {
	return "FS"
}

func (s *FS) Tune(options map[string]interface{}) error {
	path, ok := options["path"]
	if !ok {
		return fmt.Errorf("provider %s: crt not defined", s.ID())
	}

	s.path = path.(string)
	return nil
}

func (s *FS) Get(name string) ([]byte, []byte, error) {
	crtData, err := os.ReadFile(filepath.Join(s.path, name, fsCrtName))
	if err != nil {
		return nil, nil, err
	}

	keyData, err := os.ReadFile(filepath.Join(s.path, name, fsKeyName))
	if err != nil {
		return nil, nil, err
	}

	return crtData, keyData, nil
}

func (s *FS) Put(name string, crtData, keyData []byte) error {
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

	if err := os.WriteFile(crtPath, crtData, 0644); err != nil {
		return err
	}

	if err := os.WriteFile(keyPath, keyData, 0644); err != nil {
		return err
	}

	return nil
}

func (s *FS) Del(name string) error {
	// needed as os.RemoveAll does not return an error
	// when the directory does not exist
	if _, err := os.Stat(filepath.Join(s.path, name)); errors.Is(err, os.ErrNotExist) {
		return err
	}

	return os.RemoveAll(filepath.Join(s.path, name))
}

func (s *FS) Find(filters ...string) ([][]byte, error) {
	dirs, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	results := [][]byte{}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		var (
			crtPath  = filepath.Join(s.path, dir.Name(), fsCrtName)
			isResult = pki.IsValidCN(dir.Name()) && util.RegexesMatch(dir.Name(), filters...)
		)

		_, err := os.Stat(crtPath)
		isResult = isResult && !os.IsNotExist(err)
		if !isResult {
			continue
		}

		crt, _, err := s.Get(dir.Name())
		if err != nil {
			return nil, err
		}

		results = append(results, crt)
	}

	return results, nil
}

func (s *FS) Config() map[string]string {
	return map[string]string{
		"Path": s.path,
	}
}
