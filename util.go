package encjson

const hexChars = "0123456789ABCDEF"

func lastChar(b []byte) byte {
	l := len(b)
	if l == 0 {
		return 0
	}
	return b[l-1]
}

// func growBytesCap(b []byte, n int) []byte {
// 	l, c := len(b), cap(b)
// 	if l+n <= c {
// 		return b
// 	}
// 	newCap := l + n // TODO better growing
// 	newBuf := make([]byte, l, newCap)
// 	copy(newBuf, b)
// 	return newBuf
// }

// func growBytesLen(b []byte, n int) []byte {
// 	l, c := len(b), cap(b)
// 	if l+n <= c {
// 		return b[:l+n]
// 	}
// 	newCap := l + n // TODO better growing
// 	newBuf := make([]byte, l+n, newCap)
// 	copy(newBuf, b)
// 	return newBuf
// }
