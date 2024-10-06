package server

import (
	"errors"
	"fmt"

	"github.com/DecxBase/core/db"
	"github.com/DecxBase/core/logger"
	"github.com/DecxBase/core/options"
	"github.com/DecxBase/core/types"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type ComposedServer struct {
	Logger log.Logger

	opts         types.ServerOptions
	grpcHandlers []types.GrpcHandler
	httpHandlers []types.HttpHandler
}

func Create(fns ...types.ServerOptionsFunc) *ComposedServer {
	opts := types.NewServerOptions(fns...)
	godotenv.Load("local.env")

	server := &ComposedServer{
		opts:         opts,
		grpcHandlers: make([]types.GrpcHandler, 0),
		httpHandlers: make([]types.HttpHandler, 0),
		Logger:       logger.Create("app", opts.Name),
	}

	return server
}

func (s *ComposedServer) Name() string {
	return s.opts.Name
}

func (s *ComposedServer) RegisterGRPC(handlers ...types.GrpcHandler) {
	s.grpcHandlers = append(s.grpcHandlers, handlers...)
}

func (s *ComposedServer) RegisterHTTP(handlers ...types.HttpHandler) {
	s.httpHandlers = append(s.httpHandlers, handlers...)
}

func (s ComposedServer) Run(closers ...func()) error {
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

func (s ComposedServer) ResolveDBConnector() *pgdriver.Connector {
	return db.NewPGConnector(fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=%s",
		options.ReadEnv(s.Name(), "db_user", "postgres"),
		options.ReadEnv(s.Name(), "db_pass", ""),
		options.ReadEnv(s.Name(), "db_port", 5432),
		options.ReadEnv(s.Name(), "db_name", ""),
		options.ReadEnv(s.Name(), "db_ssl_mode", "disable"),
	))
}

func (s ComposedServer) ResolveDB() *bun.DB {
	return db.GetBunDB(s.ResolveDBConnector())
}
