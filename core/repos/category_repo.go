/**
 * Copyright 2015 @ 56x.net.
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
	"github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/registry"
	productImpl "github.com/ixre/go2o/core/domain/product"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"log"
	"sort"
)

var _ product.ICategoryRepo = new(categoryRepo)
var categoryPrefix = "go2o:gb:repo:cat:list_"

type categoryRepo struct {
	db.Connector
	registryRepo registry.IRegistryRepo
	_globService product.IGlobCatService
	o            orm.Orm
	storage      storage.Interface
}

func NewCategoryRepo(o orm.Orm, registryRepo registry.IRegistryRepo,
	storage storage.Interface) product.ICategoryRepo {
	return &categoryRepo{
		Connector:    o.Connector(),
		o:            o,
		registryRepo: registryRepo,
		storage:      storage,
	}
}

func (c *categoryRepo) GlobCatService() product.IGlobCatService {
	if c._globService == nil {
		c._globService = productImpl.NewCategoryManager(0, c, c.registryRepo)
	}
	return c._globService
}

func (c *categoryRepo) getCategoryCacheKey(id int) string {
	return fmt.Sprintf("go2o:repo:cat:c%d", id)
}

func (c *categoryRepo) SaveCategory(v *product.Category) (int, error) {
	id, err := orm.Save(c.o, v, int(v.Id))
	// 清理缓存
	if err == nil {
		c.storage.Delete(c.getCategoryCacheKey(id))
		c.storage.DeleteWith(categoryPrefix)
	}
	return id, err
}

// 检查分类是否关联商品
func (c *categoryRepo) CheckContainGoods(vendorId int64, catId int) bool {
	num := 0
	if vendorId <= 0 {
		c.Connector.ExecScalar(`SELECT COUNT(0) FROM product WHERE cat_id= $1`, &num, catId)
	} else {
		panic(errors.New("暂时不支持商户绑定分类"))
	}
	return num > 0
}

func (c *categoryRepo) DeleteCategory(mchId int64, id int) error {
	//删除子类
	_, err := c.Connector.ExecNonQuery("DELETE FROM product_category WHERE parent_id= $1",
		id)

	//删除分类
	_, err = c.Connector.ExecNonQuery("DELETE FROM product_category WHERE id= $1",
		id)

	// 清理缓存
	if err == nil {
		c.storage.Delete(c.getCategoryCacheKey(id))
		c.storage.DeleteWith(categoryPrefix)
	}

	return err
}

func (c *categoryRepo) GetCategory(mchId, id int) *product.Category {
	e := product.Category{}
	key := c.getCategoryCacheKey(id)
	if c.storage.Get(key, &e) != nil {
		err := c.o.Get(id, &e)
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
	var list []*product.Category
	err := c.o.Select(&list, "")
	if err != nil {
		handleError(err)
	}
	return list
}

func (c *categoryRepo) GetCategories(mchId int) []*product.Category {
	key := categoryPrefix + "data"
	var list []*product.Category
	jsonStr, err := c.storage.GetBytes(key)
	if err == nil {
		//println("---",string(jsonStr))
		err = json.Unmarshal(jsonStr, &list)
	}
	if err != nil {
		handleError(err)
		err := c.o.Select(&list, "1=1 ORDER BY sort_num DESC,id ASC")
		if err == nil {
			b, _ := json.Marshal(list)
			c.storage.Set(key, b)
		}
	}
	return list
}

// 获取关联的品牌
func (c *categoryRepo) GetRelationBrands(idArr []int) []*promodel.ProductBrand {
	var list []*promodel.ProductBrand
	if len(idArr) > 0 {
		err := c.o.Select(&list, `id IN (SELECT brand_id FROM product_model_brand
        WHERE pro_model IN (SELECT distinct pro_model FROM product_category WHERE id IN(`+
			format.IntArrStrJoin(idArr)+`)))`)
		if err != nil && err != sql.ErrNoRows {
			log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
		}
	}
	return list
}
