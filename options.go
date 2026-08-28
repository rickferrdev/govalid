package govalid

type options struct {
	StopOnFirstError      bool
	WithSilenceErrors     bool
	WithPanicOnFirstError bool
	IssueHandler          func(Issue)
	FieldSources          []FieldSource
}

// Option configures a Validator during construction.
type Option func(opts *options)

func defaultOptionsMap() options {
	return options{
		StopOnFirstError:      false,
		WithSilenceErrors:     false,
		WithPanicOnFirstError: false,
		IssueHandler:          nil,
		FieldSources:          nil,
	}
}

// WithStopOnFirstError makes validation return after the first failed rule.
func WithStopOnFirstError() Option {
	return func(opts *options) {
		opts.StopOnFirstError = true
	}
}

// WithPanicOnFirstError makes validation panic with an Issue on the first
// failed rule.
func WithPanicOnFirstError() Option {
	return func(opts *options) {
		opts.WithPanicOnFirstError = true
	}
}

// WithSilenceErrors makes validation return nil after running applicable
// rules. Combine it with WithIssueHandler to observe failures.
func WithSilenceErrors() Option {
	return func(opts *options) {
		opts.WithSilenceErrors = true
	}
}

// WithIssueHandler registers a callback invoked once for every failed rule.
// The handler runs before stop, panic, and silence behavior is applied. A
// handler used by concurrent validations must itself be concurrency-safe.
func WithIssueHandler(handler func(Issue)) Option {
	return func(opts *options) {
		opts.IssueHandler = handler
	}
}

// WithFieldSources registers field-discovery extensions that run before the
// explicit FieldSpec values passed to Validate. Sources are copied during
// construction and must be concurrency-safe when the Validator is shared.
func WithFieldSources(sources ...FieldSource) Option {
	copied := append([]FieldSource(nil), sources...)
	return func(opts *options) {
		opts.FieldSources = append(opts.FieldSources, copied...)
	}
}
