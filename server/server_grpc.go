package server

import (
	"errors"

	"github.com/DecxBase/core/logger"
	"github.com/phuslu/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s ComposedServer) BuildGrpc(l log.Logger) (*grpc.Server, error) {
	srv := grpc.NewServer()

	for _, hnd := range s.grpcHandlers {
		hnd.SetLogger(logger.CreateFrom(l, "service", hnd.GrpcIdentifier()))
		hnd.RegisterServer(srv)
	}

	return srv, nil
}

func (s ComposedServer) RunGrpc(closers ...func()) error {
	if len(s.grpcHandlers) < 1 {
		return errors.New("no grpc handlers registered")
	}

	l := logger.CreateFrom(s.Logger, "type", "grpc")

	grpcServer, err := s.BuildGrpc(l)
	if err != nil {
		return err
	}

	if s.opts.ReflectGRPC {
		reflection.Register(grpcServer)
	}

	errC, err := GracefulGrpc(grpcServer, s.opts.GrpcAddr(), l)
	if err != nil {
		return err
	}

	if err := <-errC; err != nil {
		return err
	}

	return nil
}
