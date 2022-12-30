package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/types/typeconv"
)

var _ = ti.Factory.GetAdRepo()

func TestGrantMemberAccessToken(t *testing.T) {
	var memberId int64 = 1
	s := impl.MemberService
	token, _ := s.GrantAccessToken(context.TODO(), &proto.GrantAccessTokenRequest{
		MemberId: memberId,
		Expire:   720,
	})
	if len(token.Error) > 0 {
		t.Error(token.Error)
		t.Failed()
	}
	t.Log("token is:", token.AccessToken)
	accessToken, _ := s.CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken: token.AccessToken,
	})
	if accessToken.MemberId != memberId {
		t.Error(accessToken.Error)
		t.Failed()
	}
}

func TestCheckMemberAccessToken(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI3MDIiLCJleHAiOjE2NzIzMDM4MDIsImlzcyI6ImdvMm8iLCJzdWIiOiJnbzJvLWFwaS1qd3QifQ.Ebx4PcD0KSIftqejzfbyYbUpunm3jEi0gsScipcl-lo"
	ret, _ := impl.MemberService.CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken: accessToken,
	})
	if len(ret.Error) > 0 {
		t.Log(ret.Error)
		t.FailNow()
	}
	t.Log(typeconv.MustJson(ret))
	t.Log("会员Id", ret.MemberId)
}

func TestPagingIntegralLog(t *testing.T) {
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := impl.MemberService.PagingAccountLog(context.TODO(),
		&proto.PagingAccountInfoRequest{
			MemberId:    1,
			AccountType: int32(member.AccountWallet),
			Params:      params,
		})
	t.Logf("%#v", r)
}

func TestPagingWalletLog(t *testing.T) {
	memberId := 77153
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := impl.MemberService.PagingAccountLog(context.TODO(),
		&proto.PagingAccountInfoRequest{
			MemberId:    int64(memberId),
			AccountType: int32(member.AccountWallet),
			Params:      params,
		})
	t.Logf("%#v", r)
}

// 测试检查交易密码
func TestCheckTradePassword(t *testing.T) {
	memberId := 1
	pwd := domain.Md5("123456")
	r2, _ := impl.MemberService.ModifyTradePassword(context.TODO(),
		&proto.ModifyPasswordRequest{
			MemberId:    int64(memberId),
			NewPassword: pwd,
		})
	t.Logf("%#v", r2)

	r, _ := impl.MemberService.VerifyTradePassword(context.TODO(),
		&proto.VerifyPasswordRequest{
			MemberId: int64(memberId),
			Password: pwd,
		})
	t.Logf("%#v", r)
}

func TestGetMember(t *testing.T) {
	memberId := 22149
	r, _ := impl.MemberService.GetMember(context.TODO(), &proto.MemberIdRequest{MemberId: int64(memberId)})
	t.Logf("%#v", r)
}

func TestChangeHeadPortrait(t *testing.T) {
	r, _ := impl.MemberService.ChangeHeadPortrait(context.TODO(),
		&proto.ChangePortraitRequest{
			MemberId:    702,
			PortraitUrl: "",
		})
	if r.ErrCode > 0 {
		t.Log(r.ErrMsg)
		t.FailNow()
	}
}
