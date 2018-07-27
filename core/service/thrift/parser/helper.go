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
	"github.com/jsix/gof/util"
	"go2o/core/service/auto-gen/thrift/define"
)

func Result_(data interface{}, err error) *ttype.Result_ {
	r := &ttype.Result_{}
	if err != nil {
		r.ErrCode = 1
		r.ErrMsg = err.Error()
	} else {
		r.ErrCode = 0
		if data != nil {
			switch data.(type) {
			case string, int, int32, int64, bool, float32, float64:
				r.Data1 = util.Str(data)
			default:
				d, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}
				r.Data1 = string(d)
			}
		}
	}
	return r
}

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

func PagingResult(total int, data interface{}, err error) *define.PagingResult_ {
	r := &define.PagingResult_{}
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
