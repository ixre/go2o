/**
 * Copyright 2015 @ z3q.net.
 * name : format
 * author : jarryliu
 * date : 2016-05-23 19:42
 * description :
 * history :
 */
package format

import (
	"encoding/json"
	"github.com/jsix/gof/log"
	"html/template"
)

// 强制序列化为可用于HTML的JSON
func MustHtmlJson(v interface{}) template.JS {
	d, err := json.Marshal(v)
	if err != nil {
		log.Println("[ Go2o][ Json] - ", err.Error())
	}
	return template.JS(d)
}
