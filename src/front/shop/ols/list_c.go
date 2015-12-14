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
	"github.com/jsix/gof/algorithm/iterator"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/pager"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/domain/util"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/front"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ListC struct {
}

func categoryWalk(buf *bytes.Buffer, cs []*sale.ValueCategory) {
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.ValueCategory)
		if c.Id == 0 {
			return
		}
		if level == 1 {
			buf.WriteString(fmt.Sprintf("<div class=\"cat_panel\"><div class=\"t1 t1_%d\"><span class=\"icon\"></span><a href=\"%s\"><strong>%s</strong></a></div>", c.Id, c.Url, c.Name))
		} else if level == 2 {
			buf.WriteString(fmt.Sprintf("<ul><li><a href=\"%s\"><b>%s</b></a></li>", c.Url, c.Name))
		} else if level == 3 || level > 3 {
			buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", c.Url, c.Name))
		}
	}

	var over iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.ValueCategory)
		if c.Id != 0 {
			if level == 1 {
				buf.WriteString("</div>")
			} else if level == 2 {
				buf.WriteString("</ul>")
			}
		}
	}
	util.WalkCategory(cs, &sale.ValueCategory{Id: 0}, start, over)
}

// 类目，限移动端
func (this *ListC) All_cate(ctx *echox.Context) error {
	p := getPartner(ctx)
	mm := getMember(ctx)
	siteConf := getSiteConf(ctx)

	categories := dps.SaleService.GetCategories(p.Id)
	buf := bytes.NewBufferString("")
	categoryWalk(buf, categories)
	web.SetCacheHeader(ctx.Response(), 10)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":   p,
		"member":    mm,
		"conf":      siteConf,
		"cate_html": template.HTML(buf.String()),
	}
	return ctx.RenderOK("category.html", d)

}

func (this *ListC) getIdArray(path string) []int {
	idStr := path[strings.Index(path, "-")+1 : strings.LastIndex(path, ".")]
	strArr := strings.Split(idStr, "-")
	intArr := make([]int, len(strArr))
	for i, v := range strArr {
		intArr[i], _ = strconv.Atoi(v)
	}
	return intArr
}

func (this *ListC) GetSorter(ctx *echox.Context) error {
	return nil
}

// 商品列表
func (this *ListC) List_Index(ctx *echox.Context) error {
	r := ctx.Request()
	p := getPartner(ctx)
	const size int = 20 //-1表示全部
	sortQuery := ctx.Query("sort")
	idArr := this.getIdArray(r.URL.Path)
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	categoryId := idArr[len(idArr)-1]
	cat := dps.SaleService.GetCategory(p.Id, categoryId)

	total, items := dps.SaleService.GetPagedOnShelvesGoods(p.Id, categoryId,
		(page-1)*size, page*size, sortQuery)

	var pagerHtml string
	if total > size {
		pager := pager.NewUrlPager(pager.TotalPage(total, size), page, pager.GetterDefaultPager)
		pager.RecordCount = total
		pagerHtml = pager.PagerString()
	}

	buf := bytes.NewBufferString("")

	if len(items) == 0 {
		buf.WriteString("<div class=\"no_goods\">没有找到商品!</div>")
	} else {
		for i, v := range items {
			var hasDisCls string = ""
			if v.SalePrice == v.Price {
				hasDisCls = "no-disc"
			}
			buf.WriteString(fmt.Sprintf(`
				<div class="item item-i%d">
					<div class="block">
						<a href="/goods-%d.htm" class="goods-link">
							<img class="goods-img" src="%s" alt="%s"/>
							<h3 class="name">%s</h3>
							<span class="sale-price">￥%s</span>
							<span class="market-price %s"><del>￥%s</del></span>
						</a>
					</div>
                    <div class="clear-fix"></div>
                </div>
		`, i%2, v.GoodsId, format.GetGoodsImageUrl(v.Image),
				v.Name, v.Name, format.FormatFloat(v.SalePrice),
				hasDisCls, format.FormatFloat(v.Price)))
		}
	}

	sortBar := front.GetSorterHtml(front.GoodsListSortItems,
		sortQuery,
		r.URL.RequestURI())

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"cat":      cat,
		"sort_bar": template.HTML(sortBar),
		"items":    template.HTML(buf.Bytes()),
		"pager":    template.HTML(pagerHtml),
	}
	return ctx.RenderOK("list.html", d)
}

