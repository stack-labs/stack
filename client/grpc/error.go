package grpc

import (
	"github.com/stack-labs/stack-rpc/errors"
	"google.golang.org/grpc/status"
)

func stackError(err error) error {
	// no error
	switch err {
	case nil:
		return nil
	}

	// stack error
	if v, ok := err.(*errors.Error); ok {
		return v
	}

	// grpc error
	if s, ok := status.FromError(err); ok {
		if e := errors.Parse(s.Message()); e.Code > 0 {
			return e // actually a stack error
		}
		return errors.InternalServerError("stack.rpc.client", s.Message())
	}

	// do nothing
	return err
}
