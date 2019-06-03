/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-04 20:39
 * description :
 * history :
 */
package item

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchSku *domain.DomainError = domain.NewError(
		"err_item_no_such_item_sku", "商品SKU不存在")

	ErrNotBindShop *domain.DomainError = domain.NewError(
		"err_item_not_bind_shop", "请选择商品上架的商铺")

	ErrIncorrectShopOfItem *domain.DomainError = domain.NewError(
		"err_item_not_bind_shop", "商品绑定的商铺非法")

	ErrEmptyReviewRemark *domain.DomainError = domain.NewError(
		"err_sale_empty_remark", "原因不能为空")

	ErrGoodsNum *domain.DomainError = domain.NewError(
		"err_goods_num", "商品数量错误")

	ErrOutOfSalePrice *domain.DomainError = domain.NewError(
		"out_of_sale_price", "超出商品售价")

	ErrOutOfStock *domain.DomainError = domain.NewError(
		"err_out_of_stock", "库存不足")

	ErrFullOfStock *domain.DomainError = domain.NewError(
		"err_full_of_stock", "商品已经售完")

	ErrInternalDisallow *domain.DomainError = domain.NewError(
		"err_sale_tag_internal_disallow", "不允许删除内置销售标签！")

	ErrCanNotDeleteItem *domain.DomainError = domain.NewError(
		"err_goods_can_not_delete_item", "已售出货品只允许下架。")

	ErrNotSetWholesalePrice *domain.DomainError = domain.NewError(
		"err_not_set_wholesale_price", "请先设置批发价格！")
)
