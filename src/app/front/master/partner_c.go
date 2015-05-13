/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
 package master

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"github.com/atnet/gof/web"
	"go2o/src/core/service/dps"
	"github.com/atnet/gof"
	"go2o/src/core/domain/interface/partner"
	"time"
	"go2o/src/core/infrastructure/domain"
)

type partnerC struct {
	*baseC
}

func (c *partnerC) Index(ctx *web.Context) {
	ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter,nil,"views/master/partner_index.html")
}

func (c *partnerC) AddPartner(ctx *web.Context) {
	ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter,nil, "views/master/partner_add.html")
}

func (c *partnerC) EditPartner(ctx *web.Context) {
	var entityJson template.JS
	id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	if err == nil {
		partner, err := dps.PartnerService.GetPartner(id)
		if err == nil && partner != nil {
			partner.Pwd = strings.Repeat("*", 10)
			entity, _ := json.Marshal(partner)
			entityJson = template.JS(entity)
		}
	}
	ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter, func(mp *map[string]interface{}) {
		(*mp)["entity"] = entityJson
	}, "views/admin/partner_edit.html")
}

func (c *partnerC) SavePartner_post(w http.ResponseWriter, r *http.Request) {
	var result gof.Message
	r.ParseForm()

	partner := partner.ValuePartner{}
	web.ParseFormToEntity(r.Form, &partner)

	dt := time.Now()
	anousPwd := strings.Repeat("*", 10) //匿名密码
	if partner.Pwd != anousPwd {
		partner.Pwd = domain.EncodePartnerPwd(partner.Usr, partner.Pwd)
	}

	//更新
	if partner.Id > 0 {
		original,_ := dps.PartnerService.GetPartner(partner.Id)
		partner.JoinTime = original.JoinTime
		partner.ExpiresTime = original.ExpiresTime
		partner.UpdateTime = dt.Unix()

		if partner.Pwd == anousPwd {
			partner.Pwd = original.Pwd
		}
	} else {
		partner.JoinTime = dt.Unix()
		partner.ExpiresTime = dt.AddDate(10, 0, 0).Unix()
		partner.UpdateTime = dt.Unix()
	}

	id, err := dps.PartnerService.SavePartner(partner.Id,&partner)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Data = id
		result.Result = true
	}
	w.Write(result.Marshal())
}

func (c *partnerC) PartnerList(ctx *web.Context) {
	ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter,nil,
	"views/admin/partner_list.html")
}

func (c *partnerC) DelPartner_post(w http.ResponseWriter, r *http.Request) {
//	var result gof.Message
//	r.ParseForm()
//	ptid, err := strconv.Atoi(r.Form.Get("id"))
//	if err != nil {
//		result.Message = err.Error()
//	} else {
////		err := dps.PartnerService.DeletePartner(ptid)
////		if err != nil {
////			result.Message = err.Error()
////		} else {
////			result.Result = true
////		}
//	}
//	w.Write(result.Marshal())
}

//地区Json
func (this *partnerC) ChinaJson(w http.ResponseWriter, r *http.Request) {
//	var node *tree.TreeNode = logic.GetChinaTree()
//	json, _ := json.Marshal(node)
//	w.Write(json)
}
