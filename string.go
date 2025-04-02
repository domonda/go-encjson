package encjson

import "unicode/utf8"

const hexChars = "0123456789ABCDEF"

func StringNeedsEscaping(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == '\\' || r == '"' {
			return true
		}
	}
	return false
}

func AppendString(b []byte, s string) []byte {
	switch lastChar(b) {
	case 0, '{', '[', ':':
		b = append(b, '"')
	default:
		b = append(b, ',', '"')
	}

	for i, r := range s {
		if l := utf8.RuneLen(r); l > 1 {
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
				continue
			}

			b = append(b, s[i:i+l]...)
			continue
		}

		switch r {
		case '"':
			b = append(b, '\\', '"')

		case '\\':
			b = append(b, '\\', '\\')

		case '\n':
			b = append(b, '\\', 'n')

		case '\f':
			b = append(b, '\\', 't')

		case '\b':
			b = append(b, '\\', 'b')

		case '\r':
			b = append(b, '\\', 'r')

		case '\t':
			b = append(b, '\\', 't')

		default:
			if r < 0x20 {
				b = append(b, '\\', 'u', '0', '0', hexChars[r>>4&0xF], hexChars[r&0xF])
			} else {
				b = append(b, byte(r))
			}
		}
	}

	return append(b, '"')
}

func AppendStringBytes(b []byte, s []byte) []byte {
	switch lastChar(b) {
	case 0, '{', '[', ':':
		b = append(b, '"')
	default:
		b = append(b, ',', '"')
	}

	for r, l := utf8.DecodeRune(s); l > 0; {
		if l > 1 {
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
				b = append(b, s[:l]...)
			}

			s = s[l:]
			r, l = utf8.DecodeRune(s)
			continue
		}

		switch r {
		case '"':
			b = append(b, '\\', '"')

		case '\\':
			b = append(b, '\\', '\\')

		case '\n':
			b = append(b, '\\', 'n')

		case '\f':
			b = append(b, '\\', 't')

		case '\b':
			b = append(b, '\\', 'b')

		case '\r':
			b = append(b, '\\', 'r')

		case '\t':
			b = append(b, '\\', 't')

		default:
			if r < 0x20 {
				b = append(b, '\\', 'u', '0', '0', hexChars[r>>4&0xF], hexChars[r&0xF])
			} else {
				b = append(b, byte(r))
			}
		}

		s = s[l:]
		r, l = utf8.DecodeRune(s)
	}

	return append(b, '"')
}

func lastChar(b []byte) byte {
	l := len(b)
	if l == 0 {
		return 0
	}
	return b[l-1]
}
