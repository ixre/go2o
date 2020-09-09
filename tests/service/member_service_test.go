package service

import (
	"go2o/core/domain/interface/member"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/impl"
	"go2o/core/service/proto"
	"go2o/core/service/thrift"
	"go2o/core/service/thrift/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/rsi"
	"go2o/tests/ti"
	"testing"
)

var _ = ti.Factory.GetAdRepo()

func TestPagingIntegralLog(t *testing.T) {
	params := &proto.SPagingParams{
		Opt:        nil,
		OrderField: "",
		OrderDesc:  false,
		Begin:      0,
		Over:       10,
	}
	r, _ := impl.MemberService.PagingAccountLog(context.TODO(), 1, member.AccountWallet, params)
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
	r, _ := impl.MemberService.PagingAccountLog(context.TODO(), int64(memberId), member.AccountWallet, params)
	t.Logf("%#v", r)
}

func TestCheckTradePwd(t *testing.T) {
	memberId := 22149
	pwd := domain.Md5("123456")
	//r2,_ := impl.MemberService.ModifyTradePwd(context.TODO(),int64(memberId),"",pwd)
	//t.Logf("%#v", r2)

	r, _ := impl.MemberService.CheckTradePwd(context.TODO(), int64(memberId), pwd)
	t.Logf("%#v", r)
}

func TestGetMember(t *testing.T) {
	memberId := 22149
	r, _ := impl.MemberService.GetMember(context.TODO(), &proto.Int64{Value: int64(memberId)})
	t.Logf("%#v", r)
}
