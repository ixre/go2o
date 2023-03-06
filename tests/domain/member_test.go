package domain

import (
	"log"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/types/typeconv"
)

func TestGetMember(t *testing.T) {
	var memberId int64 = 702
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(memberId)
	if m == nil {
		t.FailNow()
	}
}

func TestModifyMemberPwd(t *testing.T) {
	m := ti.Factory.GetMemberRepo().GetMember(702)
	md5 := crypto.Md5([]byte("1234567"))
	pwd := domain.Sha1Pwd(md5, m.GetValue().Salt)
	// 7c4a8d09ca3762af61e59520943dc26494f8941b
	err := m.Profile().ChangePassword(pwd, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestGetInviteUsers(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	iv := repo.CreateMember(&member.Member{Id: 0}).Invitation()
	total, rows := iv.GetInvitationMembers(0, 10)
	t.Log(total, typeconv.MustJson(rows))
}

func TestChangeMemberPhone(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(702)
	err := m.Profile().ChangePhone("18626999822")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试更改用户名
func TestChangeMemberUsername(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(729)
	err := m.ChangeUsername("哈哈")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestQueryMemberInviterArray(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.CreateMember(&member.Member{Id: 719})
	rl := m.GetRelation()
	log.Println("relation=", rl)

	arr := m.Invitation().InviterArray(719, 3)
	log.Println(arr)
}

func TestMemberSaveDefaultAddress(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
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
