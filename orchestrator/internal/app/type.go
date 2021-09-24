package app

import (
	"google.golang.org/grpc"
	"orchestrator/internal/config"
)

type (
	Server struct {
		gServer *grpc.Server
		Config *config.Config
	}
)
