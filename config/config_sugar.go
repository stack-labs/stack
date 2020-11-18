package config

import (
	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
)

var (
	sugar config.Config
)

func Get(path ...string) reader.Value {
	return sugar.Get(path...)
}

func GetDefault() *Value {
	return &Value{
		Stack: Stack{
			Broker: Broker{
				Name: "http",
			},
			Client: Client{
				Protocol: "mucp",
				Pool: Pool{
					Size: "1",
					TTL:  "60s",
				},
				Request: ClientRequest{
					Retries: "1",
					Timeout: "5s",
				},
			},
			Profile: "",
			Registry: Registry{
				Name: "mdns",
			},
			Runtime: "",
			Server: Server{
				Protocol: "mucp",
			},
			Selector: Selector{
				Name: "registry",
			},
			Transport: Transport{
				Name: "http",
			},
		},
	}
}
