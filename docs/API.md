# govalid v0.1.0 API reference

[Project home](../README.md) · [English guide](README.en.md) · [Guia em Português](README.pt-BR.md)

This reference lists every exported type, function, and method in the current
v0.1.0 API. Signatures use the internal constraints `integer` and `floating` as
they appear in Go documentation; callers do not need to name those constraints.

## Core types

```go
type Rule func(context ruleContext) error
```

A validation rule. The context is internal in v0.1.0, so consumers compose the
rules provided by this package.

```go
type FieldSpec struct {
	Path  string
	Rules []Rule
}
```

Associates a struct field path with rules.

```go
type Issue struct {
	Path    string
	Message string
}
```

Represents one failed rule.

```go
type FieldIssueError struct {
	RulesIssues []Issue
}
```

Contains rule failures produced by one validation run.

```go
type Option func(opts *options)
```

Configures a validator. The internal `options` type is intentionally hidden.

```go
type Structurer struct {
	// contains filtered or unexported fields
}
```

The validator returned by `New`.

## Core functions and methods

```go
func New(options ...Option) *Structurer
```

Creates a validator.

```go
func Field(path string, rules ...Rule) FieldSpec
```

Builds a field specification. Dot-separated paths address nested struct fields.

```go
func (sttr *Structurer) Validate(data any, fields ...FieldSpec) error
```

Validates selected fields from a struct or pointer to a struct.

```go
func (issue *FieldIssueError) Error() string
```

Implements `error`.

```go
func WithStopOnFirstError() Option
```

Stops validation after the first rule failure.

## Boolean rules

```go
func BoolTrue() Rule
```

Requires `true`.

```go
func BoolFalse() Rule
```

Requires `false`.

```go
func BoolEqual(expect bool) Rule
```

Requires equality with `expect`.

## String rules

### Presence and length

```go
func StringRequired() Rule
func StringLength(expect int) Rule
func StringMinLength(expect int) Rule
func StringMaxLength(expect int) Rule
```

Lengths are measured in Unicode runes.

### Character and case rules

```go
func StringAlpha() Rule
func StringAlphaNumeric() Rule
func StringNumeric() Rule
func StringASCII() Rule
func StringLowercase() Rule
func StringUppercase() Rule
func StringTrimmed() Rule
```

### Content rules

```go
func StringStartsWith(expect string) Rule
func StringStartsWithFold(expect string) Rule
func StringEndsWith(expect string) Rule
func StringEndsWithFold(expect string) Rule
func StringContains(expect string) Rule
func StringContainsFold(expect string) Rule
func StringNotContains(expect string) Rule
func StringNotContainsFold(expect string) Rule
func StringOneOf(expected ...string) Rule
func StringNotOneOf(unexpected ...string) Rule
```

`Fold` variants ignore case.

### Format rules

```go
func StringRegex(expression regexp.Regexp) Rule
func StringEmail() Rule
func StringURL() Rule
func StringUUID() Rule
```

Format rules accept an empty string; combine them with `StringRequired` when
presence is mandatory.

## Integer rules

The `integer` constraint accepts `int`, every signed and unsigned fixed-width
integer, `uint`, `uintptr`, and user-defined types based on them.

### Comparison and equality

```go
func IntMin[T integer](expect T) Rule
func IntMax[T integer](expect T) Rule
func IntGreaterThan[T integer](expect T) Rule
func IntGreaterThanOrEqual[T integer](expect T) Rule
func IntLessThan[T integer](expect T) Rule
func IntLessThanOrEqual[T integer](expect T) Rule
func IntEqual[T integer](expect T) Rule
func IntNotEqual[T integer](expect T) Rule
```

### Ranges and sets

```go
func IntBetween[Min integer, Max integer](minimum Min, maximum Max) Rule
func IntNotBetween[Min integer, Max integer](minimum Min, maximum Max) Rule
func IntOneOf[T integer](allowed ...T) Rule
func IntNoneOf[T integer](restricted ...T) Rule
```

### Sign

```go
func IntPositive() Rule
func IntNegative() Rule
func IntNonPositive() Rule
func IntNonNegative() Rule
func IntZero() Rule
func IntNonZero() Rule
```

### Mathematical rules

```go
func IntEven() Rule
func IntOdd() Rule
func IntMultipleOf[T integer](n T) Rule
func IntDivisibleBy[T integer](n T) Rule
func IntPrime() Rule
func IntCompound() Rule
func IntPowerOfTwo() Rule
func IntPerfectSquare() Rule
```

### Domain rules

```go
func IntPort() Rule
func IntPercentage() Rule
func IntHTTPStatus() Rule
```

The accepted inclusive ranges are respectively 1–65535, 0–100, and 100–599.

