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

func Fail(code int, msg string) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg

	return *resp
}

func Resp(err error) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	errno := unpack(err)
	resp.Code = int(errno.Code())
	resp.Msg = errno.Message()

	return *resp
}

func RespWithData(data any, err error) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	errno := unpack(err)
	resp.Code = int(errno.Code())
	resp.Msg = errno.Message()
	resp.Data = data

	return *resp
}

func RespWithCode(code int, msg string, data any) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return *resp
}

func unpack(err error) *status.Status {
	if err != nil {
		return status.New(http.StatusOK, "ok")
	}
	if st, ok := status.FromError(err); ok {
		return st
	}

	return status.New(http.StatusInternalServerError, "Internal Server Error")
}

// func Resp(data any, err error) response {
// 	if err == nil {
// 		// httpx.OkJsonCtx(r.Context(), w, resp.OkWithData(http.StatusOK, "ok", response))
// 		return Ok()
// 	}

// 	code := int(http.StatusInternalServerError)
// 	msg := err.Error()
// 	if st, ok := status.FromError(err); ok {
// 		code = int(st.Code())
// 		msg = st.Message()
// 	}
// }

func (resp response) put() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(&resp)
}
