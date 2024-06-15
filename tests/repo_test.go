package tests

import (
	"testing"

	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/inject"
)

func TestCheckExists(t *testing.T) {
	repo := inject.GetMemberRepo()
	b := repo.CheckUserExist("jarry6", 01)
	t.Log("是否已经使用:", b)
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
