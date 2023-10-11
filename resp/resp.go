package resp

import (
	"net/http"
	"sync"

	"google.golang.org/grpc/status"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

var respPool = sync.Pool{
	New: func() any {
		return &response{}
	},
}

func NewResp() *response {
	return respPool.Get().(*response)
}

func (r *response) RespNoData(err error) *response {
	r.pack(err)
	return r
}

func (r *response) RespWithData(data any, err error) *response {
	r.pack(err)
	r.Data = data
	return r
}

func RespWithCode(code int, msg string, data any) *response {
	resp := respPool.Get().(*response)
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return resp
}

func (r *response) pack(err error) {
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

func (resp *response) Release() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(resp)
}
