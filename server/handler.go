package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/DecxBase/core/types"
	"github.com/phuslu/log"
)

type ServiceHandler struct {
	Logger log.Logger
}

func NewServiceHandler() ServiceHandler {
	return ServiceHandler{}
}

func (h *ServiceHandler) SetLogger(l log.Logger) {
	h.Logger = l
}

func (h ServiceHandler) Log(msg string, v ...any) {
	h.Logger.Debug().Msgf(msg, v...)
}

func (h ServiceHandler) AccessLog(method string, path string) {
	h.Logger.Debug().Msgf("[%s] %s", method, path)
}

func (h ServiceHandler) RuntimeLog(name string, msg string, cb types.SeviceHandlerResultFunc) {
	var result types.ServiceHandlerResult
	started := time.Now()

	defer func() {
		diff := time.Since(started)

		h.Logger.Info().Msgf(
			"[%s] %s %s - %s",
			name, msg,
			diff.String(), strings.ToUpper(result.String()),
		)
	}()

	result = cb()
}

func (h ServiceHandler) HTTPRuntimeLog(r *http.Request, cb types.SeviceHandlerResultFunc) {
	h.RuntimeLog(r.Method, r.URL.Path, cb)
}

func (h ServiceHandler) HTTPHandler(handler types.ServiceHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTTPRuntimeLog(r, func() types.ServiceHandlerResult {
			return handler(w, r)
		})
	}
}

func (h ServiceHandler) GRPCRuntimeLog(name string, since time.Time, result types.ServiceHandlerResult) {
	diff := time.Since(since)

	h.Logger.Info().Msgf(
		"%s %s - %s",
		name, diff.String(), strings.ToUpper(result.String()),
	)
}

func (h ServiceHandler) ErrorToServiceResult(err error, t ...types.ServiceHandlerResult) types.ServiceHandlerResult {
	if err == nil {
		return types.ServiceOk
	}

	if len(t) > 0 {
		return t[0]
	}

	return types.ServiceErr
}
