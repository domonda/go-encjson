package encjson

import "unicode/utf8"

func StringNeedsEscaping(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == '\\' || r == '"' {
			return true
		}
	}
	return false
}

func String(s string) []byte {
	return AppendString(nil, s)
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
			b = append(b, s[:l]...)

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
