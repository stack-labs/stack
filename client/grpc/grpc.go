// Package grpc provides a gRPC client
package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/broker/http"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	selectorR "github.com/stack-labs/stack-rpc/client/selector/registry"
	"github.com/stack-labs/stack-rpc/codec"
	raw "github.com/stack-labs/stack-rpc/codec/bytes"
	"github.com/stack-labs/stack-rpc/pkg/metadata"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/registry/mdns"
	"github.com/stack-labs/stack-rpc/transport"
	"github.com/stack-labs/stack-rpc/util/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	gmetadata "google.golang.org/grpc/metadata"
)

type grpcClient struct {
	once sync.Once
	opts client.Options
	pool *pool
}

func init() {
	encoding.RegisterCodec(wrapCodec{jsonCodec{}})
	encoding.RegisterCodec(wrapCodec{protoCodec{}})
	encoding.RegisterCodec(wrapCodec{bytesCodec{}})
}

// secure returns the dial option for whether its a secure or insecure connection
func (g *grpcClient) secure() grpc.DialOption {
	if g.opts.Context != nil {
		if v := g.opts.Context.Value(tlsAuth{}); v != nil {
			tls := v.(*tls.Config)
			creds := credentials.NewTLS(tls)
			return grpc.WithTransportCredentials(creds)
		}
	}
	return grpc.WithInsecure()
}

func (g *grpcClient) next(request client.Request, opts client.CallOptions) (selector.Next, error) {
	service := request.Service()

	// get proxy
	if prx := os.Getenv("STACKPROXY"); len(prx) > 0 {
		service = prx
	}

	// get proxy address
	if prx := os.Getenv("STACKPROXY_ADDRESS"); len(prx) > 0 {
		opts.Address = []string{prx}
	}

	// return remote address
	if len(opts.Address) > 0 {
		return func() (*registry.Node, error) {
			return &registry.Node{
				Address: opts.Address[0],
			}, nil
		}, nil
	}

	// get next nodes from the selector
	next, err := g.opts.Selector.Select(service, opts.SelectOptions...)
	if err != nil {
		if err == selector.ErrNotFound {
			return nil, errors.InternalServerError("stack.rpc.client", "service %s: %s", service, err.Error())
		}
		return nil, errors.InternalServerError("stack.rpc.client", "error selecting %s node: %s", service, err.Error())
	}

	return next, nil
}

func (g *grpcClient) call(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
	address := node.Address

	header := make(map[string]string)
	if md, ok := metadata.FromContext(ctx); ok {
		for k, v := range md {
			header[k] = v
		}
	}

	// set timeout in nanoseconds
	header["timeout"] = fmt.Sprintf("%d", opts.RequestTimeout)
	// set the content type for the request
	header["x-content-type"] = req.ContentType()

	md := gmetadata.New(header)
	ctx = gmetadata.NewOutgoingContext(ctx, md)

	cf, err := g.newGRPCCodec(req.ContentType())
	if err != nil {
		return errors.InternalServerError("stack.rpc.client", err.Error())
	}

	maxRecvMsgSize := g.maxRecvMsgSizeValue()
	maxSendMsgSize := g.maxSendMsgSizeValue()

	var grr error

	grpcDialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.ForceCodec(cf)),
		grpc.WithTimeout(opts.DialTimeout),
		g.secure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize),
		),
	}

	if opts := g.getGrpcDialOptions(); opts != nil {
		grpcDialOptions = append(grpcDialOptions, opts...)
	}

	cc, err := g.pool.getConn(address, grpcDialOptions...)
	if err != nil {
		return errors.InternalServerError("stack.rpc.client", fmt.Sprintf("Error sending request: %v", err))
	}
	defer func() {
		// defer execution of release
		g.pool.release(address, cc, grr)
	}()

	ch := make(chan error, 1)

	go func() {
		grpcCallOptions := []grpc.CallOption{grpc.CallContentSubtype(cf.Name())}
		if opts := g.getGrpcCallOptions(); opts != nil {
			grpcCallOptions = append(grpcCallOptions, opts...)
		}
		err := cc.Invoke(ctx, methodToGRPC(req.Service(), req.Endpoint()), req.Body(), rsp, grpcCallOptions...)
		ch <- stackError(err)
	}()

	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = ctx.Err()
	}

	return grr
}

