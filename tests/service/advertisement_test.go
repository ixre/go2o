package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/typeconv"
)

func TestQueryHyperLinkAdvertisementData(t *testing.T) {
	ret, _ := inject.GetAdvertisementService().GetAdvertisement(context.TODO(), &proto.AdIdRequest{
		AdUserId:   0,
		AdId:       8,
		AdKey:      "",
		ReturnData: true,
	})
	t.Log("广告数据:", typeconv.MustJson(ret))
}

// 测试获取广告数据
func TestQueryAdvertisementData(t *testing.T) {
	keys := []string{
		"mobile-index-swiper-1",
		"mobile-index-image-1",
		"mobile-index-image-2",
	}
	ret, _ := inject.GetAdvertisementService().QueryAdvertisementData(context.TODO(), &proto.QueryAdvertisementDataRequest{
		AdUserId: 0,
		Keys:     keys,
	})
	t.Log("广告数据:", typeconv.MustJson(ret))
}
