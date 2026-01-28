package encjson

import (
	"bytes"
	"encoding/json"
	"testing"
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
	`日本語`:            false,
	`0123456789`:     false,
}

func TestAppendString(t *testing.T) {
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
			if expected != actual {
				t.Fatal("appending string")
			}
		})
	}
	for str := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			expected := refEncodeString(str)
			actual := AppendString(nil, str)
			if !bytes.Equal(expected, actual) {
				t.Fatal("encoding string")
			}
		})
	}
}

func TestAppendStringBytes(t *testing.T) {
	for str := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			expected := refEncodeString(str)
			actual := AppendStringBytes(nil, []byte(str))
			if !bytes.Equal(expected, actual) {
				t.Fatal("encoding string from bytes")
			}
		})
	}
}

// TestAppendStringBytesWithPrefix tests AppendStringBytes with various buffer prefixes
func TestAppendStringBytesWithPrefix(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		input    string
		expected string
	}{
		{"empty_buf", "", "hello", `"hello"`},
		{"after_colon", ":", "hello", `:"hello"`},
		{"after_open_brace", "{", "hello", `{"hello"`},
		{"after_open_bracket", "[", "hello", `["hello"`},
		{"after_close_brace", "}", "hello", `},"hello"`},
		{"after_close_bracket", "]", "hello", `],"hello"`},
		{"after_quote", `"`, "hello", `","hello"`},
		{"after_number", "1", "hello", `1,"hello"`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendStringBytes([]byte(tc.buf), []byte(tc.input)))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// TestAppendStringBytesEscapeCharacters tests all escape sequences for byte slice version
func TestAppendStringBytesEscapeCharacters(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"backslash", []byte(`\`), `"\\"`},
		{"double_quote", []byte(`"`), `"\""`},
		{"newline", []byte("\n"), `"\n"`},
		{"carriage_return", []byte("\r"), `"\r"`},
		{"tab", []byte("\t"), `"\t"`},
		{"backspace", []byte("\b"), `"\b"`},
		{"form_feed", []byte("\f"), `"\f"`},
		{"null_byte", []byte{0x00}, `"\u0000"`},
		{"control_char_1", []byte{0x01}, `"\u0001"`},
		{"control_char_1f", []byte{0x1f}, `"\u001F"`},
		{"mixed_escapes", []byte("a\nb\\c\"d\te\f\b"), `"a\nb\\c\"d\te\f\b"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendStringBytes(nil, tc.input))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
			// Verify result is valid JSON by unmarshaling
			var s string
			if err := json.Unmarshal([]byte(actual), &s); err != nil {
				t.Errorf("result is not valid JSON: %v", err)
			}
			if s != string(tc.input) {
				t.Errorf("round-trip failed: got %q, want %q", s, string(tc.input))
			}
		})
	}
}

// TestAppendStringBytesLineSeparators tests U+2028 and U+2029 escaping for byte slice
func TestAppendStringBytesLineSeparators(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"line_separator", []byte("\u2028"), `"\u2028"`},
		{"paragraph_separator", []byte("\u2029"), `"\u2029"`},
		{"mixed_separators", []byte("a\u2028b\u2029c"), `"a\u2028b\u2029c"`},
		{"separator_with_newline", []byte("\u2028\n\u2029"), `"\u2028\n\u2029"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendStringBytes(nil, tc.input))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestStringNeedsEscaping(t *testing.T) {
	for str, expected := range testStringsEscape {
		t.Run(str, func(t *testing.T) {
			actual := StringNeedsEscaping(str)
			if expected != actual {
				t.Fatal("string needs escaping")
			}
		})
	}
}

// TestAppendStringEscapeCharacters tests all JSON escape sequences
func TestAppendStringEscapeCharacters(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"backslash", `\`, `"\\"`},
		{"double_quote", `"`, `"\""`},
		{"newline", "\n", `"\n"`},
		{"carriage_return", "\r", `"\r"`},
		{"tab", "\t", `"\t"`},
		{"backspace", "\b", `"\b"`},
		{"form_feed", "\f", `"\f"`},
		{"null_byte", "\x00", `"\u0000"`},
		{"control_char_1", "\x01", `"\u0001"`},
		{"control_char_1f", "\x1f", `"\u001F"`}, // Note: encoding/json uses lowercase, but uppercase is valid JSON
		{"mixed_escapes", "a\nb\\c\"d\te", `"a\nb\\c\"d\te"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendString(nil, tc.input))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
			// Verify result is valid JSON by unmarshaling
			var s string
			if err := json.Unmarshal([]byte(actual), &s); err != nil {
				t.Errorf("result is not valid JSON: %v", err)
			}
			if s != tc.input {
				t.Errorf("round-trip failed: got %q, want %q", s, tc.input)
			}
		})
	}
}

// TestAppendStringUnicode tests various Unicode characters
func TestAppendStringUnicode(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"emoji", "Hello 👋 World 🌍"},
		{"chinese", "你好世界"},
		{"japanese", "こんにちは世界"},
		{"arabic", "مرحبا بالعالم"},
		{"hebrew", "שלום עולם"},
		{"mixed", "Hello 世界 🌍"},
		{"surrogate_pair", "𝄞"}, // Musical G clef (U+1D11E)
		{"zero_width_joiner", "👨‍👩‍👧‍👦"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := AppendString(nil, tc.input)
			ref := refEncodeString(tc.input)
			if !bytes.Equal(ref, actual) {
				t.Errorf("differs from encoding/json:\n  expected: %s\n  got:      %s", ref, actual)
			}
		})
	}
}

