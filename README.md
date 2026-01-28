# go-encjson

Performance optimized JSON encoding in Go.

This package provides low-level functions for building JSON output by appending to byte slices. It's designed for high-performance scenarios where you need fine-grained control over JSON generation without the overhead of reflection-based encoding.

## Installation

```bash
go get github.com/domonda/go-encjson
```

## Features

- Zero-allocation appending to existing byte slices
- Automatic comma separator handling
- Proper JSON string escaping (including U+2028 and U+2029)
- Support for all JSON value types

## Usage

All `Append*` functions automatically prepend a comma separator when needed based on the context (the last character of the buffer).

### Basic Types

```go
import "github.com/domonda/go-encjson"

var b []byte

// Null
b = encjson.AppendNull(b)  // null

// Boolean
b = encjson.AppendBool(nil, true)   // true
b = encjson.AppendBool(nil, false)  // false

// Numbers
b = encjson.AppendInt(nil, -42)     // -42
b = encjson.AppendUint(nil, 42)     // 42
b = encjson.AppendFloat(nil, 3.14)  // 3.14

// Strings (properly escaped)
b = encjson.AppendString(nil, "hello")           // "hello"
b = encjson.AppendString(nil, "line\nbreak")     // "line\nbreak"
b = encjson.AppendStringBytes(nil, []byte("hi")) // "hi"

// Time
b = encjson.AppendTime(nil, time.Now(), time.RFC3339) // "2024-01-15T10:30:00Z"

// UUID (as [16]byte)
var uuid [16]byte
b = encjson.AppendUUID(nil, uuid) // "00000000-0000-0000-0000-000000000000"

// Raw JSON (no encoding or validation)
b = encjson.AppendJSON(nil, []byte(`{"nested":true}`)) // {"nested":true}
```

### Arrays

```go
b := encjson.AppendArrayStart(nil)  // [
b = encjson.AppendInt(b, 1)         // [1
b = encjson.AppendInt(b, 2)         // [1,2
b = encjson.AppendInt(b, 3)         // [1,2,3
b = encjson.AppendArrayEnd(b)       // [1,2,3]
```

### Objects

```go
b := encjson.AppendObjectStart(nil)      // {
b = encjson.AppendKey(b, "name")         // {"name":
b = encjson.AppendString(b, "John")      // {"name":"John"
b = encjson.AppendKey(b, "age")          // {"name":"John","age":
b = encjson.AppendInt(b, 30)             // {"name":"John","age":30
b = encjson.AppendObjectEnd(b)           // {"name":"John","age":30}
```

For keys that are known to be safe ASCII (no escaping needed), use the faster `AppendSafeKey`:

```go
b := encjson.AppendObjectStart(nil)
b = encjson.AppendSafeKey(b, "count")  // Faster for known-safe keys
b = encjson.AppendInt(b, 42)
b = encjson.AppendObjectEnd(b)         // {"count":42}
```

### Nested Structures

```go
// {"users":[{"name":"Alice"},{"name":"Bob"}]}
b := encjson.AppendObjectStart(nil)
b = encjson.AppendSafeKey(b, "users")
b = encjson.AppendArrayStart(b)

b = encjson.AppendObjectStart(b)
b = encjson.AppendSafeKey(b, "name")
b = encjson.AppendString(b, "Alice")
b = encjson.AppendObjectEnd(b)

b = encjson.AppendObjectStart(b)
b = encjson.AppendSafeKey(b, "name")
b = encjson.AppendString(b, "Bob")
b = encjson.AppendObjectEnd(b)

b = encjson.AppendArrayEnd(b)
b = encjson.AppendObjectEnd(b)
```

### String Escaping Check

Use `StringNeedsEscaping` to check if a string contains characters that require escaping:

```go
encjson.StringNeedsEscaping("hello")      // false
encjson.StringNeedsEscaping("hello\n")    // true
encjson.StringNeedsEscaping("quote\"")    // true
```

## Special Float Values

Special float values that are not valid in JSON are encoded as quoted strings:

```go
encjson.AppendFloat(nil, math.NaN())     // "NaN"
encjson.AppendFloat(nil, math.Inf(1))    // "+Inf"
encjson.AppendFloat(nil, math.Inf(-1))   // "-Inf"
```

## License

MIT
