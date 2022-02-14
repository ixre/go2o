package tests

import (
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/tests/ti"
	"testing"
	"time"
)

func TestCreateNewMember(t *testing.T) {
	inviteCode := ""
	phone := "13162222820"
	inviterId := 22149
	repo := ti.Factory.GetMemberRepo()
	_, err := repo.GetManager().CheckInviteRegister(inviteCode)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v := &member.Member{
		User:   phone,
		Pwd:    domain.Md5("123456"),
		Avatar: "",
		Phone:  phone,
		Email:  "",
		Flag:   0,
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
	repo := ti.Factory.GetMemberRepo()
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
	repo := ti.Factory.GetMemberRepo()
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

func TestGetMember(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(502)
	t.Logf("%#v", m.GetValue())
}

func TestModifyPassword(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(2)
	NewPassword := domain.MemberSha1Pwd(domain.Md5("13268240456"), m.GetValue().Salt)
	err := m.Profile().ModifyPassword(NewPassword, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if o := m.GetValue().Pwd; o != NewPassword {
		t.Logf("登陆密码不正确")
		t.FailNow()
	}
}

func TestReceiptsCode(t *testing.T) {
	memberId := 22149
	m := ti.Factory.GetMemberRepo().GetMember(int64(memberId))
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
	m := ti.Factory.GetMemberRepo().GetMember(int64(memberId))
	err := m.Lock(1440, "测试锁定会员")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	time.Sleep(time.Second * 2)
}

// 测试更改邀请人
func TestUpdateInviter(t *testing.T) {
	memberId := 96941
	inviterId := 96949
	//memberId := 97622
	//inviterId := 97613  13630000873
	//inviterId = 96999
	//  18824894128
	m := ti.Factory.GetMemberRepo().GetMember(int64(memberId))
	err := m.BindInviter(int64(inviterId), true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试钱包
func TestMemberWallet(t *testing.T) {
	var memberId int64 = 16
	m := ti.Factory.GetMemberRepo().GetMember(memberId)
	ic := m.GetAccount()
	if ic.GetValue().WalletBalance != ic.Wallet().Get().Balance {
		t.Error("钱包金额不符合")
		t.FailNow()
	}
}
