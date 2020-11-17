package stack

import "github.com/stack-labs/stack-rpc/config"

func newDefaultConfig() *config.Value {
	return &config.Value{
		Stack: config.Stack{
			Broker: config.Broker{
				Name: "http",
			},
			Client: config.Client{
				Protocol: "mucp",
				Pool: config.Pool{
					Size: "1",
					TTL:  "60s",
				},
				Request: config.ClientRequest{
					Retries: "1",
					Timeout: "5s",
				},
			},
			Profile: "",
			Registry: config.Registry{
				Name: "mdns",
			},
			Runtime: "",
			Server: config.Server{
				Protocol: "mucp",
			},
			Selector: config.Selector{
				Name: "registry",
			},
			Transport: config.Transport{
				Name: "http",
			},
		},
	}
}
