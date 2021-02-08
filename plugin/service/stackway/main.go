package main

import (
	"github.com/stack-labs/stack"
	"github.com/stack-labs/stack/util/log"

	"github.com/stack-labs/stack/plugin/service/stackway/api"
)

func main() {
	svc := stack.NewService()

	// stackway server
	apiServer := api.NewServer(svc)
	svc.Init(apiServer.Options()...)

	// run service
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
