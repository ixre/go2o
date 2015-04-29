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
	"go2o/src/app/front"
)

type commC struct {
	*front.WebCgi
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
