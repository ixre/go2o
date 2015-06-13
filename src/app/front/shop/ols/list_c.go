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
	"strings"
	"github.com/atnet/gof/web/pager"
)

type ListC struct {
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
func (this *ListC) AllCate(ctx *web.Context) {
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

func (this *ListC) getIdArray(path string)[]int{
	idStr := path[strings.Index(path,"-")+1:strings.LastIndex(path,".")]
	strArr := strings.Split(idStr,"-")
	intArr := make([]int,len(strArr))
	for i,v := range strArr{
		intArr[i],_ = strconv.Atoi(v)
	}
	return intArr
}

func (this *ListC) List_Index(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {
		r := ctx.Request
		p := this.BaseC.GetPartner(ctx)

		const size int = 20 //-1表示全部

		idArr := this.getIdArray(r.URL.Path)
		page, _ := strconv.Atoi(r.FormValue("page"))
		if page < 1 {
			page = 1
		}
		categoryId := idArr[len(idArr)-1]
		cat := dps.SaleService.GetCategory(p.Id, categoryId)

		total, items := dps.SaleService.GetPagedOnShelvesGoods(p.Id, categoryId,(page-1)*size,page*size)
		var pagerHtml string
		if total > size {
			pager := pager.NewUrlPager(pager.TotalPage(total, size), page, pager.GetterDefaultPager)
			pager.RecordCount = total
			pagerHtml = pager.PagerString()
		}

		buf := bytes.NewBufferString("")

		if len(items) == 0 {
			buf.WriteString("<div class=\"no_goods\">没有找到商品!</div>")
		}else {
			for i, v := range items {
				buf.WriteString(fmt.Sprintf(`
				<div class="item-block item-block%d">
					<a href="/item-%d.htm" class="goods-link">
                        <img class="goods-img" src="%s" alt="%s"/>
                        <h3 class="name">%s%s</h3>
                        <span class="sale-price">￥%s</span><br />
                        <span class="market-price"><del>￥%s</del></span>
					</a>
                    <div class="clearfix"></div>
                </div>
		`, i%2, v.Id, format.GetGoodsImageUrl(v.Image),
				v.Name, v.Name, v.SmallTitle, format.FormatFloat(v.SalePrice),
				format.FormatFloat(v.Price)))
			}
		}

		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"cat":cat,
			"items":template.HTML(buf.Bytes()),
			"pager":template.HTML(pagerHtml),
		},
		"views/shop/{device}/list.html",
		"views/shop/{device}/inc/header.html",
		"views/shop/{device}/inc/footer.html")
	}
}

func (this *ListC) GoodsDetails(ctx *web.Context){
	if this.BaseC.Requesting(ctx) {
		r := ctx.Request
		p := this.BaseC.GetPartner(ctx)
		path := r.URL.Path
		goodsId,_ := strconv.Atoi(path[strings.Index(path,"-")+1:strings.Index(path,".")])

		goods := dps.SaleService.GetValueGoods(p.Id,goodsId)

		goods.Image = format.GetGoodsImageUrl(goods.Image)
		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"goods":goods,
		},
		"views/shop/{device}/goods.html",
		"views/shop/{device}/inc/header.html",
		"views/shop/{device}/inc/footer.html")
	}
}
