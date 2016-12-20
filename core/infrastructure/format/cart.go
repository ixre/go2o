/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-11 20:57
 * description :
 * history :
 */
package format

import (
	"bytes"
	"fmt"
	"go2o/core/dto"
)

// 购物车详情
func CartDetails(c *dto.ShoppingCart) string {
	buf := bytes.NewBufferString("")
	for _, vendor := range c.Vendors {
		if vendor.CheckedNum > 0 {
			buf.WriteString(fmt.Sprintf(`<div class="vendor"><div class="tit">%s</div>`,
				vendor.ShopName))
			for _, item := range vendor.Items {
				if !item.Checked {
					continue //只显示结账的
				}
				buf.WriteString(fmt.Sprintf(`
			<div class="product clearfix">
				   <a target="_blank" href="/item-%d.htm">
				     <img src="%s" class="item-image" />
				   </a>
				   	<span class="title">%s
							<span class="quantity">x%d</span>
					</span>
				   	<span class="code">商品编号：<i>%s</i></span>
				<span class="price">￥%s</span>
				<span class="fee">￥%s</span>
			</div>`,
					item.GoodsId, GetGoodsImageUrl(item.GoodsImage), item.GoodsName, item.Quantity, item.GoodsNo,
					FormatFloat(item.SalePrice), FormatFloat(item.SalePrice*float32(item.Quantity)),
				))
			}
			buf.WriteString("</div>")
		}
	}

	return buf.String()
}
