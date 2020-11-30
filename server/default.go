package server

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var (
	DefaultAddress          = ":0"
	DefaultName             = "stack.rpc.server"
	DefaultVersion          = time.Now().Format("2006.01.02.15.04")
	DefaultId               = uuid.New().String()
	DefaultRegisterCheck    = func(context.Context) error { return nil }
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)
