package service

import (
	"go2o/core/domain/interface/member"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/tests/ti"
	"testing"
)
var _ = ti.Factory.GetAdRepo()

func TestPagingIntegralLog(t *testing.T){
	params := &ttype.SPagingParams{
		Opt:        nil,
		OrderField: "",
		OrderDesc:  false,
		Begin:      0,
		Over:       10,
	}
	r,_ := rsi.MemberService.PagingAccountLog(thrift.Context,1,member.AccountWallet,params)
	t.Logf("%#v",r)
}
