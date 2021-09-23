package app

import (
	"google.golang.org/grpc"
	"orchestrator/pkg/config"
)

type (
	Server struct {
		gServer *grpc.Server
		Config *config.Config
	}
)
