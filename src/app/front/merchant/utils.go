/**
 * Copyright 2015 @ z3q.net.
 * name : utils
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"bytes"
	"fmt"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
)

func getMerchantId(ctx *echox.Context) int {
	obj := ctx.Session.Get("merchant_id")
	if obj != nil {
		return obj.(int)
	}
	return 0
}

func getSaleTagsCheckBoxHtml(tags []*sale.ValueSaleTag) string {
	if len(tags) == 0 || tags == nil {
		return `<div style="color:red">没有找到任何销售标签!</div>`
	}
	buf := bytes.NewBufferString(`<ul class="sale_tags">`)
	for i, v := range tags {
		buf.WriteString(fmt.Sprintf(`<li><input type="checkbox" id="sale_tag%d" field="SaleTags[%d]" value="%d" name="SaleTags"/>
            <label for="sale_tag%d">%s</label></li>`, i, i, v.Id, i, v.TagName))
	}
	buf.WriteString("</ul>")
	return buf.String()
}

// 获取等级下拉选项列表
func getLevelDropDownList(merchantId int) string {
	buf := bytes.NewBufferString("")
	lvs := dps.MerchantService.GetMemberLevels(merchantId)
	for _, v := range lvs {
		if v.Enabled == 1 {
			buf.WriteString(fmt.Sprintf(`<option value="%d">%s</option>`, v.Value, v.Name))
		}
	}
	return buf.String()
}

// 获取邮件模板选项
func getMailTemplateOpts(merchantId int) string {
	buf := bytes.NewBufferString("")
	list := dps.MerchantService.GetMailTemplates(merchantId)
	for _, v := range list {
		if v.Enabled == 1 {
			buf.WriteString(fmt.Sprintf(`<option value="%d">%s</option>`, v.Id, v.Name))
		}
	}
	return buf.String()
}
