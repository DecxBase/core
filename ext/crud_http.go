package ext

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
}

func (h *CrudHTTPHandler) RegisterList(hn CrudHTTPListHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(h.MakePattern("GET"), hn.FindRecords),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterFind(hn CrudHTTPFindHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(h.MakePattern("GET", "{id}"), hn.FindRecord),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterSave(hn CrudHTTPSaveHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(h.MakePattern("POST", "{id}"), hn.SaveRecord),
		HTTPEndpoint(h.MakePattern("PATCH", "{id}"), hn.SaveRecord),
	)

	return h
}

func (h *CrudHTTPHandler) RegisterDelete(hn CrudHTTPDeleteHandler) *CrudHTTPHandler {
	h.WithEndpoint(
		HTTPEndpoint(h.MakePattern("DELETE", "{id}"), hn.DeleteRecord),
	)

	return h
}

func NewCrudHTTPHandler(prefix string) CrudHTTPHandler {
	return CrudHTTPHandler{
		HTTPHandler: NewHTTPHandler(prefix),
	}
}
