package testing

import (
	"go2o/core/testing/ti"
	"testing"
)

func TestCheckExists(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	b := repo.CheckUsrExist("jarry6", 01)
	t.Log("是否已经使用:", b)
}
