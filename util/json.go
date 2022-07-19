package util

import "strings"

func BytesToJSON(data []byte) string {
	return strings.ReplaceAll(strings.TrimSuffix(string(data), "\n"), "\n", "\\n")
}
