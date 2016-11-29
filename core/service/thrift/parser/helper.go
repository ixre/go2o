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

func Result(err error) *define.Result_ {
	r := &define.Result_{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
	}
	return r
}

func DResult(err error) *define.DResult_ {
	r := &define.DResult_{}
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result_ = true
	}
	return r
}
