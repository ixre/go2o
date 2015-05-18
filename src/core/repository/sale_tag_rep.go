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
    saleImpl "go2o/src/core/domain/sale"
    "github.com/atnet/gof/db"
    "go2o/src/core/domain/interface/sale"
)

type SaleTagRep struct {
    db.Connector
}



// 创建销售标签
func (this *SaleTagRep) CreateSaleTag(v *sale.ValueSaleTag)sale.ISaleTag{
    return saleImpl.NewSaleTag(v.PartnerId,v,this)
}

// 获取销售标签值
func (this *SaleTagRep)GetValueSaleTag(partnerId int,tagId int)*sale.ValueSaleTag{
    var v *sale.ValueSaleTag
    err := this.Connector.GetOrm().GetBy(v,"partner_id=? AND id=?",partnerId,tagId)
    if err == nil{
        return v
    }
    return nil
}

// 获取销售标签
func (this *SaleTagRep)GetSaleTag(partnerId int,tagId int)sale.ISaleTag{
    return this.CreateSaleTag(this.GetValueSaleTag(partnerId,tagId))
}

// 保存销售标签
func (this *SaleTagRep)SaveSaleTag(partnerId int,v *sale.ValueSaleTag)(int,error) {
    orm := this.GetOrm()
    var err error
    v.PartnerId = partnerId
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    }else {
        _, _, err = orm.Save(nil, v)
        this.Connector.ExecScalar("SELECT MAX(id) FROM pt_sale_tag WHERE partner_id=?", &v.Id, partnerId)
    }
    return v.Id, err
}

// 获取商品
func (this *SaleTagRep)GetValueGoods(partnerId,tagId,begin,end int)[]*sale.ValueGoods {
    //todo:
    arr := []*sale.ValueGoods{}
    this.Connector.GetOrm().SelectByQuery(&arr,
    "SELECT * FROM mm_member WHERE id IN (SELECT member_id FROM mm_relation WHERE invi_member_id=?)", partnerId, tagId, begin, end)
    return arr
}