package wrapper

import (
	"context"
	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/debug/trace"
	"github.com/stack-labs/stack-rpc/server"
	"strings"
	"time"

	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/pkg/metadata"
)

type clientWrapper struct {
	client.Client
	headers metadata.Metadata
}

var (
	HeaderPrefix = "Stack-"
)

func (c *clientWrapper) setHeaders(ctx context.Context) context.Context {
	// copy metadata
	mda, _ := metadata.FromContext(ctx)
	md := metadata.Copy(mda)

	// set headers
	for k, v := range c.headers {
		if _, ok := md[k]; !ok {
			md[k] = v
		}
	}

	return metadata.NewContext(ctx, md)
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = c.setHeaders(ctx)
	return c.Client.Call(ctx, req, rsp, opts...)
}

func (c *clientWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	ctx = c.setHeaders(ctx)
	return c.Client.Stream(ctx, req, opts...)
}

func (c *clientWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	ctx = c.setHeaders(ctx)
	return c.Client.Publish(ctx, p, opts...)
}

// FromService wraps a client to inject From-Service header into metadata
func FromService(name string, c client.Client) client.Client {
	return &clientWrapper{
		c,
		metadata.Metadata{
			HeaderPrefix + "From-Service": name,
		},
	}
}

type traceWrapper struct {
	client.Client

	name  string
	trace trace.Tracer
}

// TraceCall is a call tracing wrapper
func TraceCall(name string, t trace.Tracer, c client.Client) client.Client {
	return &traceWrapper{
		name:   name,
		trace:  t,
		Client: c,
	}
}

// TraceHandler wraps a server handler to perform tracing
func TraceHandler(t trace.Tracer) server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// don't store traces for debug
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// get the span
			newCtx, s := t.Start(ctx, req.Service()+"."+req.Endpoint())
			s.Type = trace.SpanTypeRequestInbound

			err := h(newCtx, req, rsp)
			if err != nil {
				s.Metadata["error"] = err.Error()
			}

			// finish
			t.Finish(s)

			return err
		}
	}
}

type authWrapper struct {
	client.Client
	auth func() auth.Auth
}

func (a *authWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// parse the options
	var options client.CallOptions
	for _, o := range opts {
		o(&options)
	}

	// check to see if the authorization header has already been set.
	// We dont't override the header unless the ServiceToken option has
	// been specified or the header wasn't provided
	if _, ok := metadata.Get(ctx, "Authorization"); ok && !options.ServiceToken {
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// if auth is nil we won't be able to get an access token, so we execute
	// the request without one.
	aa := a.auth()
	if aa == nil {
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// check to see if we have a valid access token
	aaOpts := aa.Options()
	if aaOpts.ClientToken != nil && aaOpts.ClientToken.Expiry.Unix() > time.Now().Unix() {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+aaOpts.ClientToken.AccessToken)
		return a.Client.Call(ctx, req, rsp, opts...)
	}

	// call without an auth token
	return a.Client.Call(ctx, req, rsp, opts...)
}

// AuthClient wraps requests with the auth header
func AuthClient(auth func() auth.Auth, c client.Client) client.Client {
	return &authWrapper{c, auth}
}
