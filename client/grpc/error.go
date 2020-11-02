package grpc

import (
	"google.golang.org/grpc/status"
	"github.com/stack-labs/stack-rpc/errors"
)

func microError(err error) error {
	// no error
	switch err {
	case nil:
		return nil
	}

	// micro error
	if v, ok := err.(*errors.Error); ok {
		return v
	}

	// grpc error
	if s, ok := status.FromError(err); ok {
		if e := errors.Parse(s.Message()); e.Code > 0 {
			return e // actually a micro error
		}
		return errors.InternalServerError("stack.rpc.client", s.Message())
	}

	// do nothing
	return err
}
