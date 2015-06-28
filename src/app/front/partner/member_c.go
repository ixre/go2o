/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
)

var _ mvc.Filter = new(memberC)

type memberC struct {
	*baseC
}

func (this *memberC) LevelList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/partner/member/level_list.html")
}

//修改门店信息
func (this *memberC) EditMLevel(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.PartnerService.GetMemberLevelById(partnerId, id)
	bys, _ := json.Marshal(entity)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/member/edit_level.html")
}

func (this *memberC) CreateMLevel(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
			"entity": "{}",
		},
		"views/partner/member/create_level.html")
}

func (this *memberC) SaveMLevel_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	var result gof.Message
	r.ParseForm()

	e := valueobject.MemberLevel{}
	web.ParseFormToEntity(r.Form, &e)
	e.PartnerId = this.GetPartnerId(ctx)

	id, err := dps.PartnerService.SaveMemberLevel(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	this.JsonOutput(ctx, result)
}

func (this *memberC) DelMLevel(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.PartnerService.DelMemberLevel(partnerId, id)
	}

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	this.ResultOutput(ctx, result)
}

// 会员列表
func (this *memberC) List(ctx *web.Context) {
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
		}, "views/partner/member/member_list.html")
}


func (this *memberC) Lock_member_post(ctx *web.Context){
	id,_ := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	partnerId := this.GetPartnerId(ctx)
	var result gof.Message
	if  _,err := dps.MemberService.LockMember(partnerId,id);err != nil{
		result.Message = err.Error()
	}else{
		result.Result = true
	}
	this.ResultOutput(ctx,result)
}

func (this *memberC) Cancel(ctx *web.Context) {
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/partner/order/cancel.html")

}

func (this *memberC) Cancel_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dps.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte(`{result:false,message:"` + err.Error() + `"}`))
	}
}

func (this *memberC) View(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
	}

	member := dps.MemberService.GetMember(e.MemberId)
	if member == nil {
		w.Write([]byte("无效订单"))
		return
	}

	e.ItemsInfo = strings.Replace(e.ItemsInfo, "\n", "<br />", -1)
	if len(e.Note) == 0 {
		e.Note = "无备注"
	}

	js, _ := json.Marshal(e)

	var shopName string
	var payment string
	var orderStateText string
	if e.ShopId == 0 {
		shopName = "未指定"
	} else {
		shopName = dps.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
	}
	payment = enum.GetPaymentName(e.PaymentOpt)
	orderStateText = enum.OrderState(e.Status).String()

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":   template.JS(js),
			"member":   member,
			"shopName": shopName,
			"payment":  payment,
			"state":    orderStateText,
		}, "views/partner/order/order_view.html")
}
