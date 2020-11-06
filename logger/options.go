package logger

import (
	"context"
	"io"
)

type Option func(*Options)

type Options struct {
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields map[string]interface{}
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// Caller skip frame count for file:line info
	CallerSkipCount int
	Persistence     bool
	// Alternative options
	Context context.Context
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(args *Options) {
		args.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(args *Options) {
		args.Level = level
	}
}

// Output set default output writer for the logger
func Output(out io.Writer) Option {
	return func(args *Options) {
		args.Out = out
	}
}

func Persistence(p bool) Option {
	return func(options *Options) {
		options.Persistence = p
	}
}

// CallerSkipCount set frame count to skip
func CallerSkipCount(c int) Option {
	return func(args *Options) {
		args.CallerSkipCount = c
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
