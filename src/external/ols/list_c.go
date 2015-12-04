/**
 * Copyright 2015 @ z3q.net.
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
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/cache/apicache"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
)

type listC struct {
	*baseC
}

func (this *listC) Index(ctx *web.Context) {
	_, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	mm := this.GetMember(ctx)
	pa := this.GetPartnerApi(ctx)

	if b, siteConf := GetSiteConf(w, p, pa); b {
		categories := apicache.GetCategories(ctx.App, p.Id, pa.ApiSecret)
		ctx.App.Template().Exellcute(w, gof.TemplateDataMap{
			"partner":    p,
			"categories": template.HTML(categories),
			"member":     mm,
			"conf":       siteConf,
		},
			"views/shop/ols/{device}/list.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	}
}

func (this *listC) GetList(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	pa := this.GetPartnerApi(ctx)

	const getNum int = -1 //-1表示全部
	categoryId, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		w.Write([]byte(`{"error":"yes"}`))
		return
	}
	items, err := goclient.Partner.GetItems(p.Id, pa.ApiSecret, categoryId, getNum)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
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
