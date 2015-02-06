package partner

//关于regexp的资料
//http://www.cnblogs.com/golove/archive/2013/08/20/3270918.html

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/dps"
	"html/template"
	"net/http"
	"strconv"
)

//==========================  后台逻辑 =============================//
type shopC struct {
	Context app.Context
}

func (this *shopC) ShopList(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/partner/shop/shop_list.html", nil)
}

//修改门店信息
func (this *shopC) Create(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w,
		"views/partner/shop/create.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS("{}")
		})
}

//修改门店信息
func (this *shopC) Modify(w http.ResponseWriter, r *http.Request, partnerId int) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	shop := dps.PartnerService.GetShopValueById(partnerId, id)
	entity, _ := json.Marshal(shop)

	this.Context.Template().Render(w,
		"views/partner/shop/modify.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(entity)
		})
}

//修改门店信息
func (this *shopC) SaveShop_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result gof.JsonResult
	r.ParseForm()

	shop := partner.ValueShop{}
	web.ParseFormToEntity(r.Form, &shop)

	id, err := dps.PartnerService.SaveShop(partnerId, &shop)

	if err != nil {
		result = gof.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *shopC) Del_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	shopId, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.PartnerService.DeleteShop(partnerId, shopId)
	}

	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}
