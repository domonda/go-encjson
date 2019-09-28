package encjson

func AppendArrayStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ',', '[')
	}
	return append(b, '[')
}

func AppendArrayEnd(b []byte) []byte {
	return append(b, ']')
}
