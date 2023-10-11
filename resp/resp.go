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

func RespNoData(err error) *response {
	return pack(err)

}

func RespWithData(data any, err error) *response {
	rsp := pack(err)
	rsp.Data = data
	return rsp
}

func RespWithCode(code int, msg string, data any) *response {
	resp := respPool.Get().(*response)
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return resp
}

func pack(err error) *response {
	resp := respPool.Get().(*response)
	resp.Code = http.StatusOK
	resp.Msg = "ok"
	if err == nil {
		return resp
	}
	if st, ok := status.FromError(err); ok {
		resp.Code = int(st.Code())
		resp.Msg = st.Message()
		return resp
	}
	resp.Code = http.StatusInternalServerError
	resp.Msg = err.Error()
	return resp
}

func (resp *response) Put() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(resp)
}
