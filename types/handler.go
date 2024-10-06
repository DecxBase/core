package types

import (
	"context"
	"net/http"

	"github.com/phuslu/log"
	"google.golang.org/grpc"
)

type HTTPHandler = func(w http.ResponseWriter, r *http.Request) (any, error)

type HandlerResultFunc = func() HandlerResult

type HandlerResult int

const (
	ServiceOk      HandlerResult = iota + 1 // EnumIndex = 1
	ServiceErr                              // EnumIndex = 2
	ServiceFailed                           // EnumIndex = 3
	ServiceInvalid                          // EnumIndex = 4
)

func (r HandlerResult) String() string {
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

type HandlerService[K, L, M, N any] interface {
	FindRecords(context.Context, K) (any, error)
	FindRecord(context.Context, L) (any, error)
	SaveRecord(context.Context, bool, M) error
	DeleteRecord(context.Context, N) error
}

type HttpHandlerImpl[K, L, M, N any] interface {
	FindRecords(context.Context, K) (any, error)
	FindRecord(context.Context, L) (any, error)
	SaveRecord(context.Context, bool, M) error
	DeleteRecord(context.Context, N) error
}
