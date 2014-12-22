package partner

//关于regexp的资料
//http://www.cnblogs.com/golove/archive/2013/08/20/3270918.html

import (
	"com/ording/dao"
	"com/ording/entity"
	"encoding/json"
	"html/template"
	"net/http"
	"ops/cf"
	"ops/cf/app"
	"ops/cf/web"
	"strconv"
	"time"
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
func (this *shopC) Modify(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	shop := dao.Shop().GetShopById(id)
	entity, _ := json.Marshal(shop)

	this.Context.Template().Render(w,
		"views/partner/shop/modify.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(entity)
		})
}

//修改门店信息
func (this *shopC) SaveShop_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result cf.JsonResult
	r.ParseForm()

	shop := entity.Shop{}
	web.ParseFormToEntity(r.Form, &shop)

	//更新
	if shop.Id > 0 {
		orgialShop := dao.Shop().GetShopById(shop.Id)
		shop.CreateTime = orgialShop.CreateTime
		shop.PartnerId = orgialShop.PartnerId
	} else {
		shop.CreateTime = time.Now()
		shop.PartnerId = partnerId
	}

	id, err := dao.Shop().SaveShop(&shop)
	if err != nil {
		result = cf.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *shopC) Del_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	shopId, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dao.Shop().DeleteShop(partnerId, shopId)
	}

	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}
