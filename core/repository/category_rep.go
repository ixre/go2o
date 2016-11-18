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
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	saleImpl "go2o/core/domain/sale"
	"sort"
)

var _ sale.ICategoryRep = new(categoryRep)

type categoryRep struct {
	db.Connector
	_valRep          valueobject.IValueRep
	_globCateManager sale.ICategoryManager
	storage          storage.Interface
}

func NewCategoryRep(conn db.Connector, valRep valueobject.IValueRep,
	storage storage.Interface) sale.ICategoryRep {
	return &categoryRep{
		Connector: conn,
		_valRep:   valRep,
		storage:   storage,
	}
}

func (c *categoryRep) GetGlobManager() sale.ICategoryManager {
	if c._globCateManager == nil {
		c._globCateManager = saleImpl.NewCategoryManager(0, c, c._valRep)
	}
	return c._globCateManager
}

func (c *categoryRep) getCategoryCacheKey(id int) string {
	return fmt.Sprintf("go2o:rep:cat:c:%d", id)
}

func (c *categoryRep) SaveCategory(v *sale.Category) (int64, error) {
	id, err := orm.Save(c.GetOrm(), v, v.Id)
	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, fmt.Sprintf("go2o:rep:cat:%d:*", v.MerchantId))
	}
	return id, err
}

// 检查分类是否关联商品
func (c *categoryRep) CheckGoodsContain(mchId, id int64) bool {
	num := 0
	//清理项
	c.Connector.ExecScalar(`SELECT COUNT(0) FROM gs_item WHERE category_id IN
		(SELECT Id FROM gs_category WHERE mch_id=? AND id=?)`, &num, mchId, id)
	return num > 0
}

func (c *categoryRep) DeleteCategory(mchId, id int64) error {
	//删除子类
	_, _, err := c.Connector.Exec("DELETE FROM gs_category WHERE mch_id=? AND parent_id=?",
		mchId, id)

	//删除分类
	_, _, err = c.Connector.Exec("DELETE FROM gs_category WHERE mch_id=? AND id=?",
		mchId, id)

	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, fmt.Sprintf("go2o:rep:cat:%d:*", mchId))
	}

	return err
}

func (c *categoryRep) GetCategory(mchId, id int) *sale.Category {
	e := sale.Category{}
	key := c.getCategoryCacheKey(id)
	if c.storage.Get(key, &e) != nil {
		err := c.Connector.GetOrm().Get(id, &e)
		if err != nil {
			return nil
		}
		c.storage.Set(key, &e)
	}
	return &e
}

// 创建分类
func (c *categoryRep) CreateCategory(v *sale.Category) sale.ICategory {
	return saleImpl.NewCategory(c, v)
}

func (c *categoryRep) convertICategory(list sale.CategoryList) []sale.ICategory {
	sort.Sort(list)
	slice := make([]sale.ICategory, len(list))
	for i, v := range list {
		slice[i] = c.CreateCategory(v)
	}
	return slice
}

func (c *categoryRep) redirectGetCats(mchId int) []*sale.Category {
	list := []*sale.Category{}
	err := c.Connector.GetOrm().Select(&list, "mch_id=? ORDER BY id ASC", mchId)
	if err != nil {
		handleError(err)
	}
	return list
}

func (c *categoryRep) GetCategories(mchId int64) []*sale.Category {
	return c.redirectGetCats(mchId)
	//todo: cache
	//key := fmt.Sprintf("go2o:rep:cat:list9:%d", mchId)
	//list := []*sale.Category{}
	//if err := c.storage.Get(key, &list);err != nil {
	//    handleError(err)
	//    err := c.Connector.GetOrm().Select(&list, "mch_id=? ORDER BY id ASC", mchId)
	//    if err == nil {
	//        c.storage.SetExpire(key,list, DefaultCacheSeconds)
	//    } else {
	//        handleError(err)
	//    }
	//}
	//return list
}
