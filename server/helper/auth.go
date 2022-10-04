package helper

import "strings"

func IsValidToken(token string, acl map[string][]string) bool {
	for key := range acl {
		if strings.EqualFold(key, token) {
			return true
		}
	}
	return false
}
