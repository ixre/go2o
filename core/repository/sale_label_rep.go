/**
 * Copyright 2015 @ z3q.net.
 * name : tag_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository

import (
	"errors"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
	saleImpl "go2o/core/domain/sale"
)

type saleLabelRep struct {
	db.Connector
}

func NewTagSaleRep(c db.Connector) sale.ISaleLabelRep {
	return &saleLabelRep{c}
}

// 创建销售标签
func (t *saleLabelRep) CreateSaleLabel(v *sale.Label) sale.ISaleLabel {
	if v != nil {
		return saleImpl.NewSaleLabel(v.MerchantId, v, t)
	}
	return nil
}

// 获取所有的销售标签
func (t *saleLabelRep) GetAllValueSaleLabels(merchantId int) []*sale.Label {
	arr := []*sale.Label{}
	t.Connector.GetOrm().Select(&arr, "mch_id=?", merchantId)
	return arr
}

// 获取销售标签值
func (t *saleLabelRep) GetValueSaleLabel(merchantId int, tagId int) *sale.Label {
	var v *sale.Label = new(sale.Label)
	err := t.Connector.GetOrm().GetBy(v, "mch_id=? AND id=?", merchantId, tagId)
	if err == nil {
		return v
	}
	return nil
}

// 获取销售标签
func (t *saleLabelRep) GetSaleLabel(merchantId int, id int) sale.ISaleLabel {
	return t.CreateSaleLabel(t.GetValueSaleLabel(merchantId, id))
}

// 保存销售标签
func (t *saleLabelRep) SaveSaleLabel(mchId int, v *sale.Label) (int, error) {
	v.MerchantId = mchId
	return orm.Save(t.GetOrm(), v, v.Id)
}

// 根据Code获取销售标签
func (t *saleLabelRep) GetSaleLabelByCode(merchantId int, code string) *sale.Label {
	var v *sale.Label = new(sale.Label)
	if t.GetOrm().GetBy(v, "mch_id=? AND tag_code=?", merchantId, code) == nil {
		return v
	}
	return nil
}

// 删除销售标签
func (t *saleLabelRep) DeleteSaleLabel(merchantId int, id int) error {
	_, err := t.GetOrm().Delete(&sale.Label{}, "mch_id=? AND id=?", merchantId, id)
	return err
}

// 获取商品
func (t *saleLabelRep) GetValueGoodsBySaleLabel(mchId,
	tagId int, sortBy string, begin, end int) []*valueobject.Goods {
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	arr := []*valueobject.Goods{}
	t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_goods INNER JOIN
	       gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.review_state=? AND gs_item.shelve_state=? AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_label t
			 ON t.id = g.sale_tag_id WHERE t.mch_id=? AND t.id=?) `+sortBy+`
			LIMIT ?,?`, enum.ReviewPass, item.ShelvesOn, mchId, tagId, begin, end)
	return arr
}

// 获取商品
func (t *saleLabelRep) GetPagedValueGoodsBySaleLabel(mchId,
	tagId int, sortBy string, begin, end int) (int, []*valueobject.Goods) {
	var total int
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	t.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods
	    INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.review_state=? AND gs_item.shelve_state=? AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id=? AND t.id=?)`), &total, enum.ReviewPass,
		item.ShelvesOn, mchId, tagId)
	arr := []*valueobject.Goods{}
	if total > 0 {
		t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_goods
         INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.review_state=? AND gs_item.shelve_state=? AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id=? AND t.id=?) `+sortBy+` LIMIT ?,?`,
			enum.ReviewPass, item.ShelvesOn,
			mchId, tagId, begin, end)
	}
	return total, arr
}

// 获取商品的销售标签
func (t *saleLabelRep) GetItemSaleLabels(itemId int) []*sale.Label {
	arr := []*sale.Label{}
	t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_sale_label WHERE id IN
	(SELECT sale_tag_id FROM gs_item_tag WHERE item_id=?) AND enabled=1`, itemId)
	return arr
}

// 清理商品的销售标签
func (t *saleLabelRep) CleanItemSaleLabels(itemId int) error {
	_, err := t.ExecNonQuery("DELETE FROM gs_item_tag WHERE item_id=?", itemId)
	return err
}

// 保存商品的销售标签
func (t *saleLabelRep) SaveItemSaleLabels(itemId int, tagIds []int) error {
	var err error
	if tagIds == nil {
		return errors.New("SaleLabel Ids can't be null.")
	}

	for _, v := range tagIds {
		_, err = t.ExecNonQuery("INSERT INTO gs_item_tag (item_id,sale_tag_id) VALUES(?,?)",
			itemId, v)
	}

	return err
}
