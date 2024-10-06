package ext

import "net/http"

type httpEndpoint struct {
	Pattern string
	Data    map[string]any
	handler HTTPEndpointHandler
}

func (e *httpEndpoint) With(key string, value any) *httpEndpoint {
	e.Data[key] = value

	return e
}

func (e httpEndpoint) Handle(w http.ResponseWriter, r *http.Request) (any, error) {
	return e.handler(HTTPEndpointContext{
		Request: r,
		Writer:  w,
		Data:    e.Data,
	})
}

func HTTPEndpoint(pattern string, fn HTTPEndpointHandler) *httpEndpoint {
	return &httpEndpoint{
		Pattern: pattern,
		Data:    make(map[string]any),
		handler: fn,
	}
}
