package main

import (
	"github.com/stack-labs/stack"
	log "github.com/stack-labs/stack/logger"
)

func main() {
	service := stack.NewService()
	service.Init()
	log.Debug("hello，这是Debug级别")
	log.Info("hello，这是Info级别")
}
