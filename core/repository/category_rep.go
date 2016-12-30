/**
 * Copyright 2015 @ z3q.net.
 * name : category_repo.go
 * author : jarryliu
 * date : 2016-06-04 13:01
 * description :
 * history :
 */
package repository

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	productImpl "go2o/core/domain/product"
	"go2o/core/infrastructure/format"
	"log"
	"sort"
)

var _ product.ICategoryRepo = new(categoryRepo)

type categoryRepo struct {
	db.Connector
	_valRepo     valueobject.IValueRepo
	_globService product.IGlobCatService
	_orm         orm.Orm
	storage      storage.Interface
}

func NewCategoryRepo(conn db.Connector, valRepo valueobject.IValueRepo,
	storage storage.Interface) product.ICategoryRepo {
	return &categoryRepo{
		Connector: conn,
		_orm:      conn.GetOrm(),
		_valRepo:  valRepo,
		storage:   storage,
	}
}

func (c *categoryRepo) GlobCatService() product.IGlobCatService {
	if c._globService == nil {
		c._globService = productImpl.NewCategoryManager(0, c, c._valRepo)
	}
	return c._globService
}

func (c *categoryRepo) getCategoryCacheKey(id int32) string {
	return fmt.Sprintf("go2o:rep:cat:c%d", id)
}

func (c *categoryRepo) SaveCategory(v *product.Category) (int32, error) {
	id, err := orm.I32(orm.Save(c.GetOrm(), v, int(v.ID)))
	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, "go2o:rep:cat:*")
	}
	return id, err
}

// 检查分类是否关联商品
func (c *categoryRepo) CheckGoodsContain(mchId, id int32) bool {
	num := 0
	//清理项
	c.Connector.ExecScalar(`SELECT COUNT(0) FROM pro_product WHERE cat_id IN
		(SELECT Id FROM cat_category WHERE mch_id=? AND id=?)`, &num, mchId, id)
	return num > 0
}

func (c *categoryRepo) DeleteCategory(mchId, id int32) error {
	//删除子类
	_, _, err := c.Connector.Exec("DELETE FROM cat_category WHERE parent_id=?",
		id)

	//删除分类
	_, _, err = c.Connector.Exec("DELETE FROM cat_category WHERE id=?",
		id)

	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, "go2o:rep:cat:*")
	}

	return err
}

func (c *categoryRepo) GetCategory(mchId, id int32) *product.Category {
	e := product.Category{}
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

func (c *categoryRepo) convertICategory(list product.CategoryList) []product.ICategory {
	sort.Sort(list)
	slice := make([]product.ICategory, len(list))
	for i, v := range list {
		slice[i] = c.GlobCatService().CreateCategory(v)
	}
	return slice
}

func (c *categoryRepo) redirectGetCats() []*product.Category {
	list := []*product.Category{}
	err := c.Connector.GetOrm().Select(&list, "")
	if err != nil {
		handleError(err)
	}
	return list
}

func (c *categoryRepo) GetCategories(mchId int32) []*product.Category {
	return c.redirectGetCats()
	//todo: cache
	//key := fmt.Sprintf("go2o:rep:cat:list9:%d", mchId)
	//list := []*product.Category{}
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

// 获取关联的品牌
func (c *categoryRepo) GetRelationBrands(idArr []int32) []*promodel.ProBrand {
	list := []*promodel.ProBrand{}
	if len(idArr) > 0 {
		err := c._orm.Select(&list, `id IN (SELECT brand_id FROM pro_model_brand
        WHERE pro_model IN (SELECT pro_model FROM cat_category WHERE id IN(`+
			format.IdArrJoinStr32(idArr)+`)))`)
		if err != nil && err != sql.ErrNoRows {
			log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
		}
	}
	return list
}
