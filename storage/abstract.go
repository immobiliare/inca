package storage

import (
	"errors"
	"strings"
)

type Storage interface {
	ID() string
	Tune(options map[string]interface{}) error
	Put(name string, crtData, keyData []byte) error
	Get(name string) ([]byte, []byte, error)
	Del(name string) error
	Find(filters ...string) ([][]byte, error)
	Config() map[string]string
}

func Get(id string, options map[string]interface{}) (*Storage, error) {
	for _, storage := range []Storage{
		new(FS),
		new(S3),
		new(PostgreSQL),
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
