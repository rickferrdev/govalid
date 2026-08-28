<p align="center">
  <img src="assets/govalid-gopher-v1.png" alt="govalid v1 gopher mascot" width="270">
</p>

<h1 align="center">🐳 govalid</h1>

<p align="center">
  <strong>Composable, type-aware validation for Go structs.</strong><br>
  Explicit by default. Extensible when you need it.
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/rickferrdev/govalid"><img src="https://pkg.go.dev/badge/github.com/rickferrdev/govalid.svg" alt="Go Reference"></a>
  <a href="https://github.com/rickferrdev/govalid/releases"><img src="https://img.shields.io/github/v/release/rickferrdev/govalid?include_prereleases&sort=semver" alt="GitHub release"></a>
  <img src="https://img.shields.io/badge/go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-green.svg" alt="MIT license"></a>
</p>

<p align="center">
  <a href="#quick-start">Quick start</a> ·
  <a href="#rules">Rules</a> ·
  <a href="#errors-you-can-use">Errors</a> ·
  <a href="#built-to-extend">Extensibility</a> ·
  <a href="docs/API.md">API reference</a>
</p>

---

`govalid` validates the fields you choose with ordinary Go functions. Rules
stay visible in code, compose across data types, and work without a tag
language, global registry, or code generation.

```go
err := govalid.New().Validate(
	user,
	govalid.Field("Name", govalid.Required(), govalid.StringMinLength(3)),
	govalid.Field("Age", govalid.IntBetween(18, 130)),
	govalid.Field("Profile.Email", govalid.StringEmail()),
)
```

## Why teams choose govalid

| | |
| --- | --- |
| 🧩 **Composable** | Reuse rules inside maps, collections, bytes, and conditional flows. |
| 🔎 **Explicit** | Validation is easy to find, review, refactor, and test. |
| 🧠 **Type-aware** | Signed and unsigned integers, derived types, floats, bytes, maps, slices, and arrays are handled deliberately. |
| 🧭 **Nested** | Address fields with paths such as `"Profile.Email"`; use `FieldIfPresent` for nullable parents. |
| 🧰 **Extensible** | Build custom `Rule` functions and custom `FieldSource` discovery mechanisms. |
| ⚡ **Reusable** | A configured `Validator` can be shared across concurrent calls when callbacks and sources are concurrency-safe. |

## Install

```bash
go get github.com/rickferrdev/govalid@latest
```

The project is finalizing its `v1.0.0` contract. Until the stable tag is
published, prereleases may still include breaking cleanup.

## Quick start

```go
package main

import (
	"errors"
	"fmt"

	"github.com/rickferrdev/govalid"
)

type Profile struct {
	Email string
}

type User struct {
	Name    string
	Age     uint
	Profile *Profile
	Scores  []int
	Labels  map[string]string
	Payload []byte
}

func main() {
	user := User{
		Name:    "Alice",
		Age:     24,
		Profile: &Profile{Email: "alice@example.com"},
		Scores:  []int{98, 91, 87},
		Labels:  map[string]string{"environment": "production"},
		Payload: []byte(`{"active":true}`),
	}

	err := govalid.New().Validate(
		user,
		govalid.Field("Name", govalid.Required(), govalid.StringMinLength(3)),
		govalid.Field("Age", govalid.IntBetween(18, 130)),
		govalid.FieldIfPresent("Profile.Email", govalid.StringEmail()),
		govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
		govalid.Field("Labels", govalid.MapKeys(govalid.StringLowercase())),
		govalid.Field("Payload", govalid.BytesJSON()),
	)

	var validationErr *govalid.FieldIssueError
	if errors.As(err, &validationErr) {
		for _, issue := range validationErr.RulesIssues {
			fmt.Printf("%s: %s (rule: %s)\n", issue.Path, issue.Message, issue.Rule)
		}
	}
}
```

## Rules

