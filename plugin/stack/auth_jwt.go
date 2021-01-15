package stack

import (
	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/auth/jwt"
)

type jwtAuthPlugin struct{}

func (j *jwtAuthPlugin) Name() string {
	return "jwt"
}

func (j *jwtAuthPlugin) Options() []auth.Option {
	return nil
}

func (j *jwtAuthPlugin) New(opts ...auth.Option) auth.Auth {
	return jwt.NewAuth(opts...)
}
