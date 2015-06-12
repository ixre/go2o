/**
 * Copyright 2014 @ S1N1 Team.
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
					<td>商品</td>
					<td>价格</td>
					<td>数量</td>
					<td>总价</td>
				</tr>
			</thead>
		`)

	for _, v := range c.Items {
		byts.WriteString(fmt.Sprintf(`
			<tr>
				<td>
				   <a target="_blank" href="/item-%d.htm">
				   	<img src="%s" width="45" height="45" class="goods-thumb" />
				   	%s <span class="small-title">%s</span></a><br />
				   	商品编号：%s
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
			v.GoodsId, GetGoodsImageUrl(v.GoodsImage), v.GoodsName, v.SmallTitle, v.GoodsNo,
			FormatFloat(v.SalePrice), v.Num, FormatFloat(v.SalePrice*float32(v.Num)),
		))
	}

	byts.WriteString("</table>")

	return byts.String()
}
