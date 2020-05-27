/**
 * Copyright 2015 @ at3.net.
 * name : helper
 * author : jarryliu
 * date : 2016-11-29 17:04
 * description :
 * history :
 */
package parser

import (
	"encoding/json"
	"go2o/core/service/auto_gen/rpc/ttype"
)

func PagingResult(total int, data interface{}, err error) *ttype.SPagingResult_ {
	r := &ttype.SPagingResult_{}
	if err == nil {
		r.Count = int32(total)
		if data == nil || data == "" {
			r.Data = "[]"
		} else {
			d, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}
			r.Data = string(d)
		}
	} else {
		r.ErrCode = 1
		r.ErrMsg = err.Error()
	}
	return r
}
