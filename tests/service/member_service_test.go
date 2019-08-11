package service

import (
	"go2o/core/domain/interface/member"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/tests/ti"
	"testing"
)

var _ = ti.Factory.GetAdRepo()

func TestPagingIntegralLog(t *testing.T) {
	params := &ttype.SPagingParams{
		Opt:        nil,
		OrderField: "",
		OrderDesc:  false,
		Begin:      0,
		Over:       10,
	}
	r, _ := rsi.MemberService.PagingAccountLog(thrift.Context, 1, member.AccountWallet, params)
	t.Logf("%#v", r)
}

func TestPagingWalletLog(t *testing.T) {
	memberId := 77153
	params := &ttype.SPagingParams{
		Opt:        nil,
		OrderField: "",
		OrderDesc:  false,
		Begin:      0,
		Over:       10,
	}
	r, _ := rsi.MemberService.PagingAccountLog(thrift.Context, int64(memberId), member.AccountWallet, params)
	t.Logf("%#v", r)
}

func TestCheckTradePwd(t *testing.T) {
	memberId := 22149
	pwd := domain.Md5("123456")
	//r2,_ := rsi.MemberService.ModifyTradePwd(thrift.Context,int64(memberId),"",pwd)
	//t.Logf("%#v", r2)

	r, _ := rsi.MemberService.CheckTradePwd(thrift.Context, int64(memberId), pwd)
	t.Logf("%#v", r)
}