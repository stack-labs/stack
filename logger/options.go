package logger

import (
	"context"
	"io"
)

type PersistenceOptions struct {
	Enable    bool   `sc:"enable"`
	Dir       string `sc:"dir"`
	BackupDir string `sc:"back-dir"`
	// log file max size in megabytes
	MaxFileSize int `sc:"max-file-size"`
	// backup dir max size in megabytes
	MaxBackupSize int `sc:"max-backup-size"`
	// backup files keep max days
	MaxBackupKeepDays int `sc:"max-backup-keep-days"`
	// default pattern is ${serviceName}_${level}.log
	// todo available patterns map
	FileNamePattern string `sc:"file-name-pattern"`
	// default pattern is ${serviceName}_${level}_${yyyyMMdd_HH}_${idx}.zip
	// todo available patterns map
	BackupFileNamePattern string `sc:"backup-file-name-pattern"`
}

type Option func(*Options)

type Options struct {
	// logger's name, same to logger.String(). console, logrus, zap etc.
	Name string
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields map[string]interface{}
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// Alternative options
	Context context.Context
	// Caller skip frame count for file:line info
	CallerSkipCount int
	Persistence     *PersistenceOptions
}

func Name(n string) Option {
	return func(options *Options) {
		options.Name = n
	}
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(options *Options) {
		options.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(options *Options) {
		options.Level = level
	}
}

// Output set default output writer for the logger
func Output(out io.Writer) Option {
	return func(options *Options) {
		options.Out = out
	}
}

func Persistence(o *PersistenceOptions) Option {
	return func(options *Options) {
		options.Persistence = o
	}
}

// CallerSkipCount set frame count to skip
func CallerSkipCount(c int) Option {
	return func(options *Options) {
		options.CallerSkipCount = c
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
