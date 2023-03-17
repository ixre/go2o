/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-02-04 20:39
 * description :
 * history :
 */
package item

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
)

var (
	ErrInvalidTitle = domain.NewError("err_item_invalid_title", "商品标题不能包含特殊字符")
	ErrNoSuchSku    = domain.NewError(
		"err_item_no_such_item_sku", "商品SKU不存在")

	ErrIncorrectShopOfItem = domain.NewError(
		"err_item_not_bind_shop", "商品绑定的店铺非法")

	ErrEmptyReviewRemark = domain.NewError(
		"err_sale_empty_remark", "原因不能为空")

	ErrGoodsNum = domain.NewError(
		"err_goods_num", "商品数量错误")

	ErrOutOfSalePrice = domain.NewError(
		"out_of_sale_price", "超出商品售价")

	ErrOutOfStock = domain.NewError(
		"err_out_of_stock", "库存不足")

	ErrFullOfStock = domain.NewError(
		"err_full_of_stock", "商品\"%s\"已经售完")

	ErrInternalDisallow = domain.NewError(
		"err_sale_tag_internal_disallow", "不允许删除内置销售标签！")

	ErrCanNotDeleteItem = domain.NewError(
		"err_goods_can_not_delete_item", "已售出货品只允许下架。")

	ErrNotSetWholesalePrice = domain.NewError(
		"err_not_set_wholesale_price", "请先设置批发价格！")
)
