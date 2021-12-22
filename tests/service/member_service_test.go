package service

import (
	"context"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/tests/ti"
	"testing"
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
	t.Log("token is:", token.AccessToken, " user code:", token.UserCode)
	accessToken, _ := s.CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken: token.AccessToken,
		ExpiresTime: 0,
	})
	if accessToken.MemberId != memberId {
		t.Error(accessToken.Error)
		t.Failed()
	}
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

func TestCheckTradePassword(t *testing.T) {
	memberId := 22149
	pwd := domain.Md5("123456")
	//r2,_ := impl.MemberService.ModifyTradePassword(context.TODO(),int64(memberId),"",pwd)
	//t.Logf("%#v", r2)

	r, _ := impl.MemberService.VerifyTradePassword(context.TODO(),
		&proto.VerifyPasswordRequest{
			MemberId: int64(memberId),
			Password: pwd,
		})
	t.Logf("%#v", r)
}

func TestGetMember(t *testing.T) {
	memberId := 22149
	r, _ := impl.MemberService.GetMember(context.TODO(), &proto.Int64{Value: int64(memberId)})
	t.Logf("%#v", r)
}
