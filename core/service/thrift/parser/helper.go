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
	"go2o/gen-code/thrift/define"
)

func Result(data interface{}, err error) *define.Result_ {
	r := &define.Result_{}
	if err != nil {
		r.ErrCode = 1
		r.ErrMsg = err.Error()
	} else {
		r.ErrCode = 0
		if data != nil {
			switch data.(type) {
			case string, int, int32, int64, bool, float32, float64:
				r.Data = util.Str(data)
			default:
				d, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}
				r.Data = string(d)
			}
		}
	}
	return r
}

func Result64(id int64, err error) *define.Result64 {
	r := &define.Result64{}
	if err != nil {
		r.ErrMsg = err.Error()
	} else {
		r.Result_ = true
		r.ID = id
	}
	return r
}

func DResult(data float64, err error) *define.DResult_ {
	r := &define.DResult_{}
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
