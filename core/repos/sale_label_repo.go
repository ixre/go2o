/**
 * Copyright 2015 @ z3q.net.
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
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/valueobject"
	itemImpl "go2o/core/domain/item"
)

type saleLabelRepo struct {
	db.Connector
	_service item.ILabelService
	valRepo  valueobject.IValueRepo
}

func NewTagSaleRepo(c db.Connector,
	valRepo valueobject.IValueRepo) item.ISaleLabelRepo {
	return &saleLabelRepo{
		Connector: c,
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
func (t *saleLabelRepo) GetAllValueSaleLabels(mchId int32) []*item.Label {
	arr := []*item.Label{}
	t.Connector.GetOrm().Select(&arr, "mch_id=?", mchId)
	return arr
}

// 获取销售标签值
func (t *saleLabelRepo) GetValueSaleLabel(mchId int32, tagId int32) *item.Label {
	var v *item.Label = new(item.Label)
	err := t.Connector.GetOrm().GetBy(v, "mch_id=? AND id=?", mchId, tagId)
	if err == nil {
		return v
	}
	return nil
}

// 获取销售标签
func (t *saleLabelRepo) GetSaleLabel(mchId int32, id int32) item.ISaleLabel {
	return t.CreateSaleLabel(t.GetValueSaleLabel(mchId, id))
}

// 保存销售标签
func (t *saleLabelRepo) SaveSaleLabel(mchId int32, v *item.Label) (int32, error) {
	v.MerchantId = mchId
	return orm.I32(orm.Save(t.GetOrm(), v, int(v.Id)))
}

// 根据Code获取销售标签
func (t *saleLabelRepo) GetSaleLabelByCode(mchId int32, code string) *item.Label {
	var v *item.Label = new(item.Label)
	if t.GetOrm().GetBy(v, "mch_id=? AND tag_code=?", mchId, code) == nil {
		return v
	}
	return nil
}

// 删除销售标签
func (t *saleLabelRepo) DeleteSaleLabel(mchId int32, id int32) error {
	_, err := t.GetOrm().Delete(&item.Label{}, "mch_id=? AND id=?", mchId, id)
	return err
}

// 获取商品
func (t *saleLabelRepo) GetValueGoodsBySaleLabel(mchId, tagId int32,
	sortBy string, begin, end int) []*valueobject.Goods {
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	arr := []*valueobject.Goods{}
	t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM item_info INNER JOIN
	       pro_product ON pro_product.id = item_info.product_id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=? AND pro_product.id IN (
			SELECT g.item_id FROM pro_product_tag g INNER JOIN gs_sale_label t
			 ON t.id = g.sale_tag_id WHERE t.mch_id=? AND t.id=?) `+sortBy+`
			LIMIT ?,?`, enum.ReviewPass, item.ShelvesOn, mchId, tagId, begin, end)
	return arr
}

// 获取商品
func (t *saleLabelRepo) GetPagedValueGoodsBySaleLabel(mchId, tagId int32,
	sortBy string, begin, end int) (int, []*valueobject.Goods) {
	var total int
	if len(sortBy) > 0 {
		sortBy = "ORDER BY " + sortBy
	}
	t.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
	    INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=? AND pro_product.id IN (
			SELECT g.item_id FROM pro_product_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id=? AND t.id=?)`), &total, enum.ReviewPass,
		item.ShelvesOn, mchId, tagId)
	arr := []*valueobject.Goods{}
	if total > 0 {
		t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=? AND pro_product.id IN (
			SELECT g.item_id FROM pro_product_tag g INNER JOIN gs_sale_label t ON t.id = g.sale_tag_id
			WHERE t.mch_id=? AND t.id=?) `+sortBy+` LIMIT ?,?`,
			enum.ReviewPass, item.ShelvesOn,
			mchId, tagId, begin, end)
	}
	return total, arr
}

// 获取商品的销售标签
func (t *saleLabelRepo) GetItemSaleLabels(itemId int32) []*item.Label {
	arr := []*item.Label{}
	t.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM gs_sale_label WHERE id IN
	(SELECT sale_tag_id FROM pro_product_tag WHERE item_id=?) AND enabled=1`, itemId)
	return arr
}

// 清理商品的销售标签
func (t *saleLabelRepo) CleanItemSaleLabels(itemId int32) error {
	_, err := t.ExecNonQuery("DELETE FROM pro_product_tag WHERE item_id=?", itemId)
	return err
}

// 保存商品的销售标签
func (t *saleLabelRepo) SaveItemSaleLabels(itemId int32, tagIds []int) error {
	var err error
	if tagIds == nil {
		return errors.New("SaleLabel Ids can't be null.")
	}

	for _, v := range tagIds {
		_, err = t.ExecNonQuery("INSERT INTO pro_product_tag (item_id,sale_tag_id) VALUES(?,?)",
			itemId, v)
	}

	return err
}
