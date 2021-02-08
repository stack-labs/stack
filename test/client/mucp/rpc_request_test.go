package mucp

import (
	"github.com/stack-labs/stack/plugin"
	"testing"

	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/client/mucp"
)

func TestRequestOptions(t *testing.T) {
	c := mucp.NewClient(
		client.Selector(plugin.SelectorPlugins["cache"].New()),
	)
	r := c.NewRequest("service", "endpoint", nil, client.WithContentType("application/json"))
	if r.Service() != "service" {
		t.Fatalf("expected 'service' got %s", r.Service())
	}
	if r.Endpoint() != "endpoint" {
		t.Fatalf("expected 'endpoint' got %s", r.Endpoint())
	}
	if r.ContentType() != "application/json" {
		t.Fatalf("expected 'endpoint' got %s", r.ContentType())
	}

	r2 := c.NewRequest("service", "endpoint", nil, client.WithContentType("application/json"), client.WithContentType("application/protobuf"))
	if r2.ContentType() != "application/protobuf" {
		t.Fatalf("expected 'endpoint' got %s", r2.ContentType())
	}
}
