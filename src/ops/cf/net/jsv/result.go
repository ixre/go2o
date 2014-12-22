package jsv

import (
	"encoding/json"
)

const (
	RST_ERROR           = `{"result":false}`
	C_OK                = 1    //成功
	C_INVALID_REQUEST   = 1001 //非法请求
	C_PERMISSION_DENIED = 1002 //无权限
	C_PARAM_MISSING     = 1003 //丢失参数
	C_PARAM_ERROR       = 1004 //参数错误
)

var (
	OK                = Result{Result: true, Code: C_OK}
	INVALID_REQUEST   = Result{Code: C_INVALID_REQUEST, Message: "Invalid Request"}
	PERMISSION_DENIED = Result{Code: C_PERMISSION_DENIED}
	PARAM_MISSING     = Result{Code: C_PARAM_MISSING}
	PARAM_ERROR       = Result{Code: C_PARAM_ERROR}
)

// json-rpc无法在客户端反序化为实体，
// 所以data需要返回字符串
//输出结果
type Result struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (t Result) Marshal() []byte {
	json, _ := json.Marshal(t)
	return json
}
