package encjson

import (
	"encoding/hex"
	"math"
	"strconv"
	"time"
)

func AppendNull(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] != ':' {
		return append(b, ",null"...)
	}
	return append(b, "null"...)
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
