package domain

import (
	"fmt"
	"testing"

	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/types/typeconv"
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

func TestQueryHyperLinkAdvertisementData(t *testing.T) {
	manager := ti.Factory.GetAdRepo().GetAdManager()
	iu := manager.GetUserAd(0)
	ad := iu.GetById(8)
	v := ad.Dto()
	t.Log("广告数据:", typeconv.MustJson(v))
}
