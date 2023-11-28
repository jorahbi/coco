package resp

import (
	"net/http"

	"google.golang.org/grpc/status"
)

type response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

func NewResp[T any]() *response[T] {
	return &response[T]{}
}

func (r *response[T]) RespNoData(err error) *response[T] {
	r.pack(err)
	return r
}

func (r *response[T]) RespWithData(data T, err error) *response[T] {
	r.pack(err)
	r.Data = data
	return r
}

func (r *response[T]) RespWithCode(code int, msg string, data T) *response[T] {
	r.Code = code
	r.Msg = msg
	r.Data = data

	return r
}

func (r *response[T]) pack(err error) {
	r.Code = http.StatusOK
	r.Msg = "ok"
	if err == nil {
		return
	}
	if st, ok := status.FromError(err); ok {
		r.Code = int(st.Code())
		r.Msg = st.Message()
		return
	}
	r.Code = http.StatusInternalServerError
	r.Msg = err.Error()
}
