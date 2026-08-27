<p align="center">
  <img src="assets/govalid-gopher.png" alt="govalid gopher mascot holding a validation checklist" width="220">
</p>

# govalid

Composable validation rules for Go structs.

`govalid` validates selected struct fields without tags. Rules are ordinary
values that can be combined per field, including nested struct paths, map keys
and values, and collection items.

> Current target: **v0.1.0**. The project is pre-v1 and its public API may still
> change. Dedicated nested-struct rules are not included yet.

## Documentation

Choose a language for the project overview, or open the complete API catalog:

- [English guide](docs/README.en.md)
- [Guia em Português (Brasil)](docs/README.pt-BR.md)
- [Complete API reference](docs/API.md) — English only, with every public signature

## Install

After the v0.1.0 tag is published:

```bash
go get github.com/rickferrdev/govalid@v0.1.0
```

For the current development version:

```bash
go get github.com/rickferrdev/govalid
```

## Quick start

```go
package main

import (
	"errors"
	"fmt"

	"github.com/rickferrdev/govalid"
)

type User struct {
	Name   string
	Age    uint
	Active bool
	Scores []int
	Labels map[string]string
}

func main() {
	user := User{
		Name:   "Alice",
		Age:    24,
		Active: true,
		Scores: []int{90, 85},
		Labels: map[string]string{"environment": "production"},
	}

	err := govalid.New().Validate(
		user,
		govalid.Field("Name", govalid.StringRequired(), govalid.StringMinLength(3)),
		govalid.Field("Age", govalid.IntBetween(18, 130)),
		govalid.Field("Active", govalid.BoolTrue()),
		govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
		govalid.Field("Labels", govalid.MapKeys(govalid.StringLowercase())),
	)

	if err == nil {
		return
	}

	var validationErr *govalid.FieldIssueError
	if errors.As(err, &validationErr) {
		for _, issue := range validationErr.RulesIssues {
			fmt.Printf("%s: %s\n", issue.Path, issue.Message)
		}
	}
}
```

## Included in v0.1.0

- Boolean rules: true, false, and equality.
- String rules: presence, Unicode-aware length, content, case, regex, email,
  URL, and UUID.
- Integer rules: all signed and unsigned integer variants, comparisons,
  ranges, sets, divisibility, and mathematical properties.
- Float rules: `float32` and `float64`, comparisons, tolerance, bit width,
  `NaN`, infinity, and finiteness.
- Map rules: nil/empty, length, keys, values, nested rules, and comparisons.
- Collection rules: slices and arrays, length, contents, uniqueness, nested
  rules, and comparisons.
- Byte rules: byte slices and arrays, length, binary content, encodings,
  document formats, and per-byte validation.
- Nested field lookup using paths such as `"Profile.Email"`.
- Collection of every issue, or early exit with `WithStopOnFirstError()`.

## Scope notes

- Empty strings are accepted by format-only rules such as `StringEmail()`;
  combine them with `StringRequired()` when presence is required.
- `Collection` rules support slices and arrays. Only slices can be nil.
- Integer rules accept `int`, all `intN`/`uintN` variants, `uint`, `uintptr`,
  and user-defined types with those underlying types.
- Dedicated rules for arbitrary nested struct validation, time, duration, and
  universal presence are planned beyond the initial scope.

## Development

```bash
go test ./...
go vet ./...
```

See the language-specific guides for the complete rule catalog and behavioral
details.
