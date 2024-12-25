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
	domain.ConfigureHmacPrivateKey(key)
}

// 管理员密码
func TestMasterPwd(t *testing.T) {
	user := "master"
	pwd := "123456"
	sha1 := domain.HmacSha256(domain.Md5(pwd) + user)
	t.Log(sha1)
}

func TestMemberPwd(t *testing.T) {

	// 96a28f426f26afb6ba747626bc518f10e7a6d80ddd133232ff228aab254e3f03
	// 0a45e3f114becc7894987b3e2507a5052294796280f043de1d19b8f706145b99
	// real:
	// 1d6dffa5443e6e2dc859c3beb12f302d18537324156d85d31c2c259e2949eae7
	pwd := domain.Md5("123456")
	salt := "eAhLIi"
	t.Log("--pwd=", pwd, "\n")
	pwd = domain.MemberSha256Pwd(pwd, salt)
	t.Log("--pwd=", pwd, "\n")

	pwd = domain.MerchantSha265Pwd(pwd, salt)
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
	domain.ConfigureHmacPrivateKey(key)
	str := domain.HmacSha256(crypto.Md5([]byte("123456")))
	t.Log(str, "len:", len(str))
}