### Legacy integer aliases

```go
func IntNoneof[T integer](restricted ...T) Rule
func IsPerfectSquare() Rule
```

Prefer `IntNoneOf` and `IntPerfectSquare`. Legacy aliases may be removed before
v1.0.0.

## Float rules

The `floating` constraint accepts `float32`, `float64`, and user-defined types
based on them.

### Comparison and equality

```go
func FloatMin[T floating](expect T) Rule
func FloatMax[T floating](expect T) Rule
func FloatGreaterThan[T floating](expect T) Rule
func FloatGreaterThanOrEqual[T floating](expect T) Rule
func FloatLessThan[T floating](expect T) Rule
func FloatLessThanOrEqual[T floating](expect T) Rule
func FloatEqual[T floating](expect T) Rule
func FloatNotEqual[T floating](expect T) Rule
```

### Ranges and approximation

```go
func FloatBetween[Min floating, Max floating](minimum Min, maximum Max) Rule
func FloatNotBetween[Min floating, Max floating](minimum Min, maximum Max) Rule
func FloatApprox[T floating](expect T, tolerance float64) Rule
func FloatEqualWithin[T floating](expect T, tolerance float64) Rule
```

`FloatEqualWithin` is the preferred descriptive name for tolerance-based
comparison. `FloatApprox` provides the same behavior.

### Sign

```go
func FloatPositive() Rule
func FloatNegative() Rule
func FloatNonPositive() Rule
func FloatNonNegative() Rule
func FloatZero() Rule
func FloatNonZero() Rule
```

### Representation and special values

```go
func FloatBits(bits int) Rule
func FloatIs32() Rule
func FloatIs64() Rule
func FloatNaN() Rule
func FloatNotNaN() Rule
func FloatInfinite() Rule
func FloatNotInfinite() Rule
func FloatPositiveInfinite() Rule
func FloatNegativeInfinite() Rule
func FloatFinite() Rule
func FloatNotFinite() Rule
```

## Byte rules

Byte rules accept slices and arrays whose element type has `uint8` as its
underlying kind, including `[]byte`, `[N]byte`, and user-defined byte types.

### Type, nil, and emptiness

```go
func Bytes() Rule
func BytesNil() Rule
func BytesNotNil() Rule
func BytesEmpty() Rule
func BytesNotEmpty() Rule
```

Only slices can be nil. Byte arrays always pass `BytesNotNil` and fail
`BytesNil`.

### Length

```go
func BytesLength(expect int) Rule
func BytesMinLength(expect int) Rule
func BytesMaxLength(expect int) Rule
func BytesLengthBetween(minimum, maximum int) Rule
func BytesLengthNotBetween(minimum, maximum int) Rule
```

Negative lengths and reversed ranges produce validation errors.

### Equality and sets

```go
func BytesEqual(expected []byte) Rule
func BytesNotEqual(unexpected []byte) Rule
func BytesConstantTimeEqual(expected []byte) Rule
func BytesOneOf(expected ...[]byte) Rule
func BytesNoneOf(unexpected ...[]byte) Rule
```

Expected byte slices are cloned when the rule is created, so later mutations
by the caller do not alter rule behavior. `BytesConstantTimeEqual` uses
constant-time content comparison.

### Contents, prefixes, and suffixes

```go
func BytesContains(expected []byte) Rule
func BytesNotContains(unexpected []byte) Rule
func BytesContainsAny(expected ...[]byte) Rule
func BytesContainsAll(expected ...[]byte) Rule
func BytesHasPrefix(prefix []byte) Rule
func BytesNotHasPrefix(prefix []byte) Rule
func BytesHasSuffix(suffix []byte) Rule
func BytesNotHasSuffix(suffix []byte) Rule
```

Containment refers to byte subsequences rather than individual elements.

### Lexicographical comparison

```go
func BytesGreaterThan(expected []byte) Rule
func BytesGreaterThanOrEqual(expected []byte) Rule
func BytesLessThan(expected []byte) Rule
func BytesLessThanOrEqual(expected []byte) Rule
func BytesBetween(minimum, maximum []byte) Rule
```

These rules use lexicographical byte ordering.

### Byte state and text content

```go
func BytesNoZero() Rule
func BytesHasZero() Rule
func BytesAllZero() Rule
func BytesNotAllZero() Rule
func BytesUnique() Rule
func BytesASCII() Rule
func BytesPrintableASCII() Rule
func BytesUTF8() Rule
func BytesNotUTF8() Rule
```

`BytesPrintableASCII` accepts bytes from 32 through 126 inclusive.

### Encoded and document formats

```go
func BytesJSON() Rule
func BytesXML() Rule
func BytesPEM() Rule
func BytesHex() Rule
func BytesBase64() Rule
func BytesBase64URL() Rule
```

