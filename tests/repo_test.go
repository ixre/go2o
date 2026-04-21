package tests

import (
	"testing"

	"github.com/ixre/go2o/pkg/initial/provide"
	"github.com/ixre/go2o/pkg/inject"
)

func TestCheckExists(t *testing.T) {
	repo := inject.GetMemberRepo()
	b := repo.CheckUserExist("jarry6", 01)
	if b {
		t.Fatal("用户已经存在")
	}
}

// 测试清除缓存
func TestCleanCache(t *testing.T) {
	sto := provide.GetStorageInstance()
	i, err := sto.DeleteWith("")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("清除缓存%d条", i)
}
