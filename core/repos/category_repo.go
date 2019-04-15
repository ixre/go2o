/**
 * Copyright 2015 @ z3q.net.
 * name : category_repo.go
 * author : jarryliu
 * date : 2016-06-04 13:01
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
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
	return fmt.Sprintf("go2o:repo:cat:c%d", id)
}

func (c *categoryRepo) SaveCategory(v *product.Category) (int32, error) {
	id, err := orm.I32(orm.Save(c.GetOrm(), v, int(v.ID)))
	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, "go2o:repo:cat:list")
	}
	return id, err
}

// 检查分类是否关联商品
func (c *categoryRepo) CheckContainGoods(vendorId, catId int32) bool {
	num := 0
	if vendorId <= 0 {
		c.Connector.ExecScalar(`SELECT COUNT(0) FROM pro_product WHERE cat_id=?`, &num, catId)
	} else {
		panic(errors.New("暂时不支持商户绑定分类"))
	}
	return num > 0
}

func (c *categoryRepo) DeleteCategory(mchId, id int32) error {
	//删除子类
	_, _, err := c.Connector.Exec("DELETE FROM pro_category WHERE parent_id=?",
		id)

	//删除分类
	_, _, err = c.Connector.Exec("DELETE FROM pro_category WHERE id=?",
		id)

	// 清理缓存
	if err == nil {
		c.storage.Del(c.getCategoryCacheKey(id))
		PrefixDel(c.storage, "go2o:repo:cat:list")
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
	key := "go2o:repo:cat:list"
	list := []*product.Category{}
	jsonStr, err := c.storage.GetBytes(key)
	if err == nil {
		err = json.Unmarshal(jsonStr, &list)
	}
	if err != nil {
		handleError(err)
		err := c.Connector.GetOrm().Select(&list, "true ORDER BY sort_num DESC,id ASC")
		if err == nil {
			b, _ := json.Marshal(list)
			c.storage.Set(key, b)
		}
	}
	return list
}

// 获取关联的品牌
func (c *categoryRepo) GetRelationBrands(idArr []int32) []*promodel.ProBrand {
	list := []*promodel.ProBrand{}
	if len(idArr) > 0 {
		err := c._orm.Select(&list, `id IN (SELECT brand_id FROM pro_model_brand
        WHERE pro_model IN (SELECT distinct pro_model FROM pro_category WHERE id IN(`+
			format.I32ArrStrJoin(idArr)+`)))`)
		if err != nil && err != sql.ErrNoRows {
			log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
		}
	}
	return list
}
