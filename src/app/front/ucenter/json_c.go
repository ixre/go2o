/**
 * Copyright 2015 @ z3q.net.
 * name : json_c.go
 * author : jarryliu
 * date : 2016-04-25 23:09
 * description :
 * history :
 */
package ucenter

import (
	"encoding/gob"
	"github.com/jsix/gof/crypto"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
)

const (
	//todo: ??? 设置为可配置
	maxSeconds int64 = 10
)

func init() {
	gob.Register(map[string]map[string]interface{}{})
	gob.Register(ad.ValueGallery{})
	gob.Register(ad.ValueAdvertisement{})
	gob.Register([]*valueobject.Goods{})
	gob.Register(valueobject.Goods{})
}

type jsonC struct {
}

func getMd5(s string) string {
	return crypto.Md5([]byte(s))[8:16]
}

// 广告
func (t *jsonC) Member(ctx *echox.Context) error {
	memberId := GetSessionMemberId(ctx)
	ms := dps.MemberService.GetMemberSummary(memberId)
	return ctx.JSON(http.StatusOK, ms)
}
