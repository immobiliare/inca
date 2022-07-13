package storage

import (
	"fmt"

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
