package util

import (
	"testing"

	"github.com/matryer/is"
)

const (
	testingQueryStringOk  = "key1=value1&key2=value2&key3"
	testingQueryStringNok = "key"
)

func TestUtilHttpParseQueryStringOk(t *testing.T) {
	t.Parallel()

	var (
		test        = is.New(t)
		queryString = ParseQueryString([]byte(testingQueryStringOk))
	)

	key, ok := queryString["key1"]
	test.True(ok)
	test.Equal(key, "value1")

	key, ok = queryString["key2"]
	test.True(ok)
	test.Equal(key, "value2")

	_, ok = queryString["key3"]
	test.True(!ok)
}

func TestUtilHttpParseQueryStringNok(t *testing.T) {
	t.Parallel()

	is.New(t).True(len(ParseQueryString([]byte(testingQueryStringNok))) == 0)
}
