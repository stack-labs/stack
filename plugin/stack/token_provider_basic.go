package stack

import (
	"github.com/stack-labs/stack-rpc/auth/token"
	"github.com/stack-labs/stack-rpc/auth/token/basic"
)

type basicTokenProviderPlugin struct{}

func (j *basicTokenProviderPlugin) Name() string {
	return "basic"
}

func (j *basicTokenProviderPlugin) Options() []token.Option {
	return nil
}

func (j *basicTokenProviderPlugin) New(opts ...token.Option) token.Provider {
	return basic.NewTokenProvider(opts...)
}
