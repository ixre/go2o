package domain

import (
	"strconv"
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/util"
)

func TestGenerateAppId(t *testing.T) {
	for {
		s := strconv.Itoa(util.RandInt(8))
		t.Log(s)
		time.Sleep(1000)
	}
}

func TestUpdateRegistryValue(t *testing.T) {
	repo := ti.Factory.GetRegistryRepo()
	ir := repo.Get(registry.OrderAffiliatePushEnabled)
	if ir != nil {
		err := ir.Update("2")
		if err == nil {
			err = ir.Save()
		}
		if err != nil {
			t.Error(err)
		}
	}
}
