package api

import (
	"reflect"
	"testing"

	"github.com/stack-labs/stack-rpc"
)

func TestNewServer(t *testing.T) {
	type args struct {
		svc stack.Service
	}
	tests := []struct {
		name string
		args args
		want *httpServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
