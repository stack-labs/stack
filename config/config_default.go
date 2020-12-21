package config

var (
	DefaultConfig = NewConfig()

	// todo this doest work
	_defaultCfg = &Value{
		Stack: Stack{
			Broker: Broker{
				Name: "",
			},
			Client: Client{
				Protocol: "",
				Pool: Pool{
					Size: 10,
					TTL:  10,
				},
				Request: ClientRequest{
					Retries: 1,
					Timeout: "2s",
				},
			},
			Profile: "",
			Registry: Registry{
				Name: "",
			},
			Runtime: "",
			Server: Server{
				Protocol: "",
			},
			Selector: Selector{
				Name: "",
			},
			Transport: Transport{
				Name: "",
			},
		},
	}
)
