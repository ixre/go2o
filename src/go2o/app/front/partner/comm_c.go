/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/app"
	"go2o/app/front"
	"net/http"
)

type commC struct {
	*front.WebCgi
	app.Context
}

func (this *mainC) GeoLocation(w http.ResponseWriter, r *http.Request) {
	this.WebCgi.GeoLocation(w, r)
}

//地区Json
//func (this *mainC) ChinaJson(w http.ResponseWriter, r *http.Request) {
//	var node *tree.TreeNode = dao.Common().GetChinaTree()
//	json, _ := json.Marshal(node)
//	w.Write(json)
//}