// 销售标签列表
func (this *ListC) SaleTagGoodsList(ctx *echox.Context) error {

	r := ctx.Request()
	p := getPartner(ctx)

	const size int = 20
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	i := strings.LastIndex(r.URL.Path, "/")
	tagCode := r.URL.Path[i+1:]

	saleTag := dps.SaleService.GetSaleTagByCode(p.Id, tagCode)
	if saleTag == nil {
		http.Error(ctx.Response(), "404 file not found!", http.StatusNotFound)
		return nil
	}

	total, items := dps.SaleService.GetPagedValueGoodsBySaleTag(p.Id, saleTag.Id, (page-1)*size, page*size)
	var pagerHtml string
	if total > size {
		pager := pager.NewUrlPager(pager.TotalPage(total, size), page, pager.GetterDefaultPager)
		pager.RecordCount = total
		pagerHtml = pager.PagerString()
	}

	buf := bytes.NewBufferString("")

	if len(items) == 0 {
		buf.WriteString("<div class=\"no_goods\">没有找到商品!</div>")
	} else {
		for i, v := range items {
			buf.WriteString(fmt.Sprintf(`
				<div class="item item-i%d">
					<div class="block">
						<a href="/goods-%d.htm" class="goods-link">
							<img class="goods-img" src="%s" alt="%s"/>
							<h3 class="name">%s</h3>
							<span class="sale-price">￥%s</span>
							<span class="market-price"><del>￥%s</del></span>
						</a>
					</div>
                    <div class="clear-fix"></div>
                </div>
		`, i%2, v.GoodsId, format.GetGoodsImageUrl(v.Image),
				v.Name, v.Name, format.FormatFloat(v.SalePrice),
				format.FormatFloat(v.Price)))
		}
	}
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"sale_tag": saleTag,
		"items":    template.HTML(buf.Bytes()),
		"pager":    template.HTML(pagerHtml),
	}
	return ctx.RenderOK("sale_tag.html", d)
}

// 商品详情
func (this *ListC) GoodsView(ctx *echox.Context) error {
	r := ctx.Request()
	p := getPartner(ctx)
	path := r.URL.Path
	goodsId, _ := strconv.Atoi(path[strings.LastIndex(path, "-")+1 : strings.Index(path, ".")])

	m := getMember(ctx)
	var level int = 0
	if m != nil {
		level = m.Level
	}
	goods, proMap := dps.SaleService.GetGoodsDetails(p.Id, goodsId, level)
	goods.Image = format.GetGoodsImageUrl(goods.Image)

	// 促销价 & 销售价
	var promPrice string
	var salePrice string

	if goods.PromPrice < goods.SalePrice {
		promPrice = fmt.Sprintf(`<span class="prom-price">￥<b>%s</b></span>`,
			format.FormatFloat(goods.PromPrice))
		salePrice = fmt.Sprintf("<del>￥%s</del>", format.FormatFloat(goods.SalePrice))
	} else {
		salePrice = "￥" + format.FormatFloat(goods.SalePrice)
	}

	// 促销信息
	var promDescribe string
	var promCls string = "hidden"
	if len(proMap) != 0 {
		promCls = ""
		buf := bytes.NewBufferString("")
		var i int = 0
		for k, v := range proMap {
			i++
			buf.WriteString(fmt.Sprintf(`<div class="prom-box prom%d">
					<span class="bg_txt red">%s</span>：<span class="describe">%s</span></div>`, i, k, v))
		}
		promDescribe = buf.String()
	}

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"goods":         goods,
		"promap":        proMap,
		"prom_price":    template.HTML(promPrice),
		"sale_price":    template.HTML(salePrice),
		"prom_describe": template.HTML(promDescribe),
		"prom_cls":      promCls,
	}
	return ctx.RenderOK("goods.html", d)

}

func (this *ListC) GoodsDetails(ctx *echox.Context) error {
	goodsId, _ := strconv.Atoi(ctx.Query("id"))
	describe := dps.SaleService.GetItemDescriptionByGoodsId(GetSessionPartnerId(ctx), goodsId)

	if len(describe) == 0 {
		describe = "<div style=\"text-align:center;color:#F00\">商品暂无描述</div>"
	}
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"describe": template.HTML(describe),
	}
	return ctx.RenderOK("goods-describe.html", d)
}
