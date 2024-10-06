package ext

import (
	"net/http"

	"github.com/DecxBase/core/server"
)

type HTTPHandler struct {
	server.Handler
	endpoints []*httpEndpoint
}

func (h *HTTPHandler) WithEndpoint(endps ...*httpEndpoint) *HTTPHandler {
	h.endpoints = append(h.endpoints, endps...)

	return h
}

func (h HTTPHandler) RegisterRoutes(router *http.ServeMux) error {
	if h.endpoints != nil {
		for _, endp := range h.endpoints {
			router.HandleFunc(endp.Pattern, h.HTTPHandler(endp.Handle))
		}
	}

	return nil
}

func NewHTTPHandler() HTTPHandler {
	return HTTPHandler{
		endpoints: make([]*httpEndpoint, 0),
		Handler:   server.Handler{},
	}
}
