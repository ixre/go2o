package partner

import (
	"encoding/json"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web/ui/tree"
	"go2o/app/front"
	"go2o/core/ording/dao"
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
func (this *mainC) ChinaJson(w http.ResponseWriter, r *http.Request) {
	var node *tree.TreeNode = dao.Common().GetChinaTree()
	json, _ := json.Marshal(node)
	w.Write(json)
}
