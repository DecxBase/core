package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/DecxBase/core/utils"
)

type ServerException struct {
	err error
	msg string

	Code        string
	StatusCode  int
	Reportable  bool
	Context     string
	ContextData map[string]any
}

func (e *ServerException) WithReportable(r bool) *ServerException {
	e.Reportable = r
	return e
}

func (e *ServerException) WithContext(c string) *ServerException {
	e.Context = c
	return e
}

func (e *ServerException) WithContextData(data map[string]any) *ServerException {
	e.ContextData = data
	return e
}

func (e *ServerException) WithCode(code string) *ServerException {
	e.Code = code
	return e
}

func (e *ServerException) WithStatusCode(code int) *ServerException {
	e.StatusCode = code
	return e
}

func (e *ServerException) WithMessage(msg string) *ServerException {
	e.msg = msg
	return e
}

func (e ServerException) ErrorMsg() string {
	msg := e.err.Error()
	if len(e.msg) > 0 {
		msg = e.msg
	}

	return msg
}

func (e *ServerException) Error() string {
	tp := reflect.ValueOf(e).Elem().Type()

	return fmt.Sprintf("[%s] %s [%d:%s]", tp, e.ErrorMsg(), e.StatusCode, e.Code)
}

func (e *ServerException) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"error": e.ErrorMsg(),
		"code":  e.Code,
	})
}

func (e *ServerException) WriteToResponse(w http.ResponseWriter) error {
	return utils.WriteResponseJSON(w, e.StatusCode, e)
}

func Raise(err error) *ServerException {
	return &ServerException{
		err:        err,
		StatusCode: http.StatusInternalServerError,
	}
}
