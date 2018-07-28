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



func Result64(id int64, err error) *ttype.Result64 {
	r := &ttype.Result64{}
	if err != nil {
		r.ErrMsg = err.Error()
	} else {
		r.Result_ = true
		r.ID = id
	}
	return r
}

func DResult(data float64, err error) *ttype.DResult_ {
	r := &ttype.DResult_{}
	if err != nil {
		r.ErrMsg = err.Error()
	} else {
		r.Result_ = true
		r.Data = data
	}
	return r
}

func PagingResult(total int, data interface{}, err error) *ttype.PagingResult_ {
	r := &ttype.PagingResult_{}
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
