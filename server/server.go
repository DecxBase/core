package server

import (
	"errors"

	"github.com/DecxBase/core/logger"
	"github.com/DecxBase/core/types"
	"github.com/phuslu/log"
)

type composedServer struct {
	Logger log.Logger

	opts         types.ServerOptions
	grpcHandlers []types.GrpcHandler
	httpHandlers []types.HttpHandler
}

func Create(fns ...types.ServerOptionsFunc) *composedServer {
	opts := types.NewServerOptions(fns...)

	return &composedServer{
		opts:         opts,
		grpcHandlers: make([]types.GrpcHandler, 0),
		httpHandlers: make([]types.HttpHandler, 0),
		Logger:       logger.Create("app", opts.Name),
	}
}

func (s *composedServer) Name() string {
	return s.opts.Name
}

func (s *composedServer) RegisterGRPC(handlers ...types.GrpcHandler) {
	s.grpcHandlers = append(s.grpcHandlers, handlers...)
}

func (s *composedServer) RegisterHTTP(handlers ...types.HttpHandler) {
	s.httpHandlers = append(s.httpHandlers, handlers...)
}

func (s composedServer) Run(closers ...func()) error {
	hasGRPC := len(s.grpcHandlers) > 0
	hasHTTP := len(s.httpHandlers) > 0

	if hasGRPC && hasHTTP {
		go s.RunHttp()
		return s.RunGrpc(closers...)
	} else if hasGRPC {
		return s.RunGrpc(closers...)
	} else if hasHTTP {
		return s.RunHttp(closers...)
	}

	return errors.New("no grpc/http handlers registered")
}
