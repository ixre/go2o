package tests

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

func TestCheckExists(t *testing.T) {
	repo := inject.GetMemberRepo()
	b := repo.CheckUserExist("jarry6", 01)
	t.Log("是否已经使用:", b)
}
