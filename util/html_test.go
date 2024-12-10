package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestUtilHTML(t *testing.T) {
	t.Parallel()

	test := is.New(t)
	test.True(IsValidHTML([]byte(`<html></html>`)))
	test.True(!IsValidHTML([]byte(`<html></div>`)))
	test.True(!IsValidHTML([]byte(`<</>`)))
}