func (g *grpcClient) stream(ctx context.Context, node *registry.Node, req client.Request, opts client.CallOptions) (client.Stream, error) {
	address := node.Address

	header := make(map[string]string)
	if md, ok := metadata.FromContext(ctx); ok {
		for k, v := range md {
			header[k] = v
		}
	}

	// set timeout in nanoseconds
	header["timeout"] = fmt.Sprintf("%d", opts.RequestTimeout)
	// set the content type for the request
	header["x-content-type"] = req.ContentType()

	md := gmetadata.New(header)
	ctx = gmetadata.NewOutgoingContext(ctx, md)

	cf, err := g.newGRPCCodec(req.ContentType())
	if err != nil {
		return nil, errors.InternalServerError("stack.rpc.client", err.Error())
	}

	var dialCtx context.Context
	var cancel context.CancelFunc
	if opts.DialTimeout >= 0 {
		dialCtx, cancel = context.WithTimeout(ctx, opts.DialTimeout)
	} else {
		dialCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	wc := wrapCodec{cf}

	grpcDialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.ForceCodec(wc)),
		g.secure(),
	}

	if opts := g.getGrpcDialOptions(); opts != nil {
		grpcDialOptions = append(grpcDialOptions, opts...)
	}

	cc, err := grpc.DialContext(dialCtx, address, grpcDialOptions...)
	if err != nil {
		return nil, errors.InternalServerError("stack.rpc.client", fmt.Sprintf("Error sending request: %v", err))
	}

	desc := &grpc.StreamDesc{
		StreamName:    req.Service() + req.Endpoint(),
		ClientStreams: true,
		ServerStreams: true,
	}

	grpcCallOptions := []grpc.CallOption{}
	if opts := g.getGrpcCallOptions(); opts != nil {
		grpcCallOptions = append(grpcCallOptions, opts...)
	}
	st, err := cc.NewStream(ctx, desc, methodToGRPC(req.Service(), req.Endpoint()), grpcCallOptions...)
	if err != nil {
		return nil, errors.InternalServerError("stack.rpc.client", fmt.Sprintf("Error creating stream: %v", err))
	}

	codec := &grpcCodec{
		s: st,
		c: wc,
	}

	// set request codec
	if r, ok := req.(*grpcRequest); ok {
		r.codec = codec
	}

	rsp := &response{
		conn:   cc,
		stream: st,
		codec:  cf,
		gcodec: codec,
	}

	return &grpcStream{
		context:  ctx,
		request:  req,
		response: rsp,
		stream:   st,
		conn:     cc,
	}, nil
}

func (g *grpcClient) maxRecvMsgSizeValue() int {
	if g.opts.Context == nil {
		return DefaultMaxRecvMsgSize
	}
	v := g.opts.Context.Value(maxRecvMsgSizeKey{})
	if v == nil {
		return DefaultMaxRecvMsgSize
	}
	return v.(int)
}

func (g *grpcClient) maxSendMsgSizeValue() int {
	if g.opts.Context == nil {
		return DefaultMaxSendMsgSize
	}
	v := g.opts.Context.Value(maxSendMsgSizeKey{})
	if v == nil {
		return DefaultMaxSendMsgSize
	}
	return v.(int)
}

