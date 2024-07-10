package domain

import (
	"log"
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/inject"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/typeconv"
)

func TestGetMember(t *testing.T) {
	var memberId int64 = 702
	repo := inject.GetMemberRepo()
	m := repo.GetMember(memberId)
	if m == nil {
		t.FailNow()
	}
	t.Logf("%#v", m.GetValue())
}

// 测试修改用户信息
func TestModifyMemberProfile(t *testing.T) {
	m := inject.GetMemberRepo().GetMember(1)
	v := m.GetValue()
	v.Nickname = "哈哈"
	m.Save()
	time.Sleep(30 * time.Second)
}

func TestModifyMemberPwd(t *testing.T) {
	m := inject.GetMemberRepo().GetMember(31)
	md5 := crypto.Md5([]byte("123456"))
	pwd := domain.Sha1Pwd(md5, m.GetValue().Salt)
	// 7c4a8d09ca3762af61e59520943dc26494f8941b
	err := m.Profile().ChangePassword(pwd, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestGetInviteUsers(t *testing.T) {
	repo := inject.GetMemberRepo()
	iv := repo.CreateMember(&member.Member{Id: 0}).Invitation()
	total, rows := iv.GetInvitationMembers(0, 10)
	t.Log(total, typeconv.MustJson(rows))
}

func TestChangeMemberPhone(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.GetMember(702)
	err := m.Profile().ChangePhone("18626999822")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试更改用户名
func TestChangeMemberUsername(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.GetMember(729)
	err := m.ChangeUsername("哈哈")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestQueryMemberInviterArray(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.CreateMember(&member.Member{Id: 719})
	rl := m.GetRelation()
	log.Println("relation=", rl)

	arr := m.Invitation().InviterArray(719, 3)
	log.Println(arr)
}

func TestMemberSaveDefaultAddress(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.CreateMember(&member.Member{Id: 723})
	addrList := m.Profile().GetDeliverAddress()
	addr := addrList[0]
	av := addr.GetValue()
	err := addr.SetValue(&av)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = m.Profile().SetDefaultAddress(addr.GetDomainId())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestCreateNewMember(t *testing.T) {
	inviteCode := ""
	phone := "13162222821"
	inviterId := 6
	repo := inject.GetMemberRepo()
	_, err := repo.GetManager().CheckInviteRegister(inviteCode)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v := &member.Member{
		Username: phone,
		Password: domain.Md5("123456"),
		Portrait: "",
		Phone:    phone,
		Email:    "",
		RoleFlag: member.RoleMchStaff,
	}
	m := repo.CreateMember(v) //创建会员
	id, err := m.Save()
	if err == nil {
		err = m.BindInviter(int64(inviterId), true)
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(5 * time.Second)
	t.Logf("注册成功,Id:%d", id)
}

func TestSaveMemberGroups(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.GetManager()
	groups := m.GetAllBuyerGroups()
	oriName := groups[0].Name
	groups[0].Name = "测试修改"
	_, err := m.SaveBuyerGroup(groups[0])
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	groups[0].Name = oriName
	_, err = m.SaveBuyerGroup(groups[0])
	if err != nil {
		t.Error(err)
		return
	}
	v := m.GetBuyerGroup(groups[0].ID)
	if v.Name != oriName {
		t.Log("旧名称：", oriName, "; 当前名称:", v.Name)
	}
}

func TestToBePremium(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.GetMember(1)
	err := m.Premium(member.PremiumWhiteGold,
		time.Now().Add(time.Hour*24*365).Unix())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	m = repo.GetMember(m.GetAggregateRootId())
	v := m.GetValue()
	t.Logf("Premium: user:%d ; expires:%s", v.PremiumUser,
		time.Unix(v.PremiumExpires, 0).Format("2006-01-02 15:04:05"))
}

func TestChangePassword(t *testing.T) {
	repo := inject.GetMemberRepo()
	m := repo.GetMember(2)
	NewPassword := domain.MemberSha1Pwd(domain.Md5("13268240456"),
		m.GetValue().Salt)
	err := m.Profile().ChangePassword(NewPassword, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if o := m.GetValue().Password; o != NewPassword {
		t.Logf("登陆密码不正确")
		t.FailNow()
	}
}

func TestReceiptsCode(t *testing.T) {
	memberId := 22149
	m := inject.GetMemberRepo().GetMember(int64(memberId))
	err := m.Profile().SaveReceiptsCode(&member.ReceiptsCode{
		Identity:  "alipay",
		Name:      "刘铭",
		AccountId: "jarrysix#gmail.com",
		CodeUrl:   "1.jpg",
		State:     1,
	})
	t.Log("err:", err)
	err = m.Profile().SaveReceiptsCode(&member.ReceiptsCode{
		Id:        2,
		Identity:  "unipay",
		Name:      "刘铭",
		AccountId: "jarrysix",
		CodeUrl:   "1.jpg",
		State:     1,
	})
	err = m.Profile().SaveReceiptsCode(&member.ReceiptsCode{
		Identity:  "wepay",
		Name:      "刘铭",
		AccountId: "jarrysix",
		CodeUrl:   "1.jpg",
		State:     1,
	})
	t.Log("err:", err)
}

func TestLogin(t *testing.T) {
	pwd := "d682a6db237d3fe29f07a1545778ecf3"
	t.Log(len(pwd))
	flag := 133
	b := flag&member.FlagLocked == member.FlagLocked
	t.Log("--", b)
}

// 测试锁定会员
func TestLockMember(t *testing.T) {
	memberId := 97839
	m := inject.GetMemberRepo().GetMember(int64(memberId))
	err := m.Lock(1440, "测试锁定会员")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	time.Sleep(time.Second * 2)
}

// 测试更改邀请人
func TestUpdateInviter(t *testing.T) {
	memberId := 728
	inviterId := 710
	m := inject.GetMemberRepo().GetMember(int64(memberId))
	err := m.BindInviter(int64(inviterId), false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试钱包
func TestMemberWallet(t *testing.T) {
	var memberId int64 = 16
	m := inject.GetMemberRepo().GetMember(memberId)
	ic := m.GetAccount()
	if ic.GetValue().WalletBalance != ic.Wallet().Get().Balance {
		t.Error("钱包金额不符合")
		t.FailNow()
	}
}

// 测试更改头像
func TestChangeHeadPortrait(t *testing.T) {
	var memberId int64 = 723
	portraitUrl := "a/20230310144156396.jpeg"
	m := inject.GetMemberRepo().GetMember(memberId)
	err := m.Profile().ChangeHeadPortrait(portraitUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 　测试更改等级
func TestChangeMemberLevel(t *testing.T) {
	memberId := 821
	repo := inject.GetMemberRepo()
	m := repo.GetMember(int64(memberId))
	err := m.ChangeLevel(1, 0, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// TestReviewCertification 测试审核会员认证信息
func TestReviewCertification(t *testing.T) {
	memberId := 854
	repo := inject.GetMemberRepo()
	m := repo.GetMember(int64(memberId))
	err := m.Profile().ReviewCertification(true, "测试审核")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试钱包
func TestMemberWalletRefreeze(t *testing.T) {
	var memberId int64 = 848
	m := inject.GetMemberRepo().GetMember(memberId)
	ic := m.GetAccount()
	err := ic.Charge(member.AccountWallet, "测试充值", 10000, "", "测试充值")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tradeLogId, err := ic.Freeze(member.AccountWallet, member.AccountOperateData{
		Title:      "测试冻结",
		Amount:     10,
		OuterNo:    "xxx",
		Remark:     "",
		TradeLogId: 0,
	}, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	wr := inject.GetWalletRepo()
	l := wr.GetWalletLog_(tradeLogId)
	t.Logf("冻结金额:%d, 钱包余额:%d", l.ChangeValue, l.Balance)

	tradeLogId, _ = ic.Freeze(member.AccountWallet, member.AccountOperateData{
		Title:      "测试冻结",
		Amount:     20,
		OuterNo:    "xxx",
		Remark:     "",
		TradeLogId: tradeLogId,
	}, 0)
	l = wr.GetWalletLog_(tradeLogId)
	t.Logf("冻结金额:%d, 钱包余额:%d", l.ChangeValue, l.Balance)
}
