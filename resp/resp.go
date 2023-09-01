package resp

import (
	"sync"
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

func Ok(code int, msg string) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg

	return *resp
}

func OkWithData(code int, msg string, data any) response {
	resp := respPool.Get().(*response)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return *resp
}

func (resp response) put() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(&resp)
}
