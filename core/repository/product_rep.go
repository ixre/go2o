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
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/sale/product"
	"go2o/core/infrastructure/format"
)

var _ product.IProductRepo = new(productRepo)

type productRepo struct {
	db.Connector
}

func NewProductRepo(c db.Connector) product.IProductRepo {
	return &productRepo{
		Connector: c,
	}
}

func (p *productRepo) GetProductValue(itemId int32) *product.Product {
	var e *product.Product = new(product.Product)
	//todo: supplier_id  == -1
	if p.Connector.GetOrm().GetByQuery(e, `select * FROM pro_product
			INNER JOIN cat_category c ON c.id = pro_product.cat_id
			 WHERE pro_product.id=?`, itemId) == nil {
		return e
	}
	return nil
}

func (p *productRepo) GetProductsById(ids ...int32) ([]*product.Product, error) {
	//todo: mchId
	var items []*product.Product

	//todo:改成database/sql方式，不使用orm
	err := p.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM pro_product WHERE id IN (`+format.IdArrJoinStr32(ids)+`)`)

	return items, err
}

func (p *productRepo) SaveProductValue(v *product.Product) (int32, error) {
	return orm.I32(orm.Save(p.GetOrm(), v, int(v.Id)))
}

func (p *productRepo) GetPagedOnShelvesProduct(mchId int32, catIds []int32,
	start, end int) (total int, e []*product.Product) {
	var sql string

	var catIdStr string = format.IdArrJoinStr32(catIds)
	sql = fmt.Sprintf(`SELECT * FROM pro_product INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		WHERE merchant_id=%d AND cat_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, mchId, catIdStr, start, (end - start))

	p.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM pro_product INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		WHERE merchant_id=%d AND cat_category.id IN (%s) AND on_shelves=1`, mchId, catIdStr), &total)

	e = []*product.Product{}
	p.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (p *productRepo) GetProductSaleNum(mchId int32, id int32) int {
	var num int
	p.Connector.ExecScalar(`SELECT SUM(sale_num) FROM gs_goods WHERE item_id=?`,
		&num, id)
	return num
}

func (p *productRepo) DeleteProduct(mchId, itemId int32) error {
	_, _, err := p.Connector.Exec(`
		DELETE f FROM pro_product AS f
		INNER JOIN cat_category AS c ON f.cat_id=c.id
		WHERE f.id=? AND c.merchant_id=?`, itemId, mchId)
	return err
}
