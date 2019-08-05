package encjson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func refEncodeString(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return b
}

var testStrings = []string{
	"",
	"\t'\r\n`\u0001",
	`"Hello World!"`,
	`Österreich`,
	`日本語`,
	`0123456789`,
}

func Test_String(t *testing.T) {
	for _, str := range testStrings {
		t.Run(str, func(t *testing.T) {
			expeced := refEncodeString(str)
			actual := String(str)
			assert.Equal(t, string(expeced), string(actual), "encoding string")
		})
	}
}

func Test_AppendStringBytes(t *testing.T) {
	for _, str := range testStrings {
		t.Run(str, func(t *testing.T) {
			expeced := refEncodeString(str)
			actual := AppendStringBytes(nil, []byte(str))
			assert.Equal(t, string(expeced), string(actual), "encoding string from bytes")
		})
	}
}
