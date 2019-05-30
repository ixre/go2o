package api

import "go2o/core/service/auto_gen/rpc/ttype"

type apiUtil struct {
}

func (a apiUtil) SResult(err error) *ttype.Result_ {
	if err != nil {
		return &ttype.Result_{ErrCode: 1, ErrMsg: err.Error()}
	}
	return &ttype.Result_{}
}
