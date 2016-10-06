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
			<div class="goods-item">
				   <a target="_blank" href="/goods-%d.htm">
				     <img src="%s" class="goods-thumb" />
				   </a>
				   	<span class="goods-title">%s
							<span class="goods-num">x%d</span>
					</span>
				   	<span class="goods-no">商品编号：<i>%s</i></span>
				<span class="goods-price">￥%s</span>
				<span class="goods-fee">￥%s</span>
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
