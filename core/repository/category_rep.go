/**
 * Copyright 2015 @ z3q.net.
 * name : category_rep.go
 * author : jarryliu
 * date : 2016-06-04 13:01
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/algorithm/iterator"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	saleImpl "go2o/core/domain/sale"
)

var _ sale.ICategoryRep = new(categoryRep)

type categoryRep struct {
	db.Connector
	_valRep          valueobject.IValueRep
	_globCateManager sale.ICategoryManager
}

func NewCategoryRep(conn db.Connector, valRep valueobject.IValueRep) sale.ICategoryRep {
	return &categoryRep{
		Connector: conn,
		_valRep:   valRep,
	}
}

func (this *categoryRep) GetGlobManager() sale.ICategoryManager {
	if this._globCateManager == nil {
		this._globCateManager = saleImpl.NewCategoryManager(0, this, this._valRep)
	}
	return this._globCateManager
}

func (this *categoryRep) SaveCategory(v *sale.Category) (int, error) {
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

func (this *categoryRep) DeleteCategory(merchantId, id int) error {
	//删除子类
	_, _, err := this.Connector.Exec("DELETE FROM gs_category WHERE mch_id=? AND parent_id=?",
		merchantId, id)

	//删除分类
	_, _, err = this.Connector.Exec("DELETE FROM gs_category WHERE mch_id=? AND id=?",
		merchantId, id)

	//清理项
	this.Connector.Exec(`DELETE FROM gs_item WHERE Cid NOT IN
		(SELECT Id FROM gs_category WHERE mch_id=?)`, merchantId)

	return err
}

func (this *categoryRep) GetCategory(merchantId, id int) *sale.Category {
	var e *sale.Category = new(sale.Category)
	err := this.Connector.GetOrm().Get(id, e)
	if err == nil && e.MerchantId == merchantId {
		return e
	}
	return nil
}

func (this *categoryRep) GetCategories(merchantId int) sale.CategoryList {
	var e []*sale.Category = []*sale.Category{}
	err := this.Connector.GetOrm().Select(&e, "mch_id=? ORDER BY id ASC", merchantId)
	if err == nil {
		return e
	}
	return nil
}

// 获取与栏目相关的栏目
func (this *categoryRep) GetRelationCategories(merchantId, categoryId int) sale.CategoryList {
	var all []*sale.Category = this.GetCategories(merchantId)
	var newArr []*sale.Category = []*sale.Category{}
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
func (this *categoryRep) GetChildCategories(merchantId, categoryId int) sale.CategoryList {
	var all []*sale.Category = this.GetCategories(merchantId)
	var newArr []*sale.Category = []*sale.Category{}

	var cdt iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*sale.Category).ParentId == v.(*sale.Category).Id
	}
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.Category)
		if c.Id != categoryId {
			newArr = append(newArr, c)
		}
	}

	var arr []interface{} = make([]interface{}, len(all))
	for i, _ := range arr {
		arr[i] = all[i]
	}

	iterator.Walk(arr, &sale.Category{Id: categoryId}, cdt, start, nil, 1)

	return newArr
}
