/**
 * Copyright 2014 @ Ops.
 * name :
 * author : newmin
 * date : 2013-11-05 17:08
 * description :
 * history :
 */
package www

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/app"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/infrastructure/format"
	"go2o/core/ording/cache/apicache"
	"go2o/core/service/goclient"
	"go2o/core/share/variable"
	"html/template"
	"net/http"
	"strconv"
)

type ordingC struct {
	app.Context
}

func (this *ordingC) Index(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner, mm *member.ValueMember) {
	if b, siteConf := GetSiteConf(w, p); b {
		categories := apicache.GetCategories(this.Context, p.Id, p.Secret)
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			//(*m)["title"] = "在线订餐-" + p.Name
			(*m)["categories"] = template.HTML(categories)
			(*m)["member"] = mm
			(*m)["conf"] = siteConf
		},
			"views/web/www/ording.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *ordingC) List(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	const getNum int = -1 //-1表示全部
	categoryId, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		w.Write([]byte(`{"error":"yes"}`))
		return
	}
	items, err := goclient.Partner.GetItems(p.Id, p.Secret, categoryId, getNum)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	buf := bytes.NewBufferString("<ul>")

	noPicPath := this.Context.Config().GetString(variable.NoPicPath)

	for _, v := range items {

		if len(v.Image) == 0 {
			v.Image = noPicPath
		}

		buf.WriteString(fmt.Sprintf(`
			<li>
				<div class="gs_goodss">
                        <img src="%s/%s" alt="%s"/>
                        <h3 class="name">%s%s</h3>
                        <span class="srice">原价:￥%s</span>
                        <span class="sprice">优惠价:￥%s</span>
                        <a href="javascript:cart.add(%d,1);" class="add">&nbsp;</a>
                </div>
             </li>
		`, this.Context.Config().GetString(variable.ImageServer),
			v.Image, v.Name, v.Name, v.SmallTitle, format.FormatFloat(v.Price),
			format.FormatFloat(v.SalePrice),
			v.Id))
	}
	buf.WriteString("</ul>")
	w.Write(buf.Bytes())
}
