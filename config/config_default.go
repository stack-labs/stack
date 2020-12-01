package config

var (
	DefaultConfig = NewConfig()

	// todo this doest work
	_defaultCfg = &Value{
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
)

func GetDefault() *Value {
	return _defaultCfg
}
