package encjson

import "unicode/utf8"

const hexChars = "0123456789ABCDEF"

// noEscapeTable is a lookup table for ASCII characters that don't need escaping.
// Characters that need escaping: control chars (< 0x20), backslash, double quote.
// Characters > 0x7e are handled by the multi-byte path.
var noEscapeTable = [256]bool{
	// 0x20-0x21: space and !
	0x20: true, 0x21: true,
	// 0x22 is " - needs escaping, skip
	// 0x23-0x5b: # through [
	0x23: true, 0x24: true, 0x25: true, 0x26: true, 0x27: true,
	0x28: true, 0x29: true, 0x2a: true, 0x2b: true, 0x2c: true,
	0x2d: true, 0x2e: true, 0x2f: true,
	0x30: true, 0x31: true, 0x32: true, 0x33: true, 0x34: true,
	0x35: true, 0x36: true, 0x37: true, 0x38: true, 0x39: true,
	0x3a: true, 0x3b: true, 0x3c: true, 0x3d: true, 0x3e: true,
	0x3f: true,
	0x40: true, 0x41: true, 0x42: true, 0x43: true, 0x44: true,
	0x45: true, 0x46: true, 0x47: true, 0x48: true, 0x49: true,
	0x4a: true, 0x4b: true, 0x4c: true, 0x4d: true, 0x4e: true,
	0x4f: true,
	0x50: true, 0x51: true, 0x52: true, 0x53: true, 0x54: true,
	0x55: true, 0x56: true, 0x57: true, 0x58: true, 0x59: true,
	0x5a: true, 0x5b: true,
	// 0x5c is \ - needs escaping, skip
	// 0x5d-0x7e: ] through ~
	0x5d: true, 0x5e: true, 0x5f: true,
	0x60: true, 0x61: true, 0x62: true, 0x63: true, 0x64: true,
	0x65: true, 0x66: true, 0x67: true, 0x68: true, 0x69: true,
	0x6a: true, 0x6b: true, 0x6c: true, 0x6d: true, 0x6e: true,
	0x6f: true,
	0x70: true, 0x71: true, 0x72: true, 0x73: true, 0x74: true,
	0x75: true, 0x76: true, 0x77: true, 0x78: true, 0x79: true,
	0x7a: true, 0x7b: true, 0x7c: true, 0x7d: true, 0x7e: true,
}

// StringNeedsEscaping reports whether the string s contains any characters
// that need to be escaped when encoding as JSON. This includes control characters
// (< 0x20), backslash, double quote, and the special Unicode characters U+2028
// (LINE SEPARATOR) and U+2029 (PARAGRAPH SEPARATOR).
func StringNeedsEscaping(s string) bool {
	for i := range len(s) {
		c := s[i]
		if c >= 0x80 {
			// Multi-byte UTF-8, check for U+2028 and U+2029
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == '\u2028' || r == '\u2029' {
				return true
			}
			i += size - 1 // -1 because loop will increment
			continue
		}
		if !noEscapeTable[c] {
			return true
		}
	}
	return false
}

// AppendString appends the JSON encoding of string s to b, including surrounding
// double quotes. It handles proper escaping of control characters, backslash,
// double quote, and the special Unicode characters U+2028 and U+2029.
// A comma separator is prepended if the last character of b is not '{', '[', ':', or empty.
// Invalid UTF-8 sequences are replaced with the Unicode replacement character U+FFFD.
func AppendString(b []byte, s string) []byte {
	switch lastChar(b) {
	case 0, '{', '[', ':':
		b = append(b, '"')
	default:
		b = append(b, ',', '"')
	}

	// Fast path: scan for characters that need escaping
	for i := range len(s) {
		c := s[i]
		if c >= 0x80 {
			// Multi-byte UTF-8, need to check for U+2028/U+2029
			return appendStringComplex(b, s, i)
		}
		if !noEscapeTable[c] {
			// Found a character that needs escaping
			return appendStringComplex(b, s, i)
		}
	}

	// No escaping needed - bulk append
	b = append(b, s...)
	return append(b, '"')
}

