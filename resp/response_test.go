package resp

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResp(t *testing.T) {
	r := Fail(1, "test")
	s, e := json.Marshal(r)
	fmt.Println(string(s), e)

	r = OkWithData(1, "test", nil)
	s, e = json.Marshal(r)
	fmt.Println(string(s), e)
}
