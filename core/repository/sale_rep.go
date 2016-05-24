/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:09
 * description :
 * history :
 */

package repository

import (
	"fmt"
	"github.com/jsix/gof/algorithm/iterator"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	saleImpl "go2o/core/domain/sale"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/log"
)

var _ sale.ISaleRep = new(saleRep)

type saleRep struct {
	db.Connector
	_cache    map[int]sale.ISale
	_tagRep   sale.ISaleTagRep
	_promRep  promotion.IPromotionRep
	_goodsRep sale.IGoodsRep
}

func NewSaleRep(c db.Connector, saleTagRep sale.ISaleTagRep, goodsRep sale.IGoodsRep, promRep promotion.IPromotionRep) sale.ISaleRep {
	return (&saleRep{
		Connector: c,
		_tagRep:   saleTagRep,
		_promRep:  promRep,
		_goodsRep: goodsRep,
	}).init()
}

func (this *saleRep) init() sale.ISaleRep {
	this._cache = make(map[int]sale.ISale)
	return this
}

func (this *saleRep) GetSale(merchantId int) sale.ISale {
	v, ok := this._cache[merchantId]
	if !ok {
		v = saleImpl.NewSale(merchantId, this, this._goodsRep, this._tagRep, this._promRep)
		this._cache[merchantId] = v
	}
	return v
}

func (this *saleRep) GetValueItem(merchantId, itemId int) *sale.ValueItem {
	var e *sale.ValueItem = new(sale.ValueItem)
	err := this.Connector.GetOrm().GetByQuery(e, `select * FROM gs_item
			INNER JOIN gs_category c ON c.id = gs_item.category_id WHERE gs_item.id=?
			AND c.merchant_id=?`, itemId, merchantId)
	if err != nil {
		return nil
	}
	return e
}

