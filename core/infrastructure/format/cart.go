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
	"go2o/core/service/thrift/idl/gen-go/define"
)

// 购物车详情
func CartDetails(c *define.ShoppingCart) string {
	buf := bytes.NewBufferString("")
	for _, shop := range c.Shops {
		if shop.Checked {
			buf.WriteString(fmt.Sprintf(`<div class="vendor"><div class="tit">%s</div>`,
				shop.ShopName))
			for _, it := range shop.Items {
				if !it.Checked {
					continue //只显示结账的
				}
				buf.WriteString(fmt.Sprintf(`
			<div class="product clearfix">
				   <a target="_blank" href="/item-%d.html">
				     <img src="%s" class="item-image" />
				   </a>
				   	<span class="title">%s
							<span class="quantity">x%d</span>
					</span>
				   	<span class="code">商品编号：<i>%s</i></span>
				<span class="price">￥%s</span>
				<span class="fee">￥%s</span>
			</div>`,
					it.ItemId, GetGoodsImageUrl(it.Image), it.Title, it.Quantity, it.Code,
					FormatFloat64(it.Price), FormatFloat64(it.Price*float64(it.Quantity)),
				))
			}
			buf.WriteString("</div>")
		}
	}

	return buf.String()
}
