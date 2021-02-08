package logrus

import (
	"github.com/stack-labs/stack/config"
	"github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/plugin"
	"github.com/stack-labs/stack/plugin/logger/logrus/logrus"
	scfg "github.com/stack-labs/stack/service/config"
)

var options struct {
	Stack struct {
		Logger struct {
			scfg.Logger
			Logrus struct {
				SplitLevel      bool   `sc:"split-level"`
				ReportCaller    bool   `sc:"report-caller"`
				Formatter       string `sc:"formatter"`
				WithoutKey      bool   `sc:"without-key"`
				WithoutQuote    bool   `sc:"without-quote"`
				TimestampFormat string `sc:"timestamp-format"`
			} `sc:"slogrus"`
		} `sc:"logger"`
	} `sc:"stack"`
}

type logrusLogPlugin struct{}

func (l *logrusLogPlugin) Name() string {
	return "slogrus"
}

func (l *logrusLogPlugin) Options() []logger.Option {
	var opts []logger.Option
	lc := options.Stack.Logger.Logrus
	opts = append(opts, SplitLevel(lc.SplitLevel))
	opts = append(opts, ReportCaller(lc.ReportCaller))
	opts = append(opts, WithoutKey(lc.WithoutKey))
	opts = append(opts, WithoutQuote(lc.WithoutQuote))

	if len(lc.TimestampFormat) > 0 {
		opts = append(opts, TimestampFormat(lc.TimestampFormat))
	}

	switch lc.Formatter {
	case "text":
		opts = append(opts, TextFormatter(new(logrus.TextFormatter)))
	case "json":
		opts = append(opts, JSONFormatter(new(logrus.JSONFormatter)))
	}

	return opts
}

func (l *logrusLogPlugin) New(opts ...logger.Option) logger.Logger {
	return NewLogger(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.LoggerPlugins["slogrus"] = &logrusLogPlugin{}
}
