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
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"github.com/atnet/gof/algorithm/iterator"
	"go2o/src/core/infrastructure/domain/util"
)

type CategoryC struct {
	*BaseC
}

func categoryWalk(buf *bytes.Buffer,cs []*sale.ValueCategory) {
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.ValueCategory)
		if c.Id == 0 {
			return
		}
		if level == 1 {
			buf.WriteString(fmt.Sprintf("<div class=\"cat_panel\"><div class=\"t1 t1_%d\"><a href=\"%s\"><strong>%s</strong></a></div>", c.Id, c.Url, c.Name))
		}else if level == 2 {
			buf.WriteString(fmt.Sprintf("<ul><li><a href=\"%s\"><b>%s</b></a></li>", c.Url, c.Name))
		}else if level == 3 || level>3 {
			buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", c.Url, c.Name))
		}
	}

	var over iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.ValueCategory)
		if c.Id != 0 {
			if level == 1 {
				buf.WriteString("</div>")
			}else if level == 2 {
				buf.WriteString("</ul>")
			}
		}
	}
	util.WalkCategory(cs, &sale.ValueCategory{Id: 0}, start, over)
}

// 类目，限移动端
func (this *CategoryC) Index(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	mm := this.BaseC.GetMember(ctx)
	siteConf := this.BaseC.GetSiteConf(ctx)

	categories := dps.SaleService.GetCategories(p.Id)
	buf := bytes.NewBufferString("")
	categoryWalk(buf,categories)

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":  p,
		"member":   mm,
		"conf":     siteConf,
		"cate_html": template.HTML(buf.String()),
	},
		"views/shop/{device}/category.html",
		"views/shop/{device}/inc/header.html",
		"views/shop/{device}/inc/footer.html")
}

func (this *CategoryC) GetList(ctx *web.Context) {
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