// TestAppendStringLineSeparators tests U+2028 and U+2029 escaping
func TestAppendStringLineSeparators(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"line_separator", "\u2028", `"\u2028"`},
		{"paragraph_separator", "\u2029", `"\u2029"`},
		{"mixed_separators", "a\u2028b\u2029c", `"a\u2028b\u2029c"`},
		{"separator_with_newline", "\u2028\n\u2029", `"\u2028\n\u2029"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendString(nil, tc.input))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}

	// Also verify StringNeedsEscaping detects these
	if !StringNeedsEscaping("\u2028") {
		t.Error("StringNeedsEscaping should return true for U+2028")
	}
	if !StringNeedsEscaping("\u2029") {
		t.Error("StringNeedsEscaping should return true for U+2029")
	}
}

// TestAppendStringInvalidUTF8 tests handling of invalid UTF-8 sequences
func TestAppendStringInvalidUTF8(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"single_invalid_byte", "\x80"},
		{"invalid_continuation", "\xc0\x80"},
		{"truncated_sequence", "\xc2"},
		{"invalid_in_middle", "abc\xfexyz"},
		{"overlong_encoding", "\xc0\xaf"},
		{"invalid_start_byte", "\xff"},
		{"mixed_valid_invalid", "hello\xffworld"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			actual := AppendString(nil, tc.input)
			// Result should be valid JSON (unmarshalable)
			var s string
			err := json.Unmarshal(actual, &s)
			if err != nil {
				t.Errorf("result is not valid JSON: %v, got %q", err, actual)
			}
		})
	}
}

// TestAppendStringBytesInvalidUTF8 tests handling of invalid UTF-8 in byte slices
func TestAppendStringBytesInvalidUTF8(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{"single_invalid_byte", []byte{0x80}},
		{"invalid_continuation", []byte{0xc0, 0x80}},
		{"truncated_sequence", []byte{0xc2}},
		{"invalid_in_middle", []byte("abc\xfexyz")},
		{"overlong_encoding", []byte{0xc0, 0xaf}},
		{"invalid_start_byte", []byte{0xff}},
		{"mixed_valid_invalid", []byte("hello\xffworld")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			actual := AppendStringBytes(nil, tc.input)
			// Result should be valid JSON (unmarshalable)
			var s string
			err := json.Unmarshal(actual, &s)
			if err != nil {
				t.Errorf("result is not valid JSON: %v, got %q", err, actual)
			}
		})
	}
}

// TestAppendStringFastPath tests the fast path for strings without escaping
func TestAppendStringFastPath(t *testing.T) {
	testCases := []string{
		"hello world",
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"!@#$%^*()_+-=[]{}|;':,./",     // safe ASCII punctuation
		"!@#$%^&*()_+-=[]{}|;':,./<>?", // includes <, >, & which encoding/json escapes for HTML safety
		" ",                             // just space
		"a b c",
		"", // empty string
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			actual := AppendString(nil, input)
			// Verify result is valid JSON by unmarshaling and round-tripping
			var s string
			if err := json.Unmarshal(actual, &s); err != nil {
				t.Errorf("result is not valid JSON: %v, got %q", err, actual)
			}
			if s != input {
				t.Errorf("round-trip failed: got %q, want %q", s, input)
			}
		})
	}
}

// TestAppendStringSlowPath tests the slow path with escape at various positions
func TestAppendStringSlowPath(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"escape_at_start", "\nhello"},
		{"escape_at_end", "hello\n"},
		{"escape_in_middle", "hel\nlo"},
		{"multiple_escapes", "a\nb\nc\nd"},
		{"long_prefix_then_escape", "abcdefghijklmnopqrstuvwxyz\n"},
		{"short_prefix_then_escape", "a\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := AppendString(nil, tc.input)
			ref := refEncodeString(tc.input)
			if !bytes.Equal(ref, actual) {
				t.Errorf("differs from encoding/json:\n  expected: %s\n  got:      %s", ref, actual)
			}
		})
	}
}

// TestAppendSafeKey tests the fast key appending
func TestAppendSafeKey(t *testing.T) {
	testCases := []struct {
		name     string
		buf      string
		key      string
		expected string
	}{
		{"after_open_brace", "{", "key", `{"key":`},
		{"after_value", `{"a":1`, "key", `{"a":1,"key":`},
		{"after_close_brace", "}", "key", `},"key":`},
		{"after_string", `"`, "key", `","key":`},
		{"empty_buf", "", "key", `"key":`},
		{"numeric_key", "{", "123", `{"123":`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(AppendSafeKey([]byte(tc.buf), tc.key))
			if tc.expected != actual {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// BenchmarkAppendString benchmarks string encoding
func BenchmarkAppendString(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"short_no_escape", "hello"},
		{"medium_no_escape", "hello world this is a test message"},
		{"long_no_escape", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."},
		{"with_escapes", "hello\nworld\ttab\"quote\\backslash"},
		{"unicode", "Hello 世界 🌍 мир"},
	}

	buf := make([]byte, 0, 256)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf = AppendString(buf[:0], bm.input)
			}
		})
	}
}

// BenchmarkAppendStringBytes benchmarks byte slice encoding
func BenchmarkAppendStringBytes(b *testing.B) {
	benchmarks := []struct {
		name  string
		input []byte
	}{
		{"short_no_escape", []byte("hello")},
		{"medium_no_escape", []byte("hello world this is a test message")},
		{"with_escapes", []byte("hello\nworld\ttab\"quote\\backslash")},
	}

	buf := make([]byte, 0, 256)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf = AppendStringBytes(buf[:0], bm.input)
			}
		})
	}
}