// appendStringComplex handles strings that need character-by-character processing,
// starting from the position where special handling is needed.
func appendStringComplex(b []byte, s string, start int) []byte {
	// Append the safe prefix as a bulk copy
	b = append(b, s[:start]...)

	for i := start; i < len(s); {
		c := s[i]

		// Fast path for safe ASCII characters
		if c < 0x80 && noEscapeTable[c] {
			b = append(b, c)
			i++
			continue
		}

		// Multi-byte UTF-8
		if c >= 0x80 {
			r, size := utf8.DecodeRuneInString(s[i:])
			// U+2028 is LINE SEPARATOR.
			// U+2029 is PARAGRAPH SEPARATOR.
			// They are both technically valid characters in JSON strings,
			// but don't work in JSONP, which has to be evaluated as JavaScript,
			// and can lead to security holes there. It is valid JSON to
			// escape them, so we do so unconditionally.
			// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
			if r == '\u2028' || r == '\u2029' {
				b = append(b, `\u202`...)
				b = append(b, hexChars[r&0xF])
			} else {
				b = append(b, s[i:i+size]...)
			}
			i += size
			continue
		}

		// Handle ASCII characters that need escaping
		switch c {
		case '"':
			b = append(b, '\\', '"')
		case '\\':
			b = append(b, '\\', '\\')
		case '\n':
			b = append(b, '\\', 'n')
		case '\f':
			b = append(b, '\\', 'f')
		case '\b':
			b = append(b, '\\', 'b')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			// Control characters < 0x20
			b = append(b, '\\', 'u', '0', '0', hexChars[c>>4&0xF], hexChars[c&0xF])
		}
		i++
	}

	return append(b, '"')
}

// AppendStringBytes appends the JSON encoding of the byte slice s (interpreted as UTF-8)
// to b, including surrounding double quotes. It handles proper escaping of control
// characters, backslash, double quote, and the special Unicode characters U+2028 and U+2029.
// A comma separator is prepended if the last character of b is not '{', '[', ':', or empty.
// Invalid UTF-8 sequences are replaced with the Unicode replacement character U+FFFD.
func AppendStringBytes(b []byte, s []byte) []byte {
	switch lastChar(b) {
	case 0, '{', '[', ':':
		b = append(b, '"')
	default:
		b = append(b, ',', '"')
	}

	// Fast path: scan for characters that need escaping
	for i := range len(s) {
		c := s[i]
		if c >= 0x80 {
			// Multi-byte UTF-8, need to check for U+2028/U+2029
			return appendStringBytesComplex(b, s, i)
		}
		if !noEscapeTable[c] {
			// Found a character that needs escaping
			return appendStringBytesComplex(b, s, i)
		}
	}

	// No escaping needed - bulk append
	b = append(b, s...)
	return append(b, '"')
}

// appendStringBytesComplex handles byte slices that need character-by-character processing,
// starting from the position where special handling is needed.
func appendStringBytesComplex(b []byte, s []byte, start int) []byte {
	// Append the safe prefix as a bulk copy
	b = append(b, s[:start]...)
	s = s[start:]

	for len(s) > 0 {
		c := s[0]

		// Fast path for safe ASCII characters
		if c < 0x80 && noEscapeTable[c] {
			b = append(b, c)
			s = s[1:]
			continue
		}

		// Multi-byte UTF-8
		if c >= 0x80 {
			r, size := utf8.DecodeRune(s)
			// U+2028 is LINE SEPARATOR.
			// U+2029 is PARAGRAPH SEPARATOR.
			// They are both technically valid characters in JSON strings,
			// but don't work in JSONP, which has to be evaluated as JavaScript,
			// and can lead to security holes there. It is valid JSON to
			// escape them, so we do so unconditionally.
			// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
			if r == '\u2028' || r == '\u2029' {
				b = append(b, `\u202`...)
				b = append(b, hexChars[r&0xF])
			} else {
				b = append(b, s[:size]...)
			}
			s = s[size:]
			continue
		}

		// Handle ASCII characters that need escaping
		switch c {
		case '"':
			b = append(b, '\\', '"')
		case '\\':
			b = append(b, '\\', '\\')
		case '\n':
			b = append(b, '\\', 'n')
		case '\f':
			b = append(b, '\\', 'f')
		case '\b':
			b = append(b, '\\', 'b')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			// Control characters < 0x20
			b = append(b, '\\', 'u', '0', '0', hexChars[c>>4&0xF], hexChars[c&0xF])
		}
		s = s[1:]
	}

	return append(b, '"')
}

// lastChar returns the last byte of b, or 0 if b is empty.
func lastChar(b []byte) byte {
	l := len(b)
	if l == 0 {
		return 0
	}
	return b[l-1]
}
