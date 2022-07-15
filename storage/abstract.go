package storage

import (
	"encoding/pem"
	"errors"
)

type Storage interface {
	ID() string
	Tune(options ...string) error
	Put(name string, data *pem.Block) error
	Get(name string) ([]byte, error)
}

func Get(id string, options ...string) (*Storage, error) {
	for _, storage := range []Storage{
		new(FileSystem),
	} {
		if id != storage.ID() {
			continue
		}
		if err := storage.Tune(options...); err != nil {
			return nil, err
		}
		return &storage, nil
	}

	return nil, errors.New("storage not found")
}
