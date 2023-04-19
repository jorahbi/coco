package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type StateInterface interface {
	Code() int
	Msg() string
}

func Response(w http.ResponseWriter, resp interface{}, state StateInterface) {
	body := Body{
		Code: state.Code(),
		Msg:  state.Msg(),
		Data: resp,
	}

	httpx.OkJson(w, body)
}
