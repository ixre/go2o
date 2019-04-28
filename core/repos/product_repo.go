/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:09
 * description :
 * history :
 */

package repos

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	proImpl "go2o/core/domain/product"
	"go2o/core/infrastructure/format"
	"log"
)

var _ product.IProductRepo = new(productRepo)

type productRepo struct {
	db.Connector
	_orm      orm.Orm
	pmRepo    promodel.IProModelRepo
	valueRepo valueobject.IValueRepo
}

func NewProductRepo(c db.Connector, pmRepo promodel.IProModelRepo,
	valRepo valueobject.IValueRepo) product.IProductRepo {
	return &productRepo{
		Connector: c,
		_orm:      c.GetOrm(),
		pmRepo:    pmRepo,
		valueRepo: valRepo,
	}
}

// 创建产品
func (p *productRepo) CreateProduct(v *product.Product) product.IProduct {
	return proImpl.NewProductImpl(v, p, p.pmRepo, p.valueRepo)
}

// 根据产品编号获取货品
func (p *productRepo) GetProduct(itemId int64) product.IProduct {
	v := p.GetProductValue(itemId)
	if v != nil {
		return p.CreateProduct(v)
	}
	return nil
}

func (p *productRepo) GetProductsById(ids ...int32) ([]*product.Product, error) {
	//todo: mchId
	var items []*product.Product

	//todo:改成database/sql方式，不使用orm
	err := p.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM pro_product WHERE id IN (`+format.I32ArrStrJoin(ids)+`)`)

	return items, err
}

func (p *productRepo) GetPagedOnShelvesProduct(mchId int32, catIds []int32,
	start, end int) (total int, e []*product.Product) {
	var sql string

	var catIdStr string = format.I32ArrStrJoin(catIds)
	sql = fmt.Sprintf(`SELECT * FROM pro_product INNER JOIN pro_category ON pro_product.cat_id=pro_category.id
		WHERE merchant_id=%d AND pro_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, mchId, catIdStr, start, (end - start))

	p.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM pro_product INNER JOIN pro_category ON pro_product.cat_id=pro_category.id
		WHERE merchant_id=%d AND pro_category.id IN (%s) AND on_shelves=1`, mchId, catIdStr), &total)

	e = []*product.Product{}
	p.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (p *productRepo) GetProductSaleNum(itemId int64) int {
	var num int
	p.Connector.ExecScalar(`SELECT SUM(sale_num) FROM item_info WHERE product_id=?`,
		&num, itemId)
	return num
}

// Get Product
func (p *productRepo) GetProductValue(itemId int64) *product.Product {
	e := product.Product{}
	err := p._orm.Get(itemId, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Product")
	}
	return nil
}

// Select Product
func (p *productRepo) SelectProduct(where string, v ...interface{}) []*product.Product {
	list := []*product.Product{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Product")
	}
	return list
}

// Save Product
func (p *productRepo) SaveProduct(v *product.Product) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Product")
	}
	return id, err
}

// Delete Product
func (p *productRepo) DeleteProduct(itemId int64) error {
	err := p._orm.DeleteByPk(product.Product{}, itemId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Product")
	}
	return err
}

// Batch Delete Product
func (p *productRepo) BatchDeleteProduct(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(product.Product{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Product")
	}
	return r, err
}

// Get ProAttrInfo
func (p *productRepo) GetAttr(primary interface{}) *product.Attr {
	e := product.Attr{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProAttrInfo")
	}
	return nil
}

// Select ProAttrInfo
func (p *productRepo) SelectAttr(where string, v ...interface{}) []*product.Attr {
	list := []*product.Attr{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProAttrInfo")
	}
	return list
}

// Save ProAttrInfo
func (p *productRepo) SaveAttr(v *product.Attr) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProAttrInfo")
	}
	return id, err
}

// Delete ProAttrInfo
func (p *productRepo) DeleteAttr(primary interface{}) error {
	err := p._orm.DeleteByPk(product.Attr{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProAttrInfo")
	}
	return err
}

// Batch Delete ProAttrInfo
func (p *productRepo) BatchDeleteAttr(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(product.Attr{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProAttrInfo")
	}
	return r, err
}
