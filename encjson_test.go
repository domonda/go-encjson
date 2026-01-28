package encjson

import (
	"math"
	"testing"
	"time"
)

func TestAppendNull(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		expected string
	}{
		{"empty", "", "null"},
		{"after_colon", ":", ":null"},
		{"after_bracket", "[", "[null"},
		{"after_value", `"x"`, `"x",null`},
		{"after_number", "1", "1,null"},
		{"after_brace", "}", "},null"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendNull([]byte(tc.buf)))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendBool(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		val      bool
		expected string
	}{
		{"true_empty", "", true, "true"},
		{"false_empty", "", false, "false"},
		{"true_after_colon", ":", true, ":true"},
		{"false_after_colon", ":", false, ":false"},
		{"true_after_bracket", "[", true, "[true"},
		{"false_after_bracket", "[", false, "[false"},
		{"true_after_value", `"x"`, true, `"x",true`},
		{"false_after_value", `"x"`, false, `"x",false`},
		{"true_after_number", "1", true, "1,true"},
		{"false_after_number", "1", false, "1,false"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendBool([]byte(tc.buf), tc.val))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendInt(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		val      int64
		expected string
	}{
		{"positive_empty", "", 42, "42"},
		{"negative_empty", "", -42, "-42"},
		{"zero_empty", "", 0, "0"},
		{"after_colon", ":", 123, ":123"},
		{"after_bracket", "[", 456, "[456"},
		{"after_value", `"x"`, 789, `"x",789`},
		{"max_int64", "", math.MaxInt64, "9223372036854775807"},
		{"min_int64", "", math.MinInt64, "-9223372036854775808"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendInt([]byte(tc.buf), tc.val))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendUint(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		val      uint64
		expected string
	}{
		{"positive_empty", "", 42, "42"},
		{"zero_empty", "", 0, "0"},
		{"after_colon", ":", 123, ":123"},
		{"after_bracket", "[", 456, "[456"},
		{"after_value", `"x"`, 789, `"x",789`},
		{"max_uint64", "", math.MaxUint64, "18446744073709551615"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendUint([]byte(tc.buf), tc.val))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendFloat(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		val      float64
		expected string
	}{
		{"positive_empty", "", 3.14, "3.14"},
		{"negative_empty", "", -3.14, "-3.14"},
		{"zero_empty", "", 0.0, "0"},
		{"after_colon", ":", 1.5, ":1.5"},
		{"after_bracket", "[", 2.5, "[2.5"},
		{"after_value", `"x"`, 3.5, `"x",3.5`},
		{"nan", "", math.NaN(), `"NaN"`},
		{"pos_inf", "", math.Inf(1), `"+Inf"`},
		{"neg_inf", "", math.Inf(-1), `"-Inf"`},
		{"nan_after_value", "1", math.NaN(), `1,"NaN"`},
		{"pos_inf_after_value", "1", math.Inf(1), `1,"+Inf"`},
		{"neg_inf_after_value", "1", math.Inf(-1), `1,"-Inf"`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendFloat([]byte(tc.buf), tc.val))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendTime(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	testCases := []struct {
		name     string
		buf      string
		layout   string
		expected string
	}{
		{"rfc3339_empty", "", time.RFC3339, `"2024-01-15T10:30:00Z"`},
		{"rfc3339_after_colon", ":", time.RFC3339, `:"2024-01-15T10:30:00Z"`},
		{"rfc3339_after_bracket", "[", time.RFC3339, `["2024-01-15T10:30:00Z"`},
		{"rfc3339_after_value", `"x"`, time.RFC3339, `"x","2024-01-15T10:30:00Z"`},
		{"date_only", "", "2006-01-02", `"2024-01-15"`},
		{"time_only", "", "15:04:05", `"10:30:00"`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendTime([]byte(tc.buf), testTime, tc.layout))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendJSON(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		val      string
		expected string
	}{
		{"object_empty", "", `{"a":1}`, `{"a":1}`},
		{"array_empty", "", `[1,2,3]`, `[1,2,3]`},
		{"after_colon", ":", `{"a":1}`, `:{"a":1}`},
		{"after_bracket", "[", `{"a":1}`, `[{"a":1}`},
		{"after_value", `"x"`, `{"a":1}`, `"x",{"a":1}`},
		{"null_literal", "", `null`, `null`},
		{"nested", "", `{"a":{"b":[1,2]}}`, `{"a":{"b":[1,2]}}`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendJSON([]byte(tc.buf), []byte(tc.val)))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendArrayStart(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		expected string
	}{
		{"empty", "", "["},
		{"after_colon", ":", ":["},
		{"after_value", `"x"`, `"x",[`},
		{"after_brace", "}", "},["},
		{"after_bracket", "]", "],["},
		{"after_number", "1", "1,["},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendArrayStart([]byte(tc.buf)))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAppendArrayEnd(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		expected string
	}{
		{"after_bracket", "[", "[]"},
		{"after_value", "[1", "[1]"},
		{"after_nested", "[[1]", "[[1]]"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendArrayEnd([]byte(tc.buf)))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGrowLen(t *testing.T) {
	// Test case where capacity is sufficient
	b := make([]byte, 5, 10)
	copy(b, "hello")
	result := growLen(b, 3)
	if len(result) != 8 {
		t.Errorf("expected len 8, got %d", len(result))
	}
	if cap(result) != 10 {
		t.Errorf("expected cap 10, got %d", cap(result))
	}

	// Test case where capacity needs to grow
	b2 := make([]byte, 5)
	copy(b2, "hello")
	result2 := growLen(b2, 10)
	if len(result2) != 15 {
		t.Errorf("expected len 15, got %d", len(result2))
	}
	if string(result2[:5]) != "hello" {
		t.Errorf("expected content to be preserved, got %q", string(result2[:5]))
	}
}

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
