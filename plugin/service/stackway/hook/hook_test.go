package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/api"
	hApi "github.com/stack-labs/stack-rpc/api/handler/api"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	memSource "github.com/stack-labs/stack-rpc/pkg/config/source/memory"
	"github.com/stack-labs/stack-rpc/registry/memory"
	"github.com/stretchr/testify/assert"

	"github.com/stack-labs/stack-rpc-plugins/service/stackway/test/handler"
	test "github.com/stack-labs/stack-rpc-plugins/service/stackway/test/proto"
)

func run(ctx context.Context, t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	svc := stack.NewService(
		stack.Name("stack.rpc.api.example"),
		stack.Flags(cli.StringFlag{Name: "test.v"}),
		stack.Flags(cli.StringFlag{Name: "test.timeout"}),
		stack.Flags(cli.StringFlag{Name: "test.count"}),
		stack.Flags(cli.StringFlag{Name: "test.coverprofile"}),
		stack.Flags(cli.StringFlag{Name: "test.testlogfile"}),
	)

	yamlConf := `
stack:
  registry:
    name: memory

  stackway:
    address: :8080
    handler: "meta"
    resolver: "stack"
    rpc_path: "/rpc"
    api_path: "/"
    proxy_path: "/{service:[a-zA-Z0-9]+}"
    namespace: "stack.rpc.api"
    header_prefix: "X-Stack-"
    enable_rpc: true
    enable_acme: false
    enable_tls: false
    acme:
      provider: "autocert"
      challenge_provider: "cloudflare"
      ca: "https://acme-v02.api.letsencrypt.org/directory"
      hosts:
        - ""
    # Plugins
    example:
      key: value
`

	reg := memory.NewRegistry()
	_ = svc.Init(
		stack.Registry(reg),
		stack.Context(ctx),
		stack.Config(
			config.NewConfig(
				config.Source(
					memSource.NewSource(memSource.WithYAML([]byte(yamlConf)))),
			),
		),
		stack.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		stack.AfterStop(func() error {
			return nil
		}),
	)

	_ = test.RegisterTestHandler(svc.Server(), &handler.Handler{},
		api.WithEndpoint(&api.Endpoint{
			// The RPC method
			Name: "Test.Api",
			// The HTTP paths. This can be a POSIX regex
			Path: []string{"^/api/test/handler/api$"},
			// The HTTP Methods for this endpoint
			Method: []string{"GET"},
			// The API handler to use
			Handler: hApi.Handler,
		}),
	)

	// run service
	go func() {
		// stackway hook
		Hook(svc)

		if err := svc.Run(); err != nil {
			t.Fatalf("service run error: %v", err)
		}
	}()

	wg.Wait()
}

func TestHook(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	run(ctx, t)

	type args struct {
		url      string
		endpoint string
		want     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handler/rpc",
			args: args{
				url:  "http://localhost:8080/example/test/rpc?msg=ok",
				want: `{"msg":"ok"}`,
			},
		},
		{
			name: "handler/api",
			args: args{
				url:  "http://localhost:8080/api/test/handler/api?msg=ok",
				want: "ok",
			},
		},
		{
			name: "rpc",
			args: args{
				url:      "http://localhost:8080/rpc",
				endpoint: "Test.Rpc",
				want:     `{"msg":"ok"}`,
			},
		},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp string
			var err error
			if len(tt.args.endpoint) > 0 {
				resp, err = rpcRequest(tt.args.url, tt.args.endpoint)
			} else {
				resp, err = request(tt.args.url)
			}

			assertions.NoError(err)
			assertions.Equal(tt.args.want, resp)

			t.Log(resp)
		})
	}

	cancel()
}

func request(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	return string(b), nil
}

func rpcRequest(url, endpoint string) (string, error) {
	request := map[string]string{
		"service":  "stack.rpc.api.example",
		"endpoint": endpoint,
		"request":  `{"msg":"ok"}`,
	}

	rb, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(rb)

	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
