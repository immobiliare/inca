package util

import (
	"bytes"
)

func ParseQueryString(queryStrings []byte) map[string]string {
	queryStringMap := make(map[string]string)
	for _, queryString := range bytes.Split(queryStrings, []byte("&")) {
		pts := bytes.Split(queryString, []byte("="))
		if len(pts) != 2 {
			continue
		}
		queryStringMap[string(pts[0])] = string(pts[1])
	}
	return queryStringMap
}
