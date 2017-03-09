/**
 * Copyright 2015 @ at3.net.
 * name : helper
 * author : jarryliu
 * date : 2016-11-29 17:04
 * description :
 * history :
 */
package parser

import "go2o/core/service/thrift/idl/gen-go/define"

func Result(id int32, err error) *define.Result_ {
	r := &define.Result_{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
		r.ID = id
	}
	return r
}

func Result64(id int64, err error) *define.Result64 {
	r := &define.Result64{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
		r.ID = id
	}
	return r
}

func I64Result(id int64, err error) *define.Result_ {
	r := &define.Result_{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
		r.ID = int32(id)
	}
	return r
}

func DResult(data float64, err error) *define.DResult_ {
	r := &define.DResult_{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
		r.Data = data
	}
	return r
}
