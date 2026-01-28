package encjson

// AppendArrayStart appends the JSON array start character '[' to b.
// A comma separator is prepended if the last character of b is not ':' or empty.
func AppendArrayStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ',', '[')
	}
	return append(b, '[')
}

// AppendArrayEnd appends the JSON array end character ']' to b.
func AppendArrayEnd(b []byte) []byte {
	return append(b, ']')
}
