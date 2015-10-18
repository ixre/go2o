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
	"go2o/src/core/dto"
)

// 购物车详情
func CartDetails(c *dto.ShoppingCart) string {
	var byts *bytes.Buffer = bytes.NewBufferString("")
	byts.WriteString(`
		<table cellspacing="1" class="cart_details_table">
			<thead>
				<tr>
					<td><span class="t">商品</span></td>
					<td><span class="t">价格</span></td>
					<td><span class="t">数量</span></td>
					<td><span class="t">总价</span></td>
				</tr>
			</thead>
		`)

	for _, v := range c.Items {
		byts.WriteString(fmt.Sprintf(`
			<tr class="goods">
				<td class="goods-info">
				   <a target="_blank" href="/goods-%d.htm"><img src="%s" class="goods-thumb" />
				   	<span class="goods-title">%s</span></a>
				   	<span class="goods-no">商品编号：<i>%s</i></span>
				</td>

				<td>
					￥%s
				</td>
				<td>
					x%d
				</td>
				<td>
					￥%s
				</td>
			</tr>
		`,
			v.GoodsId, GetGoodsImageUrl(v.GoodsImage), v.GoodsName, v.GoodsNo,
			FormatFloat(v.SalePrice), v.Num, FormatFloat(v.SalePrice*float32(v.Num)),
		))
	}

	byts.WriteString("</table>")

	return byts.String()
}