Hexadecimal and Base64 rules validate encoded textual bytes; they do not
inspect already-decoded binary data.

### Per-byte validation

```go
func BytesEach(rules ...Rule) Rule
func BytesAt(index int, rules ...Rule) Rule
func BytesAtIfPresent(index int, rules ...Rule) Rule
```

Each selected byte is supplied to nested rules as a `uint8` value, so integer
rules can be reused:

```go
govalid.BytesEach(govalid.IntBetween(uint8(0), uint8(127)))
```

## Map rules

### Type, nil, and emptiness

```go
func Map() Rule
func MapNil() Rule
func MapNotNil() Rule
func MapEmpty() Rule
func MapNotEmpty() Rule
```

### Length

```go
func MapLength(expect int) Rule
func MapMinLength(expect int) Rule
func MapMaxLength(expect int) Rule
func MapLengthBetween(minimum, maximum int) Rule
func MapLengthNotBetween(minimum, maximum int) Rule
```

### Keys

```go
func MapHasKey[K comparable](expected K) Rule
func MapNotHasKey[K comparable](unexpected K) Rule
func MapHasAllKeys[K comparable](expected ...K) Rule
func MapHasAnyKey[K comparable](expected ...K) Rule
func MapHasNoneOfKeys[K comparable](unexpected ...K) Rule
func MapAllowedKeys[K comparable](allowed ...K) Rule
func MapHasExactKeys[K comparable](expected ...K) Rule
func MapKeys(rules ...Rule) Rule
```

### Values and nested validation

```go
func MapContainsValue[V any](expected V) Rule
func MapNotContainsValue[V any](unexpected V) Rule
func MapContainsAnyValue[V any](expected ...V) Rule
func MapContainsAllValues[V any](expected ...V) Rule
func MapValues(rules ...Rule) Rule
func MapValueAt[K comparable](key K, rules ...Rule) Rule
func MapValueAtIfPresent[K comparable](key K, rules ...Rule) Rule
func MapNoNilValues() Rule
func MapNoZeroValues() Rule
func MapHasNonZeroValue() Rule
```

### Comparison and set relationships

```go
func MapEqual[M ~map[K]V, K comparable, V any](expected M) Rule
func MapNotEqual[M ~map[K]V, K comparable, V any](unexpected M) Rule
func MapSubsetOf[M ~map[K]V, K comparable, V any](expected M) Rule
func MapSupersetOf[M ~map[K]V, K comparable, V any](expected M) Rule
```

Subset and superset compare entries, including values.

### Legacy map aliases

```go
func MapIsMap() Rule
func MapHasAllKey[K comparable](expected ...K) Rule
func MapHasAnyAllKey[K comparable](expected ...K) Rule
func MapHasNoneAllKey[K comparable](unexpected ...K) Rule
```

Prefer `Map`, `MapHasAllKeys`, `MapHasAnyKey`, and `MapHasNoneOfKeys`.

## Collection rules

Collection rules accept slices and arrays.

### Type, nil, and emptiness

```go
func Collection() Rule
func CollectionNil() Rule
func CollectionNotNil() Rule
func CollectionEmpty() Rule
func CollectionNotEmpty() Rule
```

Only slices can be nil.

### Length

```go
func CollectionLength(expect int) Rule
func CollectionMinLength(expect int) Rule
func CollectionMaxLength(expect int) Rule
func CollectionLengthBetween(minimum, maximum int) Rule
func CollectionLengthNotBetween(minimum, maximum int) Rule
```

### Contents and item state

```go
func CollectionContains[T any](expected T) Rule
func CollectionNotContains[T any](unexpected T) Rule
func CollectionContainsAny[T any](expected ...T) Rule
func CollectionContainsAll[T any](expected ...T) Rule
func CollectionUnique() Rule
func CollectionNoNilItems() Rule
func CollectionNoZeroItems() Rule
func CollectionHasNonZeroItem() Rule
```

### Nested item validation

```go
func CollectionEach(rules ...Rule) Rule
func CollectionAll(rules ...Rule) Rule
func CollectionAny(rules ...Rule) Rule
func CollectionNone(rules ...Rule) Rule
func CollectionItemAt(index int, rules ...Rule) Rule
func CollectionItemAtIfPresent(index int, rules ...Rule) Rule
```

`CollectionAll` is an alias for `CollectionEach`. Every item must satisfy every
provided rule.

### Comparison

```go
func CollectionEqual[T any](expected T) Rule
func CollectionNotEqual[T any](unexpected T) Rule
```

## API stability

This reference describes the pre-v1 v0.1.0 surface. Canonical names should be
preferred over legacy aliases. Breaking cleanup remains possible until v1.0.0.
