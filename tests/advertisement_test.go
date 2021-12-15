package tests

import (
	"fmt"
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

func TestQueryAdvertisementData(t *testing.T) {
	manager := ti.Factory.GetAdRepo().GetAdManager()
	iu := manager.GetUserAd(0)
	keys := []string{"mobile-index-swiper-1",
		"mobile-index-image-1",
		"mobile-index-image-2"}
	advertisement := iu.QueryAdvertisement(keys)
	for _, v := range advertisement {
		t.Log(fmt.Sprintf("%#v", v.Dto()))
	}
	if l := len(advertisement); l == 0 || l == 3 {
		t.Fail()
	}
}
