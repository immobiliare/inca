package util

import "regexp"

func RegexesMatch(payload string, regexes ...string) bool {
	match := true
	for _, filter := range regexes {
		match = match && ErrWrap(false)(regexp.MatchString(filter, payload))
	}
	return match
}
