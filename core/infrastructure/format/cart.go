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
	var byts *bytes.Buffer = bytes.NewBufferString("")
	//byts.WriteString(`
	//	<table cellspacing="1" class="cart_details_table">
	//		<thead>
	//			<tr>
	//				<td><span class="t">商品</span></td>
	//				<td><span class="t">价格</span></td>
	//				<td><span class="t">数量</span></td>
	//				<td><span class="t">总价</span></td>
	//			</tr>
	//		</thead>
	//	`)

	for _, v := range c.Items {
		byts.WriteString(fmt.Sprintf(`
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
			</div>
		`,
			v.GoodsId, GetGoodsImageUrl(v.GoodsImage), v.GoodsName, v.Num, v.GoodsNo,
			FormatFloat(v.SalePrice), FormatFloat(v.SalePrice*float32(v.Num)),
		))
	}

	byts.WriteString("</table>")

	return byts.String()
}
