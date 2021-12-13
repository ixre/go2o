package tests

import (
	"go2o/tests/ti"
	"testing"
)

func TestGetGroups(t *testing.T) {
	manager := ti.Factory.GetAdRepo().GetAdManager()
	arr := manager.GetGroups()
	t.Log(arr)
}

func TestQueryAd(t *testing.T) {
	manager := ti.Factory.GetAdRepo().GetAdManager()
	ad := manager.QueryAd("", 10)
	t.Log(len(ad))
}
