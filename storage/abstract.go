package storage

import (
	"encoding/pem"
	"errors"
	"strings"
)

type Storage interface {
	ID() string
	Tune(options map[string]interface{}) error
	Put(name string, crtData *pem.Block, keyData *pem.Block) error
	Get(name string) ([]byte, []byte, error)
	Del(name string) error
	Find(filters ...string) ([][]byte, error)
}

func Get(id string, options map[string]interface{}) (*Storage, error) {
	for _, storage := range []Storage{
		new(FS),
		new(S3),
	} {
		if !strings.EqualFold(id, storage.ID()) {
			continue
		}

		if err := storage.Tune(options); err != nil {
			return nil, err
		}

		return &storage, nil
	}

	return nil, errors.New("storage not found")
}
