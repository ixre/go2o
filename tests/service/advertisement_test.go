package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/types/typeconv"
)

func TestQueryHyperLinkAdvertisementData(t *testing.T) {
	ret, _ := inject.GetAdService().GetAdvertisement(context.TODO(), &proto.AdIdRequest{
		AdUserId:   0,
		AdId:       8,
		AdKey:      "",
		ReturnData: true,
	})
	t.Log("广告数据:", typeconv.MustJson(ret))
}
