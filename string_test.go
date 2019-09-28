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

var testStringsEscape = map[string]bool{
	"":               false,
	"\t'\r\n`\u0001": true,
	`"Hello World!"`: true,
	`Österreich`:     false,
	`日本語`:            false, // broken gofmt
	`0123456789`:     false,
}

func Test_String(t *testing.T) {
	for str := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			expected := refEncodeString(str)
			actual := String(str)
			assert.Equal(t, string(expected), string(actual), "encoding string")
		})
	}
}

func Test_AppendStringBytes(t *testing.T) {
	for str := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			expected := refEncodeString(str)
			actual := AppendStringBytes(nil, []byte(str))
			assert.Equal(t, string(expected), string(actual), "encoding string from bytes")
		})
	}
}

func Test_StringNeedsEscaping(t *testing.T) {
	for str, expected := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			actual := StringNeedsEscaping(str)
			assert.Equal(t, expected, actual, "string needs escaping")
		})
	}
}

func Test_AppendString(t *testing.T) {
	str := `"Hello World!"`
	testCases := map[string]string{
		``:  `"\"Hello World!\""`,
		`}`: `},"\"Hello World!\""`,
		`{`: `{"\"Hello World!\""`,
		`]`: `],"\"Hello World!\""`,
		`[`: `["\"Hello World!\""`,
		`"`: `","\"Hello World!\""`,
	}
	for buf, expected := range testCases {
		t.Run(buf, func(t *testing.T) {
			actual := string(AppendString([]byte(buf), str))
			assert.Equal(t, expected, actual, "appending string")
		})
	}
}
