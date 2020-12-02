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
					Size: "",
					TTL:  "",
				},
				Request: ClientRequest{
					Retries: "",
					Timeout: "",
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

func GetDefault() *Value {
	return _defaultCfg
}
