package tests

import (
	"testing"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/crypto"
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
	sha1 := domain.Sha1(domain.Md5(pwd) + user + domain.Sha1OffSet)
	t.Log(sha1)
	t.Log(domain.Sha1OffSet)
}

func TestMemberPwd(t *testing.T) {
	pwd := domain.Md5("594488")
	t.Log("--pwd=", pwd, "\n")
	pwd = domain.Sha1(pwd)
	t.Log("--pwd=", pwd, "\n")
}

// 商户密码
func TestMerchantPwd(t *testing.T) {
	//user := "zy"
	pwd := "123456"
	salt := ""
	encPwd := domain.MerchantSha1Pwd(domain.Md5(pwd), salt)
	t.Log(encPwd)
}

func TestMd5Sign(t *testing.T) {
	sign := `{"classId":2,"wip":"LSP:S9240927001427853","subject":"需要","content":"","hopeDesc":"","photoList":"","contactWay":""}1727869234lsp_salt_202406`
	str := crypto.Md5([]byte(sign))
	t.Log(str)
}
