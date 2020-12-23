package http

import (
	"errors"
	"net/http"

	"github.com/stack-labs/stack-rpc/client/selector"
)

type roundTripper struct {
	rt   http.RoundTripper
	st   selector.Strategy
	opts Options
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s, err := r.opts.Registry.GetService(req.URL.Host)
	if err != nil {
		return nil, err
	}

	// rudimentary retry 3 times
	for i := 0; i < 3; i++ {
		n, err := r.st(s)
		if err != nil {
			continue
		}
		req.URL.Host = n.Address
		w, err := r.rt.RoundTrip(req)
		if err != nil {
			continue
		}
		return w, nil
	}

	return nil, errors.New("failed request")
}
