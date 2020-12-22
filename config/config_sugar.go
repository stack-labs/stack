package config

import (
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
)

var (
	_sugar Config
)

func Get(path ...string) reader.Value {
	return _sugar.Get(path...)
}
