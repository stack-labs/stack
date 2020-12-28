package jwt

import (
	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/plugin"
)

type jwtAuthPlugin struct {
}

func (j *jwtAuthPlugin) Name() string {
	return "jwt"
}

func (j *jwtAuthPlugin) Options() []auth.Option {
	return nil
}

func (j *jwtAuthPlugin) New(opts ...auth.Option) auth.Auth {
	return NewAuth(opts...)
}

func init() {
	plugin.AuthPlugins["jwt"] = &jwtAuthPlugin{}
}
