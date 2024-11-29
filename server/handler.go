package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/DecxBase/core/exception"
	"github.com/DecxBase/core/types"
	"github.com/DecxBase/core/utils"
	"github.com/phuslu/log"
)

type Handler struct {
	Logger log.Logger
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) SetLogger(l log.Logger) {
	h.Logger = l
}

func (h Handler) Log(msg string, v ...any) {
	h.Logger.Debug().Msgf(msg, v...)
}

func (h Handler) AccessLog(method string, path string) {
	h.Logger.Debug().Msgf("[%s] %s", method, path)
}

func (h Handler) RuntimeLog(name string, msg string, cb types.HandlerResultFunc) {
	var result types.HandlerResult
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

func (h Handler) HTTPRuntimeLog(r *http.Request, cb types.HandlerResultFunc) {
	h.RuntimeLog(r.Method, r.URL.Path, cb)
}

func (h Handler) HTTPHandler(handler types.HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTTPRuntimeLog(r, func() types.HandlerResult {
			res, rErr := handler(w, r)

			if rErr != nil {
				err, ok := rErr.(*exception.ServerException)

				if ok {
					err.WriteToResponse(w)

					if err.StatusCode == http.StatusInternalServerError {
						return types.ServiceErr
					}

					return types.ServiceFailed
				}

				utils.WriteResponseError(w, http.StatusInternalServerError, rErr)

				return types.ServiceFailed
			}

			if res != nil {
				utils.WriteResponseData(w, res)
			}
			return types.ServiceOk
		})
	}
}

func (h Handler) GRPCRuntimeLog(name string, since time.Time, result types.HandlerResult) {
	diff := time.Since(since)

	h.Logger.Info().Msgf(
		"%s %s - %s",
		name, diff.String(), strings.ToUpper(result.String()),
	)
}

func (h Handler) ErrorToServiceResult(err error, t ...types.HandlerResult) types.HandlerResult {
	if err == nil {
		return types.ServiceOk
	}

	if len(t) > 0 {
		return t[0]
	}

	return types.ServiceErr
}

func (h Handler) ParseQuery(r *http.Request, target any) {
	if err := utils.ParseQuery(r, target); err != nil {
		h.Logger.Warn().Msgf("parse query failed: %s", err.Error())
	}
}

func (h Handler) ParseBody(r *http.Request, target any) error {
	if err := utils.ParseBody(r, target); err != nil {
		return exception.Raise(err).WithCode("parse_err")
	}

	return nil
}

func (h Handler) ParseRecordID(r *http.Request, name string) int64 {
	return utils.StringToInt64(r.PathValue(name))
}

func (h Handler) RParseRecordID(r *http.Request, name string) *int64 {
	val := h.ParseRecordID(r, name)
	return &val
}