func (g *grpcClient) newGRPCCodec(contentType string) (encoding.Codec, error) {
	codecs := make(map[string]encoding.Codec)
	if g.opts.Context != nil {
		if v := g.opts.Context.Value(codecsKey{}); v != nil {
			codecs = v.(map[string]encoding.Codec)
		}
	}
	if c, ok := codecs[contentType]; ok {
		return wrapCodec{c}, nil
	}
	if c, ok := defaultGRPCCodecs[contentType]; ok {
		return wrapCodec{c}, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (g *grpcClient) Init(opts ...client.Option) error {
	size := g.opts.PoolSize
	ttl := g.opts.PoolTTL

	for _, o := range opts {
		o(&g.opts)
	}

	// update pool configuration if the options changed
	if size != g.opts.PoolSize || ttl != g.opts.PoolTTL {
		g.pool.Lock()
		g.pool.size = g.opts.PoolSize
		g.pool.ttl = int64(g.opts.PoolTTL.Seconds())
		g.pool.Unlock()
	}

	return nil
}

func (g *grpcClient) Options() client.Options {
	return g.opts
}

func (g *grpcClient) NewMessage(topic string, msg interface{}, opts ...client.MessageOption) client.Message {
	return newGRPCEvent(topic, msg, g.opts.ContentType, opts...)
}

func (g *grpcClient) NewRequest(service, method string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return newGRPCRequest(service, method, req, g.opts.ContentType, reqOpts...)
}

func (g *grpcClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// make a copy of call opts
	callOpts := g.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := g.next(req, callOpts)
	if err != nil {
		return err
	}

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		opt := client.WithRequestTimeout(time.Until(d))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return errors.New("stack.rpc.client", fmt.Sprintf("%v", ctx.Err()), 408)
	default:
	}

	// make copy of call method
	gcall := g.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		gcall = callOpts.CallWrappers[i-1](gcall)
	}

	// return errors.New("stack.rpc.client", "request timeout", 408)
	call := func(i int) error {
		// call backoff first. Someone may want an initial start delay
		t, err := callOpts.Backoff(ctx, req, i)
		if err != nil {
			return errors.InternalServerError("stack.rpc.client", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// select next node
		node, err := next()
		service := req.Service()
		if err != nil {
			if err == selector.ErrNotFound {
				return errors.InternalServerError("stack.rpc.client", "service %s: %s", service, err.Error())
			}
			return errors.InternalServerError("stack.rpc.client", "error selecting %s node: %s", service, err.Error())
		}

		// make the call
		err = gcall(ctx, node, req, rsp, callOpts)
		g.opts.Selector.Mark(service, node, err)
		return err
	}

	ch := make(chan error, callOpts.Retries+1)
	var gerr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return errors.New("stack.rpc.client", fmt.Sprintf("%v", ctx.Err()), 408)
		case err := <-ch:
			// if the call succeeded lets bail early
			if err == nil {
				return nil
			}

			retry, rerr := callOpts.Retry(ctx, req, i, err)
			if rerr != nil {
				return rerr
			}

			if !retry {
				return err
			}

			gerr = err
		}
	}

	return gerr
}

func (g *grpcClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	// make a copy of call opts
	callOpts := g.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := g.next(req, callOpts)
	if err != nil {
		return nil, err
	}

	// #200 - streams shouldn't have a request timeout set on the context

	// should we noop right here?
	select {
	case <-ctx.Done():
		return nil, errors.New("stack.rpc.client", fmt.Sprintf("%v", ctx.Err()), 408)
	default:
	}

	call := func(i int) (client.Stream, error) {
		// call backoff first. Someone may want an initial start delay
		t, err := callOpts.Backoff(ctx, req, i)
		if err != nil {
			return nil, errors.InternalServerError("stack.rpc.client", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		node, err := next()
		service := req.Service()
		if err != nil {
			if err == selector.ErrNotFound {
				return nil, errors.InternalServerError("stack.rpc.client", "service %s: %s", service, err.Error())
			}
			return nil, errors.InternalServerError("stack.rpc.client", "error selecting %s node: %s", service, err.Error())
		}

		stream, err := g.stream(ctx, node, req, callOpts)
		g.opts.Selector.Mark(service, node, err)
		return stream, err
	}

	type response struct {
		stream client.Stream
		err    error
	}

	ch := make(chan response, callOpts.Retries+1)
	var grr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func(i int) {
			s, err := call(i)
			ch <- response{s, err}
		}(i)

		select {
		case <-ctx.Done():
			return nil, errors.New("stack.rpc.client", fmt.Sprintf("%v", ctx.Err()), 408)
		case rsp := <-ch:
			// if the call succeeded lets bail early
			if rsp.err == nil {
				return rsp.stream, nil
			}

			retry, rerr := callOpts.Retry(ctx, req, i, err)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, rsp.err
			}

			grr = rsp.err
		}
	}

	return nil, grr
}

