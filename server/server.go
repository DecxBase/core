package server

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/DecxBase/core/db"
	"github.com/DecxBase/core/logger"
	"github.com/DecxBase/core/options"
	"github.com/DecxBase/core/types"
	"github.com/DecxBase/core/utils"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
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

func (s ComposedServer) ResolveDSN() (string, string) {
	customDrivers := []string{"mysql"}

	dsn := options.ReadEnv(s.Name(), "DSN", "")
	driver := options.ReadEnv(s.Name(), "db_driver", "postgres")

	if len(dsn) > 0 {
		segments := strings.Split(dsn, "://")

		if utils.CheckContains(customDrivers, segments[0]) {
			return segments[0], segments[1]
		}
		return segments[0], dsn
	}

	if utils.CheckContains(customDrivers, driver) {
		return driver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s",
			options.ReadEnv(s.Name(), "db_user", "root"),
			options.ReadEnv(s.Name(), "db_pass", ""),
			options.ReadEnv(s.Name(), "db_protocol", "tcp"),
			options.ReadEnv(s.Name(), "db_host", "localhost"),
			options.ReadEnv(s.Name(), "db_port", 3306),
			options.ReadEnv(s.Name(), "db_name", ""),
			options.ReadEnv(s.Name(), "db_extra_params", ""),
		)
	}

	return driver, fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s&%s", driver,
		options.ReadEnv(s.Name(), "db_user", "postgres"),
		options.ReadEnv(s.Name(), "db_pass", ""),
		options.ReadEnv(s.Name(), "db_host", "localhost"),
		options.ReadEnv(s.Name(), "db_port", 5432),
		options.ReadEnv(s.Name(), "db_name", ""),
		options.ReadEnv(s.Name(), "db_ssl_mode", "disable"),
		options.ReadEnv(s.Name(), "db_extra_params", ""),
	)
}

func (s ComposedServer) DBOpen(dialect schema.Dialect) *bun.DB {
	driver, dsn := s.ResolveDSN()
	sqlDB, err := db.Open(driver, dsn)

	if err != nil {
		panic(err)
	}

	return db.Transform(sqlDB, dialect)
}

func (s ComposedServer) DBOpenDB(connector driver.Connector, dialect schema.Dialect) *bun.DB {
	sqlDB := db.OpenDB(connector)
	return db.Transform(sqlDB, dialect)
}
