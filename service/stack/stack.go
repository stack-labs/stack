package stack

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/debug/profile/pprof"
	"github.com/stack-labs/stack-rpc/debug/service/handler"
	"github.com/stack-labs/stack-rpc/env"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
)

type stackService struct {
	opts service.Options

	once sync.Once
}

func (s *stackService) Name() string {
	return s.opts.Name
}

func (s *stackService) Options() service.Options {
	return s.opts
}

func (s *stackService) Client() client.Client {
	return s.opts.Client
}

func (s *stackService) Server() server.Server {
	return s.opts.Server
}

func (s *stackService) String() string {
	return "stack"
}

func (s *stackService) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *stackService) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	if err := s.opts.Config.Close(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *stackService) Run() error {
	// register the debug handler
	if s.opts.Server.Options().EnableDebug {
		if err := s.opts.Server.Handle(
			s.opts.Server.NewHandler(
				handler.DefaultHandler,
				server.InternalHandler(true),
			),
		); err != nil {
			return err
		}
	}

	// start the profiler
	// TODO: set as an option to the service, don't just use pprof
	if prof := os.Getenv(env.StackDebugProfile); len(prof) > 0 {
		service := s.opts.Server.Options().Name
		version := s.opts.Server.Options().Version
		id := s.opts.Server.Options().Id
		profiler := pprof.NewProfile(
			profile.Name(service + "." + version + "." + id),
		)
		if err := profiler.Start(); err != nil {
			return err
		}
		defer profiler.Stop()
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

func NewService(opts ...service.Option) service.Service {
	options := newOptions(opts...)
	for _, o := range opts {
		o(&options)
	}
	return &stackService{
		opts: options,
	}
}
