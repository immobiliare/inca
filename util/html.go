package util

import (
	"encoding/xml"
	"io"
	"strings"
)

func IsValidHTML(body []byte) bool {
	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	decoder.Strict = true
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity
	for {
		_, err := decoder.Token()
		switch err {
		case io.EOF:
			return true
		case nil:
		default:
			return false
		}
	}
}
