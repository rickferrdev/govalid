# govalid

[Project home](../README.md) · [Português (Brasil)](README.pt-BR.md) · [Complete API reference](API.md)

`govalid` provides composable rules for validating selected Go struct fields
without tags. The v0.1.0 scope includes booleans, strings, integers, floats,
maps, collections, and byte sequences.

> This is a pre-v1 project. Public names and behavior may still change before
> API stabilization.

## Installation

```bash
go get github.com/rickferrdev/govalid@v0.1.0
```

Use the module path without a version to install the current development
revision.

## Basic usage

```go
validator := govalid.New()

err := validator.Validate(
	user,
	govalid.Field("Name", govalid.StringRequired(), govalid.StringMinLength(3)),
	govalid.Field("Age", govalid.IntBetween(18, 130)),
	govalid.Field("Active", govalid.BoolTrue()),
	govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
)
```

Nested struct fields can be addressed with paths such as `"Profile.Email"`.
By default, all rule failures are collected. Use
`New(WithStopOnFirstError())` to stop at the first issue.

```go
var validationErr *govalid.FieldIssueError

if errors.As(err, &validationErr) {
	for _, issue := range validationErr.RulesIssues {
		fmt.Printf("%s: %s\n", issue.Path, issue.Message)
	}
}
```

## Rule summary

### Boolean

`BoolTrue`, `BoolFalse`, and `BoolEqual` validate boolean state and equality.

### String

String rules cover:

- required values and Unicode-aware length;
- alphabetic, numeric, alphanumeric, and ASCII content;
- lowercase, uppercase, and trimmed text;
- prefixes, suffixes, containment, and allowed sets;
- regular expressions, email, absolute URL, and UUID formats.

Format-only rules accept an empty string. Combine them with `StringRequired`
when the field is mandatory.

### Integer

Integer rules accept all signed and unsigned integer variants, `uintptr`, and
user-defined integer types. They cover:

- comparison, equality, inclusive ranges, and sets;
- positive, negative, and zero states;
- parity, divisibility, primes, composite values, powers of two, and perfect
  squares;
- ports, percentages, and HTTP status ranges.

Mixed signed/unsigned comparisons preserve the complete `uint64` range.

### Float

Float rules support `float32`, `float64`, and derived types. They cover:

- comparison, equality, ranges, and tolerance;
- sign and zero states;
- 32-bit and 64-bit representation;
- `NaN`, positive/negative infinity, and finiteness.

### Map

Map rules cover:

- nil, empty, and length checks;
- required, forbidden, allowed, and exact keys;
- value containment and zero/nil values;
- rules applied to every key, every value, or a selected key;
- equality, subset, and superset relationships.

```go
govalid.Field(
	"Labels",
	govalid.MapNotEmpty(),
	govalid.MapKeys(govalid.StringLowercase()),
	govalid.MapValues(govalid.StringRequired()),
)
```

### Collection

Collection rules support slices and arrays. They cover:

- nil, empty, and length checks;
- item containment, uniqueness, zero, and nil values;
- equality and inequality;
- rules applied to every item, any item, no item, or a selected index.

Arrays are never nil; only slices can satisfy `CollectionNil`.

### Bytes

Byte rules support byte slices, byte arrays, and user-defined byte types. They
cover:

- nil, empty, and length checks;
- equality, constant-time equality, containment, prefixes, and suffixes;
- lexicographical comparisons and ranges;
- zero bytes, uniqueness, ASCII, printable ASCII, and UTF-8;
- JSON, XML, PEM, hexadecimal, Base64, and Base64URL content;
- rules applied to every byte or a selected index.

## v0.1.0 limitations

- No recursive `Struct*` rules yet; nested field paths are supported.
- No dedicated rules for time, duration, pointers, or universal presence yet.
- Map iteration order is unspecified, so the first nested map failure may vary.

For every method signature and legacy alias, see the
[complete API reference](API.md).
