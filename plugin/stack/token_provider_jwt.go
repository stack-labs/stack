package stack

import (
	"github.com/stack-labs/stack-rpc/auth/token"
	"github.com/stack-labs/stack-rpc/auth/token/jwt"
)

type jwtTokenProviderPlugin struct{}

func (j *jwtTokenProviderPlugin) Name() string {
	return "jwt"
}

func (j *jwtTokenProviderPlugin) Options() []token.Option {
	return nil
}

func (j *jwtTokenProviderPlugin) New(opts ...token.Option) token.Provider {
	return jwt.NewTokenProvider(opts...)
}
