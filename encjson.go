package encjson

import (
	"encoding/hex"
	"math"
	"strconv"
	"time"
)

func AppendObjectStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ',', '{')
	}
	return append(b, '{')
}

func AppendObjectEnd(b []byte) []byte {
	return append(b, '}')
}

func AppendArrayStart(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ',', '[')
	}
	return append(b, '[')
}

func AppendArrayEnd(b []byte) []byte {
	return append(b, ']')
}

func AppendKey(b []byte, key string) []byte {
	return append(AppendString(b, key), ':')
}

func AppendBool(b []byte, val bool) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
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

func AppendNull(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ",null"...)
	}
	return append(b, "null"...)
}

func AppendInt(b []byte, val int64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		b = append(b, ',')
	}
	return strconv.AppendInt(b, val, 10)
}

func AppendUint(b []byte, val uint64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		b = append(b, ',')
	}
	return strconv.AppendUint(b, val, 10)
}

func AppendFloat(b []byte, val float64) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
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

func AppendTime(b []byte, t time.Time, layout string) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		b = append(b, ',')
	}
	return t.AppendFormat(b, layout)
}

func AppendUUID(b []byte, val [16]byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		b = append(b, ',')
	}

	// TODO use grow len first then write directly to slice
	var a [38]byte
	a[0] = '"'
	hex.Encode(a[1:9], val[0:4])
	a[9] = '-'
	hex.Encode(a[10:14], val[4:6])
	a[14] = '-'
	hex.Encode(a[15:19], val[6:8])
	a[19] = '-'
	hex.Encode(a[20:24], val[8:10])
	a[2] = '-'
	hex.Encode(a[25:37], val[10:16])
	a[37] = '"'

	return append(b, a[:]...)
}
