package encjson

// AppendKey appends a JSON object key (the key string with quotes and a trailing colon) to b.
// A comma separator is prepended if needed based on the context (previous character).
// The key is properly escaped for JSON using AppendString.
func AppendKey(b []byte, key string) []byte {
	return append(AppendString(b, key), ':')
}

// AppendSafeKey appends a JSON object key that is known to contain only safe ASCII characters
// (no escaping needed). This is faster than AppendKey for known-safe keys like struct field names.
// The key must only contain printable ASCII characters (0x20-0x7e) excluding '"' and '\'.
func AppendSafeKey(b []byte, key string) []byte {
	if l := len(b); l > 0 && b[l-1] != '{' {
		b = append(b, ',')
	}
	b = append(b, '"')
	b = append(b, key...)
	return append(b, '"', ':')
}

// AppendObjectStart appends the JSON object start character '{' to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendObjectStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		return append(b, ',', '{')
	}
	return append(b, '{')
}

// AppendObjectEnd appends the JSON object end character '}' to b.
func AppendObjectEnd(b []byte) []byte {
	return append(b, '}')
}
