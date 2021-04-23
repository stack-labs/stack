package main

import (
	"time"

	"github.com/stack-labs/stack"
	"github.com/stack-labs/stack/config"
	log "github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/pkg/config/source/file"
)

type source struct {
	DemoA        string `sc:"demoA"`
	NumberString string `sc:"number-string"`
}

type Value struct {
	Source source `sc:"source"`
}

var (
	value Value
)

func init() {
	config.RegisterOptions(&value)
}

func main() {
	service := stack.NewService(
		stack.Name("stack.config.demo"),
		stack.Config(config.NewConfig(config.Source(file.NewSource(file.WithPath("./source.yml"))))),
	)
	service.Init()

	log.Infof("demoA: %s", value.Source.DemoA)
	log.Infof("NumberString: %s", value.Source.NumberString)

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				// try to change DemoA value in source.yml
				// there will log the new value
				log.Infof("demoA: %s", value.Source.DemoA)
			}
		}
	}()
	service.Run()
}
