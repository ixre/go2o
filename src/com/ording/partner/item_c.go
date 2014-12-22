package partner

import (
	"com/ording/cache"
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

type itemC struct {
	Context app.Context
}

//食物列表
func (this *itemC) List(w http.ResponseWriter, r *http.Request) {
	/*
		'''
		菜单列表
		'''
		req=web.input(cid=-1,returnUri='')
		_dataurl='index?m=food&act=foods&ajax=1&cid=%s'%(req.cid)

		return render.foods(dataurl=_dataurl)
	*/
	r.ParseForm()
	//cid,_:= strconv.Atoi(r.FormValue("cid"))
	this.Context.Template().Render(w,
		"views/partner/item/list.html",
		nil)
}

func (this *itemC) Create(w http.ResponseWriter, r *http.Request, ptId int) {
	shopChks := cache.GetShopCheckboxs(ptId, "")
	cateOpts := cache.GetDropOptionsOfCategory(ptId)

	this.Context.Template().Render(w,
		"views/partner/item/create_item.html",
		func(m *map[string]interface{}) {
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *itemC) Edit(w http.ResponseWriter, r *http.Request, ptId int) {
	var e entity.FoodItem
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e = dao.Item().GetFoodItemById(ptId, id)
	js, _ := json.Marshal(e)

	shopChks := cache.GetShopCheckboxs(ptId, e.ApplySubs)
	cateOpts := cache.GetDropOptionsOfCategory(ptId)

	this.Context.Template().Render(w,
		"views/partner/item/update_item.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *itemC) SaveItem_post(w http.ResponseWriter, r *http.Request, ptId int) {
	var result cf.JsonResult
	r.ParseForm()

	e := entity.FoodItem{}
	web.ParseFormToEntity(r.Form, &e)

	t := time.Now()
	e.UpdateTime = t

	//更新
	if e.Id > 0 {
		origin := dao.Item().GetFoodItemById(ptId, e.Id)
		e.CreateTime = origin.CreateTime
	} else {
		e.CreateTime = t
	}

	id, err := dao.Item().SaveFoodItem(&e)

	if err != nil {
		result = cf.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *itemC) Del_post(w http.ResponseWriter, r *http.Request, ptId int) {
	var result cf.JsonResult

	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	_, err := dao.Item().DelFoodItem(ptId, id)

	if err != nil {
		result = cf.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}
