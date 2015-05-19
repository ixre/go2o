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
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/sale"
	saleImpl "go2o/src/core/domain/sale"
)

type SaleTagRep struct {
	db.Connector
}

func NewTagSaleRep(c db.Connector) sale.ISaleTagRep {
	return &SaleTagRep{c}
}

// 创建销售标签
func (this *SaleTagRep) CreateSaleTag(v *sale.ValueSaleTag) sale.ISaleTag {
	return saleImpl.NewSaleTag(v.PartnerId, v, this)
}

// 获取销售标签值
func (this *SaleTagRep) GetValueSaleTag(partnerId int, tagId int) *sale.ValueSaleTag {
	var v *sale.ValueSaleTag
	err := this.Connector.GetOrm().GetBy(v, "partner_id=? AND id=?", partnerId, tagId)
	if err == nil {
		return v
	}
	return nil
}

// 获取销售标签
func (this *SaleTagRep) GetSaleTag(partnerId int, tagId int) sale.ISaleTag {
	return this.CreateSaleTag(this.GetValueSaleTag(partnerId, tagId))
}

// 保存销售标签
func (this *SaleTagRep) SaveSaleTag(partnerId int, v *sale.ValueSaleTag) (int, error) {
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
func (this *SaleTagRep) GetSaleTagByCode(partnerId int, code string) *sale.ValueSaleTag {
	var v *sale.ValueSaleTag = new(sale.ValueSaleTag)
	if this.GetOrm().GetBy(v, "partner_id=? AND tag_code=?", partnerId, code) == nil {
		return v
	}
	return nil
}

// 删除销售标签
func (this *SaleTagRep) DeleteSaleTag(partnerId int, id int) error {
	_, err := this.GetOrm().Delete(&sale.ValueSaleTag{}, "partner_id=? AND id=?", partnerId, id)
	return err
}

// 获取商品
func (this *SaleTagRep) GetValueGoods(partnerId, tagId, begin, end int) []*sale.ValueGoods {
	//todo:
	arr := []*sale.ValueGoods{}
	this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM mm_member WHERE id IN (SELECT member_id FROM mm_relation WHERE invi_member_id=?)",
		partnerId, tagId, begin, end)
	return arr
}
