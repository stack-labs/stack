package hook

import (
	"github.com/stack-labs/stack"
	"github.com/stack-labs/stack/plugin/service/stackway/api"
	"github.com/stack-labs/stack/plugin/service/stackway/plugin"
	"github.com/stack-labs/stack/service"
	"github.com/stack-labs/stack/util/log"
)

func Hook(svc service.Service) {
	apiServer := api.NewServer(svc)

	// stackway options
	_ = svc.Init(api.Options()...)

	// stackway hook
	_ = svc.Init(
		stack.AfterStart(apiServer.Start),
		stack.AfterStop(apiServer.Stop),
	)

	// plugin tags
	plugins := plugin.Plugins()
	for _, p := range plugins {
		log.Debugf("plugin: %s", p.String())
		if flags := p.Flags(); len(flags) > 0 {
			log.Debugf("flags: %+#s", flags)
			_ = svc.Init(stack.Flags(flags...))
		}
	}

	return
}
