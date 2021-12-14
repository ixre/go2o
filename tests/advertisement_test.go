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

func TestQueryAdvertisementData(t *testing.T){
	manager := ti.Factory.GetAdRepo().GetAdManager()
	iu := manager.GetUserAd(0)
	keys := []string{"MOBI-SHOP-INDEX-SCROLLER","mobi-index-hot-win-banner","mobile-index-a2"}
	advertisement := iu.QueryAdvertisement(keys)
	if l := len(advertisement) ;l== 0 || l ==3{
		t.Fail()
	}
}
