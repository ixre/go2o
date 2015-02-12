/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package apicache

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/app"
	"go2o/core/domain/interface/enum"
	"go2o/core/service/goclient"
)

func GetShops(c app.Context, partnerId int, secret string) []byte {
	//分店
	var buf *bytes.Buffer = bytes.NewBufferString("")
	shops, err := goclient.Partner.GetShops(partnerId, secret)
	if shops == nil {
		if err != nil {
			c.Log().Panicf("[Error]:%s", err.Error())
		}
		return []byte("<div class=\"nodata noshop\">还未添加分店</div>")
	}
	buf.WriteString("<ul class=\"shops\">")
	for i, v := range shops {
		buf.WriteString(fmt.Sprintf(`<li class="s%d">
			<div class="name"><span><strong>%s</strong></div>
			<span class="shop-state shopstate%d">%s</span>
			<div class="phone">%s</div>
			<div class="address">%s</div>
			</li>`, i+1, v.Name, v.State, enum.GetFrontShopStateName(v.State), v.Phone, v.Address))
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}
