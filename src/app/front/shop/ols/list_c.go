/**
 * Copyright 2015 @ S1N1 Team.
 * name : list_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
)

type listC struct {
	*baseC
}

func (this *listC) Index(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	p := this.GetPartner(ctx)
	mm := this.GetMember(ctx)
	pa := this.GetPartnerApi(ctx)

	siteConf := this.GetSiteConf(ctx)
	categories := GetCategories(ctx.App, p.Id, pa.ApiSecret)
	ctx.App.Template().Execute(w, gof.TemplateDataMap{
		"partner":    p,
		"title":      "在线订餐-" + p.Name,
		"categories": template.HTML(categories),
		"member":     mm,
		"conf":       siteConf,
	},
		"views/shop/ols/list.html",
		"views/shop/ols/inc/header.html",
		"views/shop/ols/inc/footer.html")
}

func (this *listC) GetList(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	p := this.GetPartner(ctx)

	const getNum int = -1 //-1表示全部
	categoryId, _ := strconv.Atoi(r.URL.Query().Get("cid"))
	items := dps.SaleService.GetOnShelvesGoodsByCategoryId(p.Id, categoryId, getNum)
	if len(items) == 0 {
		w.Write([]byte(`无商品`))
		return
	}
	buf := bytes.NewBufferString("<ul>")

	for _, v := range items {

		buf.WriteString(fmt.Sprintf(`
			<li>
				<div class="gs_goodss">
                        <img src="%s" alt="%s"/>
                        <h3 class="name">%s%s</h3>
                        <span class="srice">原价:￥%s</span>
                        <span class="sprice">优惠价:￥%s</span>
                        <a href="javascript:cart.add(%d,1);" class="add">&nbsp;</a>
                </div>
             </li>
		`, format.GetGoodsImageUrl(v.Image), v.Name, v.Name, v.SmallTitle, format.FormatFloat(v.Price),
			format.FormatFloat(v.SalePrice),
			v.Id))
	}
	buf.WriteString("</ul>")
	w.Write(buf.Bytes())
}
