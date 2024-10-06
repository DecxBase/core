package ext

import (
	"fmt"
	"strings"
)

type CrudHTTPListHandler interface {
	FindRecords(ctx HTTPEndpointContext) (any, error)
}

type CrudHTTPFindHandler interface {
	FindRecord(ctx HTTPEndpointContext) (any, error)
}

type CrudHTTPSaveHandler interface {
	SaveRecord(ctx HTTPEndpointContext) (any, error)
}

type CrudHTTPDeleteHandler interface {
	DeleteRecord(ctx HTTPEndpointContext) (any, error)
}

type CrudHTTPHandler struct {
	HTTPHandler
	Prefix string
}

func (h CrudHTTPHandler) MakePath(path ...string) string {
	paths := []string{"/" + h.Prefix}
	if len(path) > 0 {
		paths = append(paths, path...)
	}

	return strings.Join(paths, "/")
}

func (h *CrudHTTPHandler) RegisterList(hn CrudHTTPListHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(fmt.Sprintf("GET %s", h.MakePath()), hn.FindRecords),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterFind(hn CrudHTTPFindHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(fmt.Sprintf("GET %s", h.MakePath("{id}")), hn.FindRecord),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterSave(hn CrudHTTPSaveHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(fmt.Sprintf("POST %s", h.MakePath()), hn.SaveRecord),
		HTTPEndpoint(fmt.Sprintf("PATCH %s", h.MakePath("{id}")), hn.SaveRecord),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterDelete(hn CrudHTTPDeleteHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(fmt.Sprintf("DELETE %s", h.MakePath("{id}")), hn.DeleteRecord),
	)

	return h
}

func NewCrudHTTPHandler(prefix string) CrudHTTPHandler {
	return CrudHTTPHandler{
		Prefix:      prefix,
		HTTPHandler: NewHTTPHandler(),
	}
}
