// Code generated by "stringer -type State -linecomment"; DO NOT EDIT.

package errorx

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UNPACK_TOKEN_ERROR-1002]
	_ = x[MISS_PARAMS-1003]
	_ = x[REPETITION_SUBMIT-1004]
	_ = x[UNAME_OR_PWD_IS_ERROR-1005]
	_ = x[UPLOAD_FILE_MAX-1005]
}

const _State_name = "非法请求缺少参数重复提交用户名或密码错误"

var _State_index = [...]uint8{0, 12, 24, 36, 60}

func (i State) String() string {
	i -= 1002
	if i < 0 || i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i+1002), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}