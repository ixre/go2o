/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type partnerC struct {
	*baseC
}

func (c *partnerC) Index(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/master/partner_partner_index.html")
}

func (c *partnerC) CreatePartner(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/master/partner/partner_create.html")
}

func (c *partnerC) CreatePartner_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.Message
	r.ParseForm()

	partner := partner.ValuePartner{}
	web.ParseFormToEntity(r.Form, &partner)

	dt := time.Now()
	anousPwd := strings.Repeat("*", 10) //匿名密码
	if len(partner.Pwd) != 0 && partner.Pwd != anousPwd {
		partner.Pwd = domain.EncodePartnerPwd(partner.Usr, partner.Pwd)
	}

	//更新
	if partner.Id > 0 {
		original, _ := dps.PartnerService.GetPartner(partner.Id)
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

	id, err := dps.PartnerService.SavePartner(partner.Id, &partner)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Data = id
		result.Result = true
	}
	w.Write(result.Marshal())
}

// 商户配置管理
func (this *partnerC) PartnerConf(ctx *web.Context) {
	var partnerId int
	partnerId, _ = strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"partnerId": partnerId,
	}, "views/master/partner/partner_create.html")
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
	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"entity": entityJson,
	}, "views/master/partner/partner_edit.html")
}
func (c *partnerC) List(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil,
		"views/master/partner/partner_list.html")
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
