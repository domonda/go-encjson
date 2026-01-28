package encjson

import (
	"encoding/hex"
	"math"
	"strconv"
	"time"
)

// AppendNull appends the JSON null literal to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendNull(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		return append(b, ",null"...)
	}
	return append(b, "null"...)
}

// AppendBool appends the JSON boolean literal "true" or "false" to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendBool(b []byte, val bool) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		switch val {
		case true:
			return append(b, ",true"...)
		default:
			return append(b, ",false"...)
		}
	}

	switch val {
	case true:
		return append(b, "true"...)
	default:
		return append(b, "false"...)
	}
}

// AppendInt appends the decimal string representation of val to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendInt(b []byte, val int64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',')
	}
	return strconv.AppendInt(b, val, 10)
}

// AppendUint appends the decimal string representation of val to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendUint(b []byte, val uint64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',')
	}
	return strconv.AppendUint(b, val, 10)
}

// AppendFloat appends the string representation of val to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
// Special values NaN, +Inf, and -Inf are encoded as quoted strings "NaN", "+Inf", and "-Inf"
// since JSON does not support these values natively.
func AppendFloat(b []byte, val float64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',')
	}
	// Special floats like NaN and Inf have to be written as quoted strings
	switch {
	case math.IsNaN(val):
		return append(b, `"NaN"`...)
	case val > math.MaxFloat64:
		return append(b, `"+Inf"`...)
	case val < -math.MaxFloat64:
		return append(b, `"-Inf"`...)
	default:
		return strconv.AppendFloat(b, val, 'f', -1, 64)
	}
}

// AppendTime appends the time t formatted according to layout as a quoted JSON string to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendTime(b []byte, t time.Time, layout string) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',', '"')
	} else {
		b = append(b, '"')
	}
	b = t.AppendFormat(b, layout)
	return append(b, '"')
}

// AppendUUID appends the UUID val formatted as a quoted JSON string in the standard
// 8-4-4-4-12 hexadecimal format (e.g., "550e8400-e29b-41d4-a716-446655440000") to b.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
func AppendUUID(b []byte, val [16]byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',')
	}

	i := len(b)
	b = growLen(b, 38)

	b[i+0] = '"'
	hex.Encode(b[i+1:i+9], val[0:4])
	b[i+9] = '-'
	hex.Encode(b[i+10:i+14], val[4:6])
	b[i+14] = '-'
	hex.Encode(b[i+15:i+19], val[6:8])
	b[i+19] = '-'
	hex.Encode(b[i+20:i+24], val[8:10])
	b[i+24] = '-'
	hex.Encode(b[i+25:i+37], val[10:16])
	b[i+37] = '"'

	return b
}

// AppendJSON appends raw JSON bytes val to b without any encoding or validation.
// A comma separator is prepended if the last character of b is not ':', '[', or empty.
// The caller must ensure val contains valid JSON.
func AppendJSON(b []byte, val []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' && b[l-1] != '[' {
		b = append(b, ',')
	}
	return append(b, val...)
}

// growLen grows b by n bytes and returns the extended slice.
func growLen(b []byte, n int) []byte {
	l, c := len(b), cap(b)
	if l+n <= c {
		return b[:l+n]
	}
	newCap := l + n // TODO better growing
	newBuf := make([]byte, l+n, newCap)
	copy(newBuf, b)
	return newBuf
}
