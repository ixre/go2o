/**
 * Copyright 2015 @ 56x.net.
 * name : tag_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	"errors"
	"fmt"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	itemImpl "github.com/ixre/go2o/core/domain/item"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type saleLabelRepo struct {
	db.Connector
	_service item.ILabelService
	valRepo  valueobject.IValueRepo
	o        orm.Orm
}

func NewTagSaleRepo(o orm.Orm,
	valRepo valueobject.IValueRepo) item.ISaleLabelRepo {
	return &saleLabelRepo{
		Connector: o.Connector(),
		o:         o,
		valRepo:   valRepo,
	}
}

// 获取商品标签服务
func (t *saleLabelRepo) LabelService() item.ILabelService {
	if t._service == nil {
		t._service = itemImpl.NewLabelManager(0, t, t.valRepo)
	}
	return t._service
}

// 创建销售标签
func (t *saleLabelRepo) CreateSaleLabel(v *item.Label) item.ISaleLabel {
	if v != nil {
		return itemImpl.NewSaleLabel(v.MerchantId, v, t)
	}
	return nil
}

// 获取所有的销售标签
func (t *saleLabelRepo) GetAllValueSaleLabels(mchId int64) []*item.Label {
	arr := []*item.Label{}
	t.o.Select(&arr, "mch_id= $1", mchId)
	return arr
}

// 获取销售标签值
func (t *saleLabelRepo) GetValueSaleLabel(mchId int64, tagId int32) *item.Label {
	var v = new(item.Label)
	err := t.o.GetBy(v, "mch_id= $1 AND id= $2", mchId, tagId)
	if err == nil {
		return v
	}
	return nil
}

// 获取销售标签
func (t *saleLabelRepo) GetSaleLabel(mchId int64, id int32) item.ISaleLabel {
	return t.CreateSaleLabel(t.GetValueSaleLabel(mchId, id))
}

// 保存销售标签
func (t *saleLabelRepo) SaveSaleLabel(mchId int64, v *item.Label) (int32, error) {
	v.MerchantId = mchId
	return orm.I32(orm.Save(t.o, v, int(v.Id)))
}

// 根据Code获取销售标签
func (t *saleLabelRepo) GetSaleLabelByCode(mchId int64, code string) *item.Label {
	var v = new(item.Label)
	if t.o.GetBy(v, "mch_id= $1 AND tag_code= $2", mchId, code) == nil {
		return v
	}
	return nil
}

// 删除销售标签
func (t *saleLabelRepo) DeleteSaleLabel(mchId int64, id int32) error {
	_, err := t.o.Delete(&item.Label{}, "mch_id= $1 AND id= $2", mchId, id)
	return err
}

// 获取商品
func (t *saleLabelRepo) GetValueGoodsBySaleLabel(mchId int64, tagId int32,
	sortBy string, begin, end int) []*valueobject.Goods {
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	arr := []*valueobject.Goods{}
	t.o.SelectByQuery(&arr, `SELECT * FROM item_info INNER JOIN
	       product ON product.id = item_info.product_id
		 WHERE product.review_state= $1 AND product.shelve_state= $2 AND product.id IN (
			SELECT g.item_id FROM product_tag g INNER JOIN gs_sale_label t
			 ON t.id = g.sale_tag_id WHERE t.mch_id= $3 AND t.id= $4) `+sortBy+`
			LIMIT $6 OFFSET $5`, enum.ReviewPass, item.ShelvesOn, mchId, tagId, begin, end)
	return arr
}

// 获取商品
func (t *saleLabelRepo) GetPagedValueGoodsBySaleLabel(mchId int64, tagId int32,
	sortBy string, begin, end int) (int, []*valueobject.Goods) {
	var total int
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	t.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
	    INNER JOIN product ON product.id = item_info.product_id
		 WHERE product.review_state= $1 AND product.shelve_state= $2 AND product.id IN (
			SELECT g.item_id FROM product_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id= $3 AND t.id= $4)`), &total, enum.ReviewPass,
		item.ShelvesOn, mchId, tagId)
	var arr []*valueobject.Goods
	if total > 0 {
		t.o.SelectByQuery(&arr, `SELECT * FROM item_info
         INNER JOIN product ON product.id = item_info.product_id
		 WHERE product.review_state= $1 AND product.shelve_state= $2 AND product.id IN (
			SELECT g.item_id FROM product_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id= $3 AND t.id= $4) `+sortBy+` LIMIT $6 OFFSET $5`,
			enum.ReviewPass, item.ShelvesOn,
			mchId, tagId, begin, end)
	}
	return total, arr
}

// 获取商品的销售标签
func (t *saleLabelRepo) GetItemSaleLabels(itemId int32) []*item.Label {
	arr := []*item.Label{}
	t.o.SelectByQuery(&arr, `SELECT * FROM gs_sale_label WHERE id IN
	(SELECT sale_tag_id FROM product_tag WHERE item_id= $1) AND enabled=1`, itemId)
	return arr
}

// 清理商品的销售标签
func (t *saleLabelRepo) CleanItemSaleLabels(itemId int32) error {
	_, err := t.ExecNonQuery("DELETE FROM product_tag WHERE item_id= $1", itemId)
	return err
}

// 保存商品的销售标签
func (t *saleLabelRepo) SaveItemSaleLabels(itemId int32, tagIds []int) error {
	var err error
	if tagIds == nil {
		return errors.New("SaleLabel Ids can't be null.")
	}

	for _, v := range tagIds {
		_, err = t.ExecNonQuery("INSERT INTO product_tag (item_id,sale_tag_id) VALUES($1,$2)",
			itemId, v)
	}

	return err
}
