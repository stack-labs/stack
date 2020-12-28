package jwt

import (
	"github.com/stack-labs/stack-rpc/auth/token"
	"github.com/stack-labs/stack-rpc/plugin"
)

type jwtTokenProviderPlugin struct {
}

func (j *jwtTokenProviderPlugin) Name() string {
	return "jwt"
}

func (j *jwtTokenProviderPlugin) Options() []token.Option {
	return nil
}

func (j *jwtTokenProviderPlugin) New(opts ...token.Option) token.Provider {
	return NewTokenProvider(opts...)
}

func init() {
	plugin.AuthTokenProviderPlugins["jwt"] = &jwtTokenProviderPlugin{}
}
