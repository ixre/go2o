package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/crypto"
)

// 测试检查交易密码
func TestCheckTradePassword(t *testing.T) {
	memberId := 1
	pwd := domain.Md5("123456")
	r2, _ := inject.GetMemberService().ChangeTradePassword(context.TODO(),
		&proto.ChangePasswordRequest{
			MemberId:    int64(memberId),
			NewPassword: pwd,
		})
	t.Logf("%#v", r2)

	r, _ := inject.GetMemberService().VerifyTradePassword(context.TODO(),
		&proto.VerifyPasswordRequest{
			MemberId: int64(memberId),
			Password: pwd,
		})
	t.Logf("%#v", r)
}

func TestGetMember(t *testing.T) {
	memberId := 22149
	r, _ := inject.GetMemberService().GetMember(context.TODO(), &proto.MemberIdRequest{MemberId: int64(memberId)})
	t.Logf("%#v", r)
}

func TestChangeProfilePhoto(t *testing.T) {
	r, _ := inject.GetMemberService().ChangeProfilePhoto(context.TODO(),
		&proto.ChangeProfilePhotoRequest{
			MemberId:        702,
			ProfilePhotoUrl: "",
		})
	if r.Code > 0 {
		t.Log(r.Message)
		t.FailNow()
	}
}

func TestChangeMemberLevel(t *testing.T) {
	r, _ := inject.GetMemberService().ChangeLevel(context.TODO(),
		&proto.ChangeLevelRequest{
			MemberId:       702,
			Level:          0,
			LevelCode:      "agent1",
			Review:         false,
			PaymentOrderId: 0,
		})
	if r.Code > 0 {
		t.Log(r.Message)
		t.FailNow()
	}
}
func TestChangeUsername(t *testing.T) {
	ret, _ := inject.GetMemberService().ChangeUsername(context.TODO(), &proto.ChangeUsernameRequest{
		MemberId: 729,
		Username: "哈哈",
	})
	if ret.Code > 0 {
		t.Log(ret.Message)
	}
}
func TestChangePasswordAndCheckLogin(t *testing.T) {
	pwd := domain.Md5("123456")
	t.Log("md5=", pwd)
	r, _ := inject.GetMemberService().ChangePassword(context.TODO(),
		&proto.ChangePasswordRequest{
			MemberId:    1,
			NewPassword: pwd,
		})
	if r.Code > 0 {
		t.Error(r.Message)
	}
	ret, _ := inject.GetMemberService().CheckLogin(context.TODO(), &proto.LoginRequest{
		Username: "13162222872",
		Password: pwd,
	})
	if ret.Code > 0 {
		t.Error(ret.Message)
	}
}

func TestCheckUserLogin(t *testing.T) {
	ret, _ := inject.GetMemberService().CheckLogin(context.TODO(), &proto.LoginRequest{
		Username: "13162222872",
		Password: "14e1b600b1fd579f47433b88e8d85291",
	})
	if ret.Code > 0 {
		t.Error(ret.Message)
	}
}

func Test_memberService_InviterArray(t *testing.T) {
	ret, _ := inject.GetMemberService().InviterArray(context.TODO(), &proto.DepthRequest{
		MemberId: 710,
		Depth:    2,
	})
	t.Log(ret)
}

// 测试绑定邀请人
func TestSetInviter(t *testing.T) {
	ret, _ := inject.GetMemberService().SetInviter(context.TODO(),
		&proto.SetInviterRequest{
			MemberId:    771,
			InviterCode: "f2PWIo",
			AllowChange: false,
		})
	t.Log(ret)
}

// 测试注册商户员工
func TestRegisterMerchantStaff(t *testing.T) {
	ms := inject.GetMemberService()
	phone := "13900000001"
	ret, _ := ms.Register(context.TODO(), &proto.RegisterMemberRequest{
		Username:    phone,
		Password:    crypto.Md5([]byte("123456")),
		Nickname:    "一号律师",
		Phone:       phone,
		Email:       fmt.Sprintf("%s@qq.com", phone),
		InviterCode: "",
		UserType:    2,
		Ext: map[string]string{
			"mchId": "1",
		},
	})
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	} else {
		t.Log("ok")
	}
}

// 测试提交充值订单
func TestSubmitRechargeOrd(t *testing.T) {
	ret, _ := inject.GetMemberService().SubmitRechargePaymentOrder(context.TODO(),
		&proto.SubmitRechargePaymentOrderRequest{
			MemberId: 702,
			Amount:   10000,
		})
	if ret.Code != 0 {
		t.Error(ret.Message)
		t.FailNow()
	} else {
		t.Log("ok:", ret.OrderNo)
	}
}
