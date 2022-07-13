package util

import "os"

func IsDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	fileStat, err := file.Stat()
	if err != nil {
		return false
	}
	return fileStat.IsDir()
}
