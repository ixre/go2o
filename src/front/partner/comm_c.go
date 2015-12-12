/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"go2o/src/front"
	"go2o/src/x/echox"
)


type commC struct {
	*front.WebCgi
}

func (this *mainC) GeoLocation(ctx *echox.Context) error {
	//this.WebCgi.GeoLocation(ctx)
	//todo:???
	return nil
}

//地区Json
//func (this *mainC) ChinaJson(ctx *echox.Context)error{
//	var node *tree.TreeNode = dao.Common().GetChinaTree()
//	json, _ := json.Marshal(node)
//	w.Write(json)
//}
