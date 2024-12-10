package util

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestUtilZStdLogger(t *testing.T) {
	t.Parallel()

	const data = "hello"

	var (
		test   = is.New(t)
		buffer = new(bytes.Buffer)
		logger = NewZStdLogger(buffer)
	)

	logger.Println(data)
	t.Log(buffer.String())

	bufferJson := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &bufferJson)
	test.NoErr(err)

	message, ok := bufferJson["message"]
	test.True(ok)
	test.Equal(message, data)
}
