package encjson

import (
	"testing"
)

func Test_AppendUUID(t *testing.T) {
	uuid := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	testCases := map[string]string{
		``:  `"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`,
		`:`: `:"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`,
		`"`: `","6ba7b810-9dad-11d1-80b4-00c04fd430c8"`,
	}
	for buf, expected := range testCases {
		t.Run(buf, func(t *testing.T) {
			actual := string(AppendUUID([]byte(buf), uuid))
			if expected != actual {
				t.Fatal("appending UUID")
			}
		})
	}
}
