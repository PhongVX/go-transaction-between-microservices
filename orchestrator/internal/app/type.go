package app

import (
	"google.golang.org/grpc"
)

type (
	Server struct {
		gServer *grpc.Server
	}
)
