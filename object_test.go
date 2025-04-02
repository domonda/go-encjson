package encjson

import (
	"testing"
)

func TestAppendObjectStart(t *testing.T) {
	var b []byte
	b = AppendObjectStart(b)

	// Empty object
	b = AppendKey(b, "emptyObject")
	b = AppendObjectStart(b)
	b = AppendObjectEnd(b)

	b = AppendKey(b, "mixedObjectArray")
	b = AppendArrayStart(b)
	b = AppendObjectStart(b)
	b = AppendObjectEnd(b)
	b = AppendObjectStart(b)
	b = AppendObjectEnd(b)
	b = AppendString(b, "String")
	b = AppendObjectStart(b)
	b = AppendKey(b, "val1")
	b = AppendInt(b, 1)
	b = AppendObjectEnd(b)
	b = AppendArrayEnd(b)

	b = AppendKey(b, "val2")
	b = AppendFloat(b, 2)

	b = AppendObjectEnd(b)

	expected := `{"emptyObject":{},"mixedObjectArray":[{},{},"String",{"val1":1}],"val2":2}`
	if string(b) != expected {
		t.Fatalf("expected: %s, got: %s", expected, string(b))
	}
}