func (g *grpcClient) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	var options client.PublishOptions
	for _, o := range opts {
		o(&options)
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	md["Content-Type"] = p.ContentType()
	md["Stack-Topic"] = p.Topic()

	cf, err := g.newGRPCCodec(p.ContentType())
	if err != nil {
		return errors.InternalServerError("stack.rpc.client", err.Error())
	}

	var body []byte

	// passed in raw data
	if d, ok := p.Payload().(*raw.Frame); ok {
		body = d.Data
	} else {
		// set the body
		b, err := cf.Marshal(p.Payload())
		if err != nil {
			return errors.InternalServerError("stack.rpc.client", err.Error())
		}
		body = b
	}

	g.once.Do(func() {
		g.opts.Broker.Connect()
	})

	topic := p.Topic()

	// get proxy topic
	if prx := os.Getenv("STACK_PROXY"); len(prx) > 0 {
		options.Exchange = prx
	}

	// get the exchange
	if len(options.Exchange) > 0 {
		topic = options.Exchange
	}

	return g.opts.Broker.Publish(topic, &broker.Message{
		Header: md,
		Body:   body,
	})
}

func (g *grpcClient) String() string {
	return "grpc"
}

func (g *grpcClient) getGrpcDialOptions() []grpc.DialOption {
	if g.opts.CallOptions.Context == nil {
		return nil
	}

	v := g.opts.CallOptions.Context.Value(grpcDialOptions{})

	if v == nil {
		return nil
	}

	opts, ok := v.([]grpc.DialOption)

	if !ok {
		return nil
	}

	return opts
}

func (g *grpcClient) getGrpcCallOptions() []grpc.CallOption {
	if g.opts.CallOptions.Context == nil {
		return nil
	}

	v := g.opts.CallOptions.Context.Value(grpcCallOptions{})

	if v == nil {
		return nil
	}

	opts, ok := v.([]grpc.CallOption)

	if !ok {
		return nil
	}

	return opts
}

func newClient(opts ...client.Option) client.Client {
	options := client.Options{
		Codecs: make(map[string]codec.NewCodec),
		CallOptions: client.CallOptions{
			Backoff:        client.DefaultBackoff,
			Retry:          client.DefaultRetry,
			Retries:        client.DefaultRetries,
			RequestTimeout: client.DefaultRequestTimeout,
			DialTimeout:    transport.DefaultDialTimeout,
		},
		PoolSize: client.DefaultPoolSize,
		PoolTTL:  client.DefaultPoolTTL,
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.ContentType) == 0 {
		options.ContentType = "application/grpc+proto"
	}

	if options.Broker == nil {
		options.Broker = http.NewBroker()
	}

	if options.Registry == nil {
		options.Registry = mdns.NewRegistry()
	}

	if options.Selector == nil {
		options.Selector = selectorR.NewSelector(
			selector.Registry(options.Registry),
		)
	}

	rc := &grpcClient{
		once: sync.Once{},
		opts: options,
		pool: newPool(options.PoolSize, options.PoolTTL),
	}

	c := client.Client(rc)

	// wrap in reverse
	for i := len(options.Wrappers); i > 0; i-- {
		c = options.Wrappers[i-1](c)
	}

	return c
}

func NewClient(opts ...client.Option) client.Client {
	return newClient(opts...)
}