| Category | What you can validate |
| --- | --- |
| **String** | Unicode length, content, case, regex, email, URL, and UUID |
| **Integer** | Every signed/unsigned variant, comparisons, sets, arithmetic, and domains |
| **Float** | Comparisons, tolerance, bit width, `NaN`, infinity, and finiteness |
| **Boolean** | True, false, and equality |
| **Bytes** | Length, equality, content, UTF-8, JSON, XML, PEM, Hex, and Base64 |
| **Map** | Keys, values, length, equality, subsets, and nested rules |
| **Collection** | Content, uniqueness, length, equality, and nested item rules |
| **Universal** | Required, nil, non-nil, zero, and non-zero values |
| **Conditional** | `When`, `Unless`, context-aware conditions, and `Optional` |

Integer comparisons preserve the complete `uint64` range even when signed and
unsigned values are mixed. Defined types with supported underlying kinds work
as well.

### Compose instead of repeating

```go
govalid.Field(
	"Labels",
	govalid.Required(),
	govalid.MapKeys(govalid.StringLowercase()),
	govalid.MapValues(govalid.StringRequired()),
)

govalid.Field(
	"Scores",
	govalid.Optional(
		govalid.CollectionEach(govalid.IntBetween(0, 100)),
	),
)
```

Context-aware conditions can inspect the selected value and root struct:

```go
govalid.WhenContext(func(ctx govalid.RuleContext) bool {
	return ctx.RootAny().(User).Age >= 18
}, govalid.StringRequired())
```

## Errors you can use

Validation collects every issue by default. `errors.As` can extract the full
`*FieldIssueError` or its first `*Issue`.

```go
type Issue struct {
	Path    string
	Value   any
	Rule    string
	Message string
}
```

Choose the execution flow that fits the boundary you are validating:

```go
govalid.New()                                // collect every issue
govalid.New(govalid.WithStopOnFirstError())  // return after the first issue
govalid.New(govalid.WithPanicOnFirstError()) // panic with *govalid.Issue
govalid.New(govalid.WithSilenceErrors())     // run rules and return nil
```

Observe failures independently of the return mode:

```go
validator := govalid.New(
	govalid.WithIssueHandler(func(issue govalid.Issue) {
		log.Printf("validation %s: %s", issue.Path, issue.Message)
	}),
)
```

Large values remain available in `Issue.Value`, while their formatted error
representation is summarized to keep logs manageable.

## Built to extend

Custom rules use the same context as built-in rules:

```go
isSlug := govalid.Rule(func(ctx govalid.RuleContext) error {
	value, ok := ctx.ValueAny().(string)
	if !ok || value == "" {
		return errors.New("should be a slug")
	}
	return nil
})
```

`FieldSource` separates field discovery from rule execution. A future tag,
schema, or generated integration can return ordinary `FieldSpec` values and
reuse the same validator:

```go
type FieldSource interface {
	Fields(structType reflect.Type) ([]govalid.FieldSpec, error)
}
```

Sources are additive and run before fields passed directly to `Validate`.
Shared sources and issue handlers must be concurrency-safe.

## v1 scope

The first stable release focuses on explicit field validation. Built-in struct
tags and recursive `Struct*` rules are intentionally outside the initial v1
scope. The `FieldSource` boundary keeps those integrations possible later
without replacing `Validator`, `Rule`, or `FieldSpec`.

The suite includes external API tests, concurrent reuse tests, and large
payload tests with 100,000 collection items, 10,000 map entries, nested
matrices, and megabyte-sized byte documents.

## Documentation

- [English overview](docs/README.en.md)
- [Português (Brasil)](docs/README.pt-BR.md)
- [Complete API reference](docs/API.md) — every exported signature, in English
- [Package documentation](https://pkg.go.dev/github.com/rickferrdev/govalid)

## Development

```bash
go test ./...
go vet ./...
```

## License

Released under the [MIT License](LICENSE).
