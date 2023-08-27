package resp

import (
	"sync"
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

var respPool = sync.Pool{
	New: func() any {
		return &Resp{}
	},
}

func Response(code int, msg string, data any) Resp {
	resp := respPool.Get().(*Resp)
	defer resp.put()
	resp.Code = code
	resp.Msg = msg
	resp.Data = data

	return *resp
}

func (resp Resp) put() {
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	respPool.Put(&resp)
}
