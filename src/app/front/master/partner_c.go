/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type partnerC struct {
}

func (c *partnerC) Index(ctx *echox.Context) error {
	return ctx.RenderOK("partner.index.html", ctx.NewData())
}

func (c *partnerC) CreateMerchant(ctx *echox.Context) error {
	return ctx.RenderOK("partner.create.html", ctx.NewData())
}

// 保存商户(POST)
func (c *partnerC) SaveMerchant(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Message
		var isCreate bool
		r.ParseForm()

		partner := merchant.MerchantValue{}
		web.ParseFormToEntity(r.Form, &partner)

		dt := time.Now()
		anousPwd := strings.Repeat("*", 10) //匿名密码
		if len(partner.Pwd) != 0 && partner.Pwd != anousPwd {
			partner.Pwd = domain.MerchantSha1Pwd(partner.Usr, partner.Pwd)
		}

		//更新
		if partner.Id > 0 {
			original, _ := dps.MerchantService.GetMerchant(partner.Id)
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
			isCreate = true
		}

		id, err := dps.MerchantService.SaveMerchant(partner.Id, &partner)
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Data = id
			result.Result = true
			if isCreate {
				// 初始化商户信息
			}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 商户配置管理
func (this *partnerC) MerchantConf(ctx *echox.Context) error {
	var merchantId int
	merchantId, _ = strconv.Atoi(ctx.Query("id"))
	d := ctx.NewData()
	d.Map["merchantId"] = merchantId
	return ctx.RenderOK("partner.create.html", d)
}

func (c *partnerC) EditMerchant(ctx *echox.Context) error {
	var entityJson template.JS
	id, err := strconv.Atoi(ctx.Query("id"))
	if err == nil {
		partner, err := dps.MerchantService.GetMerchant(id)
		if err == nil && partner != nil {
			partner.Pwd = strings.Repeat("*", 10)
			entity, _ := json.Marshal(partner)
			entityJson = template.JS(entity)
		}
	}
	d := ctx.NewData()
	d.Map["entity"] = entityJson
	return ctx.RenderOK("partner.edit.html", d)
}
func (c *partnerC) List(ctx *echox.Context) error {
	return ctx.RenderOK("partner.list.html", ctx.NewData())
}

func (c *partnerC) DelMerchant(w http.ResponseWriter, r *http.Request) {
	//	var result gof.Message
	//	r.ParseForm()
	//	ptid, err := strconv.Atoi(r.Form.Get("id"))
	//	if err != nil {
	//		result.Message = err.Error()
	//	} else {
	////		err := dps.MerchantService.DeleteMerchant(ptid)
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
