package domain

import (
	"testing"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/crypto"
)

func TestModifyMemberPwd(t *testing.T) {
	m := ti.Factory.GetMemberRepo().GetMember(699)
	md5 := crypto.Md5([]byte("123456"))
	pwd := domain.Sha1Pwd(md5, m.GetValue().Salt)
	// 7c4a8d09ca3762af61e59520943dc26494f8941b
	err := m.Profile().ModifyPassword(pwd, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
