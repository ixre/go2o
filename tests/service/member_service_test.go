package service

import (
	"context"
	"testing"
	"time"

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
		MemberId:    memberId,
		ExpiresTime: time.Now().Unix() + 720,
	})
	if len(token.Error) > 0 {
		t.Error(token.Error)
		t.Failed()
	}
	t.Log("token is:", token.AccessToken)
	now := time.Now().Unix()
	token.AccessToken = "Bearer " // + token.AccessToken
	accessToken, _ := s.CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken:      token.AccessToken,
		CheckExpireTime:  now + 800,
		RenewExpiresTime: now + 900,
	})
	if accessToken.MemberId != memberId {
		t.Error(accessToken.Error)
		t.Failed()
	}
}

func TestCheckMemberAccessToken(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI3MTAiLCJleHAiOjE2Nzc3NTM5MjAsImlzcyI6ImdvMm8iLCJzdWIiOiJnbzJvLWFwaS1qd3QifQ.Vrd0NuAT5AfKM-C5NespEyhAiyMVDugoKDDgwv5hr_g"
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

// 测试检查交易密码
func TestCheckTradePassword(t *testing.T) {
	memberId := 1
	pwd := domain.Md5("123456")
	r2, _ := impl.MemberService.ChangeTradePassword(context.TODO(),
		&proto.ChangePasswordRequest{
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

func TestChangeMemberLevel(t *testing.T) {
	r, _ := impl.MemberService.ChangeLevel(context.TODO(),
		&proto.ChangeLevelRequest{
			MemberId:       702,
			Level:          0,
			LevelCode:      "agent1",
			Review:         false,
			PaymentOrderId: 0,
		})
	if r.ErrCode > 0 {
		t.Log(r.ErrMsg)
		t.FailNow()
	}
}
func TestChangeUsername(t *testing.T) {
	ret, _ := impl.MemberService.ChangeUsername(context.TODO(), &proto.ChangeUsernameRequest{
		MemberId: 702,
		Username: "18924140900",
	})
	if ret.ErrCode > 0 {
		t.Log(ret.ErrMsg)
	}

}
func TestChangePasswordAndCheckLogin(t *testing.T) {
	pwd := domain.Md5("123456")
	t.Log("md5=", pwd)
	r, _ := impl.MemberService.ChangePassword(context.TODO(),
		&proto.ChangePasswordRequest{
			MemberId:    1,
			NewPassword: pwd,
		})
	if r.ErrCode > 0 {
		t.Error(r.ErrMsg)
	}
	ret, _ := impl.MemberService.CheckLogin(context.TODO(), &proto.LoginRequest{
		Username: "13162222872",
		Password: pwd,
	})
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	}
}

func TestCheckUserLogin(t *testing.T) {
	ret, _ := impl.MemberService.CheckLogin(context.TODO(), &proto.LoginRequest{
		Username: "13162222872",
		Password: "14e1b600b1fd579f47433b88e8d85291",
	})
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	}
}

func Test_memberService_InviterArray(t *testing.T) {
	ret, _ := impl.MemberService.InviterArray(context.TODO(), &proto.DepthRequest{
		MemberId: 710,
		Depth:    2,
	})
	t.Log(ret)
}
