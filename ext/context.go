package ext

import (
	"context"
	"net/http"

	"github.com/DecxBase/core/utils"
)

type HTTPEndpointContext struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Data    map[string]any
}

func (c HTTPEndpointContext) Context() context.Context {
	return c.Request.Context()
}

func (c HTTPEndpointContext) Method() string {
	return c.Request.Method
}

func (c HTTPEndpointContext) FromQuery(target any) error {
	return utils.ParseQuery(c.Request, target)
}

func (c HTTPEndpointContext) FromBody(target any) error {
	return utils.ParseBody(c.Request, target)
}

func (c HTTPEndpointContext) PathValue(name string) string {
	return c.Request.PathValue(name)
}

func (c HTTPEndpointContext) PathID(name string) int64 {
	return utils.StringToInt64(c.PathValue(name))
}

func (c HTTPEndpointContext) RPathID(name string) *int64 {
	val := c.PathID(name)
	return &val
}

func (c HTTPEndpointContext) WriteError(err error) {
	c.WriteStatusError(http.StatusInternalServerError, err)
}

func (c HTTPEndpointContext) WriteStatusError(statusCode int, err error) {
	utils.WriteResponseError(c.Writer, statusCode, err)
}

func (c HTTPEndpointContext) WriteData(data any) {
	utils.WriteResponseData(c.Writer, data)
}

type HTTPEndpointHandler = func(HTTPEndpointContext) (any, error)
