# govalid · English overview

[← Project home](../README.md) · [Português](README.pt-BR.md) · [API reference](API.md)

> Explicit, composable validation for selected fields of Go structs.

## Start here

```bash
go get github.com/rickferrdev/govalid@latest
```

```go
validator := govalid.New()

err := validator.Validate(
	user,
	govalid.Field("Name", govalid.Required(), govalid.StringMinLength(3)),
	govalid.Field("Age", govalid.IntBetween(18, 130)),
	govalid.FieldIfPresent("Profile.Email", govalid.StringEmail()),
	govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
)
```

The model is intentionally small:

1. Create or reuse a `Validator`.
2. Associate paths with rules through `Field` or `FieldIfPresent`.
3. Inspect the returned structured issues.

## What is included

| Family | Coverage |
| --- | --- |
| Boolean | State and equality |
| String | Unicode length, content, case, regex, email, URL, UUID |
| Integer | Signed/unsigned comparison, sets, arithmetic, domain ranges |
| Float | Comparison, tolerance, representation, special values |
| Bytes | Binary content, encodings, document formats, per-byte rules |
| Map | Keys, values, length, relationships, nested rules |
| Collection | Content, length, uniqueness, comparison, nested item rules |
| Universal | Presence, nil, and zero-value rules |
| Conditional | Boolean conditions, context conditions, optional values |

Defined types with supported underlying kinds are accepted. Integer
comparisons preserve the complete `uint64` range when signed and unsigned
values are mixed.

## Composition patterns

Validate every map key and value:

```go
govalid.Field(
	"Labels",
	govalid.MapKeys(govalid.StringLowercase()),
	govalid.MapValues(govalid.StringRequired()),
)
```

Skip nested rules when a value is absent:

```go
govalid.Field(
	"Scores",
	govalid.Optional(
		govalid.CollectionEach(govalid.IntBetween(0, 100)),
	),
)
```

Apply rules based on the root struct:

```go
govalid.WhenContext(func(ctx govalid.RuleContext) bool {
	return ctx.RootAny().(Account).Enabled
}, govalid.StringRequired())
```

## Structured errors

By default, all failures are returned in `*FieldIssueError`:

```go
var validationErr *govalid.FieldIssueError

if errors.As(err, &validationErr) {
	for _, issue := range validationErr.RulesIssues {
		fmt.Printf("%s: %s\n", issue.Path, issue.Message)
	}
}
```

Each `Issue` exposes the field `Path`, rejected `Value`, `Rule`, and `Message`.
Use `WithStopOnFirstError`, `WithPanicOnFirstError`, or `WithSilenceErrors` to
change execution. `WithIssueHandler` observes failures independently of the
selected return behavior.

## Extensibility

- Implement custom rules with `Rule` and `RuleContext`.
- Use `FieldIfPresent` when a nested pointer or interface may be nil.
- Implement `FieldSource` to discover `FieldSpec` values from another format.
- Reuse a configured validator concurrently when callbacks and sources are
  concurrency-safe.

`FieldSource` is the compatibility boundary for future tag, schema, or
generated integrations. No built-in tag syntax is part of the initial v1
scope.

## v1 scope

The initial stable API focuses on explicit rules for primitive values, maps,
collections, bytes, nested paths, conditional composition, and structured
errors. Dedicated recursive `Struct*` rules, time/duration families, and a
built-in tag parser can evolve in later releases.

Map iteration order is unspecified, so the first nested map failure may vary.
Format-only string rules accept an empty value; combine them with `Required`
or `StringRequired` when presence is mandatory.

## Next links

- [Complete API reference](API.md)
- [Root README and full example](../README.md)
- [Package documentation](https://pkg.go.dev/github.com/rickferrdev/govalid)
