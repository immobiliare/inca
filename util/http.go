package util

import (
	"bytes"
	"net/url"
)

func ParseQueryString(queryStrings []byte) map[string]string {
	queryStringMap := make(map[string]string)
	for _, queryString := range bytes.Split(queryStrings, []byte("&")) {
		pts := bytes.Split(queryString, []byte("="))
		if len(pts) != 2 {
			continue
		}
		data, err := url.QueryUnescape(string(pts[1]))
		if err != nil {
			data = string(pts[1])
		}
		queryStringMap[string(pts[0])] = data
	}
	return queryStringMap
}
