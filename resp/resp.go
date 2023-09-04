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

func RespNoData(err error) response {
	return pack(nil, err)
}

func RespWithData(data any, err error) response {
	return pack(data, err)
}

func RespWithCode(code int, msg string, data any) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return *resp
}

func pack(data any, err error) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = http.StatusOK
	resp.Msg = "ok"
	if err == nil {
		resp.Data = data
		return *resp
	}
	if st, ok := status.FromError(err); ok {
		resp.Code = int(st.Code())
		resp.Msg = st.Message()
		return *resp
	}
	resp.Code = http.StatusInternalServerError
	resp.Msg = err.Error()
	return *resp
}

func (resp response) put() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(&resp)
}
