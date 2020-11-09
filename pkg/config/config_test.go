package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stack-labs/stack-rpc/pkg/config/source/env"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
)

func createFile(t *testing.T, content string, format string) *os.File {
	data := []byte(content)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%s", format))
	fh, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
	}

	if err := fh.Close(); err != nil {
		t.Error(err)
	}

	return fh
}

func createFileForIssue18(t *testing.T, content string) *os.File {
	data := []byte(content)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
	}

	return fh
}

func createFileForTest(t *testing.T) *os.File {
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

	return fh
}

func TestConfigLoadWithGoodFile(t *testing.T) {
	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	// Create new config
	conf, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf.Close()

	// Load file source
	if err := conf.Load(file.NewSource(
		file.WithPath(path),
	)); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
}

func TestConfigLoadWithInvalidFile(t *testing.T) {
	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	// Create new config
	conf, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf.Close()

	// Load file source
	err = conf.Load(file.NewSource(
		file.WithPath(path),
		file.WithPath("/i/do/not/exists.json"),
	))

	if err == nil {
		t.Fatal("Expected error but none !")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "/i/do/not/exists.json") {
		t.Fatalf("Expected error to contain the unexisting file but got %v", err)
	}
}

func TestConfigMerge(t *testing.T) {
	fh := createFileForIssue18(t, `{
  "amqp": {
    "host": "rabbit.platform",
    "port": 80
  },
  "handler": {
    "exchange": "springCloudBus"
  }
}`)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()
	os.Setenv("AMQP_HOST", "rabbit.testing.com")

	conf, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf.Close()
	if err := conf.Load(
		file.NewSource(
			file.WithPath(path),
		),
		env.NewSource(),
	); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	actualHost := conf.Get("amqp", "host").String("backup")
	if actualHost != "rabbit.testing.com" {
		t.Fatalf("Expected %v but got %v",
			"rabbit.testing.com",
			actualHost)
	}
}

func TestConfigLoadFromBackupFile(t *testing.T) {
	fh := createFileForIssue18(t, `{
  "amqp": {
    "host": "rabbit.platform",
    "port": 80
  },
  "handler": {
    "exchange": "springCloudBus"
  }
}`)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	conf, err := NewConfig(Storage(true), StorageDir(os.TempDir()))
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf.Close()

	if err := conf.Load(
		file.NewSource(
			file.WithPath(path),
		),
	); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	conf2, err := NewConfig(Storage(true), StorageDir(os.TempDir()))
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf2.Close()
	if err := conf2.Load(
		file.NewSource(
			file.WithPath("/i/do/not/exists.json"),
		),
	); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	actualHost := conf.Get("amqp", "host").String("backup")
	if actualHost != "rabbit.platform" {
		t.Fatalf("Expected %v but got %v",
			"rabbit.platform",
			actualHost)
	}
}

func TestYmlConfigLoadFromBackupFile(t *testing.T) {
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

	conf, err := NewConfig(Storage(true), StorageDir(os.TempDir()))
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf.Close()

	if err := conf.Load(
		file.NewSource(
			file.WithPath(path),
		),
	); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	conf2, err := NewConfig(Storage(true), StorageDir(os.TempDir()))
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	defer conf2.Close()
	if err := conf2.Load(
		file.NewSource(
			file.WithPath("/i/do/not/exists.json"),
		),
	); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	port := conf.Get("stack", "service", "rpc-port").Int(1)
	if port != 8081 {
		t.Fatalf("Expected %d but got %d",
			8081,
			port)
	}
}

func TestSomeConfigLoad(t *testing.T) {
	formats := []string{"json", "yaml", "toml", "xml", "hcl", "yml"}
	for _, v := range formats {
		fh := createFile(t, "", v)
		path := fh.Name()
		defer func() {
			os.Remove(path)
		}()

		conf, err := NewConfig()
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
		defer conf.Close()

		if err := conf.Load(
			file.NewSource(
				file.WithPath(path),
			),
		); err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
	}
}
