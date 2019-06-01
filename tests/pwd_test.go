package tests

import (
	"github.com/ixre/gof/crypto"
	"go2o/core/infrastructure/domain"
	"testing"
)

// 管理员密码
func TestMasterPwd(t *testing.T) {
	user := "master"
	pwd := "123456"
	sha1 := domain.Sha1(domain.Md5(pwd) + user + domain.Sha1OffSet)
	t.Log(sha1)
}

func TestMasterPwd2(t *testing.T) {
	user := "master"
	pwd := "fs888888@txxfmall"
	sha1 := crypto.Sha1([]byte(
		pwd + domain.Sha1OffSet))
	encPwd := domain.Md5Pwd(sha1, user)
	t.Log(sha1)
	t.Log(domain.Sha1OffSet)
	t.Log(encPwd)
}

func TestMemberPwd(t *testing.T) {
	pwd := domain.Md5("594488")
	t.Log("--pwd=", pwd, "\n")
	pwd = domain.Sha1(pwd)
	t.Log("--pwd=", pwd, "\n")
}

// 商户密码
func TestMerchantPwd(t *testing.T) {
	user := "zy"
	pwd := "123456"
	encPwd := domain.MerchantSha1Pwd(user, pwd)
	t.Log(encPwd)
}
