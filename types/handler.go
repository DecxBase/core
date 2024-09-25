package types

import (
	"net/http"

	"github.com/phuslu/log"
	"google.golang.org/grpc"
)

type ServiceHandlerFunc = func(http.ResponseWriter, *http.Request) ServiceHandlerResult

type SeviceHandlerResultFunc = func() ServiceHandlerResult

type ServiceHandlerResult int

const (
	ServiceOk      ServiceHandlerResult = iota + 1 // EnumIndex = 1
	ServiceErr                                     // EnumIndex = 2
	ServiceFailed                                  // EnumIndex = 3
	ServiceInvalid                                 // EnumIndex = 4
)

func (r ServiceHandlerResult) String() string {
	return []string{"ok", "err", "failed", "invalid"}[r-1]
}

type GrpcHandler interface {
	GrpcIdentifier() string
	RegisterServer(*grpc.Server) error
	SetLogger(log.Logger)
}

type HttpHandler interface {
	HttpIdentifier() string
	RegisterRoutes(*http.ServeMux) error
	SetLogger(log.Logger)
}
