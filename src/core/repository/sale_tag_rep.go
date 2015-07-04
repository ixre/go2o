/**
 * Copyright 2015 @ S1N1 Team.
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
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/valueobject"
	saleImpl "go2o/src/core/domain/sale"
)

type saleTagRep struct {
	db.Connector
}

func NewTagSaleRep(c db.Connector) sale.ISaleTagRep {
	return &saleTagRep{c}
}

// 创建销售标签
func (this *saleTagRep) CreateSaleTag(v *sale.ValueSaleTag) sale.ISaleTag {
	if v != nil {
		return saleImpl.NewSaleTag(v.PartnerId, v, this)
	}
	return nil
}

// 获取所有的销售标签
func (this *saleTagRep) GetAllValueSaleTags(partnerId int) []*sale.ValueSaleTag {
	arr := []*sale.ValueSaleTag{}
	this.Connector.GetOrm().Select(&arr, "partner_id=?", partnerId)
	return arr
}

// 获取销售标签值
func (this *saleTagRep) GetValueSaleTag(partnerId int, tagId int) *sale.ValueSaleTag {
	var v *sale.ValueSaleTag = new(sale.ValueSaleTag)
	err := this.Connector.GetOrm().GetBy(v, "partner_id=? AND id=?", partnerId, tagId)
	if err == nil {
		return v
	}
	return nil
}

// 获取销售标签
func (this *saleTagRep) GetSaleTag(partnerId int, id int) sale.ISaleTag {
	return this.CreateSaleTag(this.GetValueSaleTag(partnerId, id))
}

// 保存销售标签
func (this *saleTagRep) SaveSaleTag(partnerId int, v *sale.ValueSaleTag) (int, error) {
	orm := this.GetOrm()
	var err error
	v.PartnerId = partnerId
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM gs_sale_tag WHERE partner_id=?", &v.Id, partnerId)
	}
	return v.Id, err
}

// 根据Code获取销售标签
func (this *saleTagRep) GetSaleTagByCode(partnerId int, code string) *sale.ValueSaleTag {
	var v *sale.ValueSaleTag = new(sale.ValueSaleTag)
	if this.GetOrm().GetBy(v, "partner_id=? AND tag_code=?", partnerId, code) == nil {
		return v
	}
	return nil
}

// 删除销售标签
func (this *saleTagRep) DeleteSaleTag(partnerId int, id int) error {
	_, err := this.GetOrm().Delete(&sale.ValueSaleTag{}, "partner_id=? AND id=?", partnerId, id)
	return err
}

// 获取商品
func (this *saleTagRep) GetValueGoodsBySaleTag(partnerId, tagId, begin, end int) []*valueobject.Goods {
	arr := []*valueobject.Goods{}
	this.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.state=1  AND gs_item.on_shelves=1 AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_tag t ON t.id = g.sale_tag_id
			WHERE t.partner_id=? AND t.id=?) LIMIT ?,?`, partnerId, tagId, begin, end)
	return arr
}

// 获取商品
func (this *saleTagRep) GetPagedValueGoodsBySaleTag(partnerId, tagId, begin, end int) (int, []*valueobject.Goods) {
	var total int
	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.state=1  AND gs_item.on_shelves=1 AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_tag t ON t.id = g.sale_tag_id
			WHERE t.partner_id=? AND t.id=?)`), &total, partnerId, tagId)
	arr := []*valueobject.Goods{}
	if total > 0 {
		this.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 WHERE gs_item.state=1  AND gs_item.on_shelves=1 AND gs_item.id IN (
			SELECT g.item_id FROM gs_item_tag g INNER JOIN gs_sale_tag t ON t.id = g.sale_tag_id
			WHERE t.partner_id=? AND t.id=?) LIMIT ?,?`, partnerId, tagId, begin, end)
	}
	return total, arr
}

// 获取商品的销售标签
func (this *saleTagRep) GetItemSaleTags(itemId int) []*sale.ValueSaleTag {
	arr := []*sale.ValueSaleTag{}
	this.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_sale_tag WHERE id IN
	(SELECT sale_tag_id FROM gs_item_tag WHERE item_id=?) AND enabled=1`, itemId)
	return arr
}

// 清理商品的销售标签
func (this *saleTagRep) CleanItemSaleTags(itemId int) error {
	_, err := this.ExecNonQuery("DELETE FROM gs_item_tag WHERE item_id=?", itemId)
	return err
}

// 保存商品的销售标签
func (this *saleTagRep) SaveItemSaleTags(itemId int, tagIds []int) error {
	var err error
	if tagIds == nil {
		return errors.New("SaleTag Ids can't be null.")
	}

	for _, v := range tagIds {
		_, err = this.ExecNonQuery("INSERT INTO gs_item_tag (item_id,sale_tag_id) VALUES(?,?)",
			itemId, v)
	}

	return err
}
