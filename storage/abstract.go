package storage

import (
	"encoding/pem"
	"errors"
	"strings"
)

type Storage interface {
	ID() string
	Tune(options ...string) error
	Put(name string, crtData *pem.Block, keyData *pem.Block) error
	Get(name string) ([]byte, []byte, error)
	Del(name string) error
	Find(filters ...string) ([][]byte, error)
}

func Get(id string, options ...string) (*Storage, error) {
	for _, storage := range []Storage{
		new(FileSystem),
		new(S3),
	} {
		if !strings.EqualFold(id, storage.ID()) {
			continue
		}
		if err := storage.Tune(options...); err != nil {
			return nil, err
		}
		return &storage, nil
	}

	return nil, errors.New("storage not found")
}
