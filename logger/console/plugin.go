package console

import (
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugin"
)

type consoleLogPlugin struct {
}

func (q *consoleLogPlugin) Name() string {
	return "console"
}

func (q *consoleLogPlugin) Options() []logger.Option {
	return nil
}

func (q *consoleLogPlugin) New(opts ...logger.Option) logger.Logger {
	return logger.NewLogger(opts...)
}

func init() {
	plugin.LoggerPlugins["console"] = &consoleLogPlugin{}
}
