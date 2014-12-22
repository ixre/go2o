package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//客户端消息
type ClientMessage struct {
	Result  bool
	Message string
	Data    interface{}
}

//序列化消息
func serializeMessage(result bool, message string, data interface{}) ([]byte, error) {
	clientMsg := new(ClientMessage)
	clientMsg.Result = result
	clientMsg.Message = message
	clientMsg.Data = data
	return json.Marshal(clientMsg)
}

//输出到客户端
func Seria2json(w http.ResponseWriter, result bool, message string, data interface{}) {
	var jsonStr string
	jso, err := serializeMessage(result, message, data)

	if err != nil {
		jsonStr = "{Result:false,Message:'" + err.Error() + "',Data:null}"
	} else {
		jsonStr = string(jso)
	}
	fmt.Fprintf(w, jsonStr)
	w.Header().Set("Content-Type", "application/json")
}