func (this *saleRep) GetItemByIds(ids ...int) ([]*sale.ValueItem, error) {
	//todo: merchantId
	var items []*sale.ValueItem

	//todo:改成database/sql方式，不使用orm
	err := this.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_item WHERE id IN (`+format.GetCategoryIdStr(ids)+`)`)

	return items, err
}

func (this *saleRep) SaveValueItem(v *sale.ValueItem) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, id, err := orm.Save(nil, v)
		return int(id), err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

func (this *saleRep) GetPagedOnShelvesItem(merchantId int, catIds []int, start, end int) (total int, e []*sale.ValueItem) {
	var sql string

	var catIdStr string = format.GetCategoryIdStr(catIds)
	sql = fmt.Sprintf(`SELECT * FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, merchantId, catIdStr, start, (end - start))

	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1`, merchantId, catIdStr), &total)

	e = []*sale.ValueItem{}
	this.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (this *saleRep) GetItemSaleNum(merchantId int, id int) int {
	var num int
	this.Connector.ExecScalar(`SELECT SUM(sale_num) FROM gs_goods WHERE item_id=?`, &num, id)
	return num
}

func (this *saleRep) DeleteItem(merchantId, itemId int) error {
	_, _, err := this.Connector.Exec(`
		DELETE f FROM gs_item AS f
		INNER JOIN gs_category AS c ON f.category_id=c.id
		WHERE f.id=? AND c.merchant_id=?`, itemId, merchantId)
	return err
}

func (this *saleRep) SaveCategory(v *sale.ValueCategory) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, _, err := orm.Save(nil, v)
		if err == nil {
			this.Connector.ExecScalar(`SELECT MAX(id) FROM gs_category`, &v.Id)
		}
		return v.Id, err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

func (this *saleRep) DeleteCategory(merchantId, id int) error {
	//删除子类
	_, _, err := this.Connector.Exec("DELETE FROM gs_category WHERE merchant_id=? AND parent_id=?",
		merchantId, id)

	//删除分类
	_, _, err = this.Connector.Exec("DELETE FROM gs_category WHERE merchant_id=? AND id=?",
		merchantId, id)

	//清理项
	this.Connector.Exec(`DELETE FROM gs_item WHERE Cid NOT IN
		(SELECT Id FROM gs_category WHERE merchant_id=?)`, merchantId)

	return err
}

func (this *saleRep) GetCategory(merchantId, id int) *sale.ValueCategory {
	var e *sale.ValueCategory = new(sale.ValueCategory)
	err := this.Connector.GetOrm().Get(id, e)
	log.Println("--", e)
	if err == nil && e.MerchantId == merchantId {
		return e
	}
	return nil
}

func (this *saleRep) GetCategories(merchantId int) sale.CategoryList {
	var e []*sale.ValueCategory = []*sale.ValueCategory{}
	err := this.Connector.GetOrm().Select(&e, "merchant_id=? ORDER BY id ASC", merchantId)
	if err == nil {
		return e
	}
	return nil
}

// 获取与栏目相关的栏目
func (this *saleRep) GetRelationCategories(merchantId, categoryId int) sale.CategoryList {
	var all []*sale.ValueCategory = this.GetCategories(merchantId)
	var newArr []*sale.ValueCategory = []*sale.ValueCategory{}
	var isMatch bool
	var pid int
	var l int = len(all)

	for i := 0; i < l; i++ {
		if !isMatch && all[i].Id == categoryId {
			isMatch = true
			pid = all[i].ParentId
			newArr = append(newArr, all[i])
			i = -1
		} else {
			if all[i].Id == pid {
				newArr = append(newArr, all[i])
				pid = all[i].ParentId
				i = -1
				if pid == 0 {
					break
				}
			}
		}
	}
	return newArr
}

// 获取子栏目
func (this *saleRep) GetChildCategories(merchantId, categoryId int) sale.CategoryList {
	var all []*sale.ValueCategory = this.GetCategories(merchantId)
	var newArr []*sale.ValueCategory = []*sale.ValueCategory{}

	var cdt iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*sale.ValueCategory).ParentId == v.(*sale.ValueCategory).Id
	}
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.ValueCategory)
		if c.Id != categoryId {
			newArr = append(newArr, c)
		}
	}

	var arr []interface{} = make([]interface{}, len(all))
	for i, _ := range arr {
		arr[i] = all[i]
	}

	iterator.Walk(arr, &sale.ValueCategory{Id: categoryId}, cdt, start, nil, 1)

	return newArr
}

// 保存快照
func (this *saleRep) SaveSnapshot(v *sale.GoodsSnapshot) (int, error) {
	var id int
	_, _, err := this.Connector.GetOrm().Save(nil, v)
	if err == nil {
		err = this.Connector.ExecScalar(`SELECT MAX(id) FROM gs_snapshot where goods_id=?`, &id, v.GoodsId)
	}

	return id, err
}

// 获取最新的商品快照
func (this *saleRep) GetLatestGoodsSnapshot(goodsId int) *sale.GoodsSnapshot {
	var e *sale.GoodsSnapshot = new(sale.GoodsSnapshot)
	if this.Connector.GetOrm().GetBy(e, "goods_id=? ORDER BY id DESC", goodsId) == nil {
		return e
	}
	return nil
}

// 获取指定的商品快照
func (this *saleRep) GetGoodsSnapshot(id int) *sale.GoodsSnapshot {
	var e *sale.GoodsSnapshot = new(sale.GoodsSnapshot)
	err := this.Connector.GetOrm().Get(id, e)
	if err != nil {
		log.Error(err)
		e = nil
	}
	return e
}

// 根据Key获取商品快照
func (this *saleRep) GetGoodsSnapshotByKey(key string) *sale.GoodsSnapshot {
	var e *sale.GoodsSnapshot = new(sale.GoodsSnapshot)
	err := this.Connector.GetOrm().GetBy(e, "key=?", key)
	if err != nil {
		log.Error(err)
		e = nil
	}
	return e
}
