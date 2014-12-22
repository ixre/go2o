package jsv

import (
	"encoding/json"
	"errors"
	"strings"
)

func LogErr(err error) {
	//log := Context.Log()
	//log.AddDepth(0)
	Context.Log().PrintErr(err)
	//log.ResetDepth()
}

func Println(v ...interface{}) {
	//log := Context.Log()
	//log.AddDepth(0)
	Context.Log().Println(v...)
	//log.ResetDepth()
}

func Printf(s string, v ...interface{}) {
	//log := Context.Log()
	//log.AddDepth(0)
	Context.Log().Printf(s, v...)
	//log.ResetDepth()
}

func MarshalString(e interface{}) string {
	if e != nil {
		js, _ := json.Marshal(e)
		return string(js)
	}
	return ""
}

// 序列化为Json字符串带转义
func MarshalFJ(e interface{}) string {
	if e != nil {
		js, _ := json.Marshal(e)
		return strings.Replace(string(js), `""`, `\"`, -1)
	}
	return ""
}

func UnmarshalMap(in interface{}, to interface{}) error {
	if in != nil {
		js, _ := json.Marshal(in)
		Println("[Client][MAP]:", string(js))
		err := json.Unmarshal(js, &to)
		Println("[Client][MAP-RESULT]:", to)
		return err
	}
	return errors.New("nil point refrence.")
}

func UnmarshalString(s string, e interface{}) error {
	if e != nil {
		return json.Unmarshal([]byte(s), e)
	}
	return errors.New("entity is null.")
}

func Unmarshal(b []byte, e interface{}) error {
	if e != nil {
		err := json.Unmarshal(b, e)
		if debugMode && err != nil {
			Printf("[Codec][Error]:", err, string(b))
		}
		return err
	}
	return errors.New("entity is null.")
}
