package logrus

import (
	"github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/plugin/logger/logrus/logrus"
)

type Options struct {
	logger.Options
	Formatter logrus.Formatter
	// Flag for whether to log caller info (off by default)
	ReportCaller    bool
	SplitLevel      bool
	WithoutKey      bool
	WithoutQuote    bool
	TimestampFormat string
	// Exit Function to call when FatalLevel log
	ExitFunc func(int)
}

type formatterKey struct{}
type splitLevelKey struct{}
type reportCallerKey struct{}
type exitKey struct{}
type withoutKeyKey struct{}
type withoutQuoteKey struct{}
type timestampFormat struct{}

func TextFormatter(formatter *logrus.TextFormatter) logger.Option {
	return logger.SetOption(formatterKey{}, formatter)
}

func JSONFormatter(formatter *logrus.JSONFormatter) logger.Option {
	return logger.SetOption(formatterKey{}, formatter)
}

func ExitFunc(exit func(int)) logger.Option {
	return logger.SetOption(exitKey{}, exit)
}

func SplitLevel(s bool) logger.Option {
	return logger.SetOption(splitLevelKey{}, s)
}

func WithoutKey(w bool) logger.Option {
	return logger.SetOption(withoutKeyKey{}, w)
}

func WithoutQuote(w bool) logger.Option {
	return logger.SetOption(withoutQuoteKey{}, w)
}

func TimestampFormat(format string) logger.Option {
	return logger.SetOption(timestampFormat{}, format)
}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func ReportCaller(r bool) logger.Option {
	return logger.SetOption(reportCallerKey{}, r)
}
