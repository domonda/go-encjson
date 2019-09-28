package encjson

func AppendObjectStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ',', '{')
	}
	return append(b, '{')
}

func AppendObjectEnd(b []byte) []byte {
	return append(b, '}')
}

func AppendKey(b []byte, key string) []byte {
	return append(AppendString(b, key), ':')
}
