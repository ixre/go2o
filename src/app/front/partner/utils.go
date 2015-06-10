/**
 * Copyright 2015 @ S1N1 Team.
 * name : utils
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import (
	"bytes"
	"fmt"
	"go2o/src/core/domain/interface/sale"
)

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
