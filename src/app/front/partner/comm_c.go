/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/front"
)

var _ mvc.Filter = new(commC)

type commC struct {
	Base *baseC
	*front.WebCgi
}

func (this *commC) Requesting(ctx *web.Context) bool {
	return this.Base.Requesting(ctx)
}
func (this *commC) RequestEnd(ctx *web.Context) {
	this.Base.RequestEnd(ctx)
}

func (this *mainC) GeoLocation(ctx *web.Context) {
	this.WebCgi.GeoLocation(ctx)
}

//地区Json
//func (this *mainC) ChinaJson(ctx *web.Context) {
//	var node *tree.TreeNode = dao.Common().GetChinaTree()
//	json, _ := json.Marshal(node)
//	w.Write(json)
//}
