package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

//go:generate stringer -type State -linecomment
const (
	SUCCESS State = 0  //请求成功
	FAIL    State = -1 //请求失败
)

type State int

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

func (e State) Code() int {
	return int(e)
}
