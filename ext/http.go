package ext

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DecxBase/core/server"
)

type HTTPHandler struct {
	server.Handler
	Prefix    string
	endpoints []*httpEndpoint
}

func (h HTTPHandler) MakePath(path ...string) string {
	paths := []string{"/" + h.Prefix}
	if len(path) > 0 {
		paths = append(paths, path...)
	}

	return strings.Join(paths, "/")
}

func (h *HTTPHandler) MakePattern(method string, paths ...string) string {
	return fmt.Sprintf("%s %s", method, h.MakePath(paths...))
}

func (h HTTPHandler) TransformEOF(err error) error {
	if err.Error() == "EOF" {
		return errors.New("body payload is empty")
	}

	return err
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

func NewHTTPHandler(prefix string) HTTPHandler {
	return HTTPHandler{
		endpoints: make([]*httpEndpoint, 0),
		Prefix:    prefix,
		Handler:   server.Handler{},
	}
}
