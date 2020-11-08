package stack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stack-labs/stack-rpc/config/source/file"
)

func TestStackConfig_Config(t *testing.T) {
	data := []byte(`
stack:
  service:
    name: demo.service
    rpc-port: 8081
    http-port: 8082`)
	path := filepath.Join(os.TempDir(), "file.yml")
	fh, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	c, err := newConfig(file.NewSource(file.WithPath(path)))
	if err != nil {
		t.Fatal(err)
	}

	if c.Config().Stack.Service.Name != "demo.service" {
		t.Fatal()
	}
	if c.Config().Stack.Service.RPCPort != 8081 {
		t.Fatal()
	}
	if c.Config().Stack.Service.HTTPPort != 8082 {
		t.Fatal()
	}
}
