package govalid

type options struct {
	StopOnFirstError bool
}

type Option func(opts *options)

func defaultOptionsMap() options {
	return options{
		StopOnFirstError: false,
	}
}

func WithStopOnFirstError() Option {
	return func(opts *options) {
		opts.StopOnFirstError = true
	}
}
