package tests

import (
	"testing"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/gof/crypto"
)

func init() {
	repo := inject.GetRegistryRepo()
	key, _ := repo.GetValue(registry.SysPrivateKey)
	domain.ConfigPrivateKey(key)
}

// 管理员密码
func TestMasterPwd(t *testing.T) {
	user := "master"
	pwd := "123456"
	sha1 := domain.HmacSha256(domain.Md5(pwd) + user)
	t.Log(sha1)
}

func TestMemberPwd(t *testing.T) {
	pwd := domain.Md5("594488")
	t.Log("--pwd=", pwd, "\n")
	pwd = domain.HmacSha256(pwd)
	t.Log("--pwd=", pwd, "\n")
}

// 商户密码
func TestMerchantPwd(t *testing.T) {
	//user := "zy"
	pwd := "123456"
	salt := ""
	encPwd := domain.MerchantSha265Pwd(domain.Md5(pwd), salt)
	t.Log(encPwd)
}

func TestMd5Sign(t *testing.T) {
	sign := `{"classId":2,"wip":"LSP:S9240927001427853","subject":"需要","content":"","hopeDesc":"","photoList":"","contactWay":""}1727869234lsp_salt_202406`
	str := crypto.Md5([]byte(sign))
	t.Log(str)
}

func TestHmacSha256(t *testing.T) {
	//  \ouput: c95b645cc57b18b3a080505dde495f97f5a2ae5755c663815be64b9d7727ff00 len: 64
	repo := inject.GetRegistryRepo()
	key, _ := repo.GetValue(registry.SysPrivateKey)
	domain.ConfigPrivateKey(key)
	str := domain.HmacSha256(crypto.Md5([]byte("123456")))
	t.Log(str, "len:", len(str))
}
