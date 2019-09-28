package encjson

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/domonda/go-types/uu"
)

func Test_AppendUUID(t *testing.T) {
	uuid := uu.IDMustFromString("7d11f476-72c9-4df4-b17f-0074e3cb2e7a")
	testCases := map[string]string{
		``:  `"7d11f476-72c9-4df4-b17f-0074e3cb2e7a"`,
		`:`: `:"7d11f476-72c9-4df4-b17f-0074e3cb2e7a"`,
		`"`: `","7d11f476-72c9-4df4-b17f-0074e3cb2e7a"`,
	}
	for buf, expected := range testCases {
		t.Run(buf, func(t *testing.T) {
			actual := string(AppendUUID([]byte(buf), uuid))
			assert.Equal(t, expected, actual, "appending UUID")
		})
	}
}
