package config

import (
	"github.com/micro/go-config/reader"
	"github.com/micro/go-micro/config"
)

func GetConfig(names ...string) reader.Value {
	names2 := append([]string{"micro", "platform_web"}, names...)
	return config.Get(names2...)
}
