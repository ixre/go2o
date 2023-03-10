package repos

import (
	"database/sql"
	"log"

	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	pmImpl "github.com/ixre/go2o/core/domain/pro_model"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ promodel.IProductModelRepo = new(proModelRepo)

type proModelRepo struct {
	conn         db.Connector
	brandService promodel.IBrandService
	attrService  promodel.IAttrService
	specService  promodel.ISpecService
	o            orm.Orm
}

// Create new ProBrandRepo
func NewProModelRepo(o orm.Orm) promodel.IProductModelRepo {
	return &proModelRepo{
		o:    o,
		conn: o.Connector(),
	}
}

// 创建商品模型
func (p *proModelRepo) CreateModel(v *promodel.ProductModel) promodel.IProductModel {
	return pmImpl.NewModel(v, p, p.AttrService(), p.SpecService(),
		p.BrandService())
}

// 获取商品模型
func (p *proModelRepo) GetModel(id int) promodel.IProductModel {
	v := p.GetProModel(id)
	if v != nil {
		return p.CreateModel(v)
	}
	return nil
}

// 属性服务
func (p *proModelRepo) AttrService() promodel.IAttrService {
	if p.attrService == nil {
		p.attrService = pmImpl.NewAttrService(p)
	}
	return p.attrService
}

// 规格服务
func (p *proModelRepo) SpecService() promodel.ISpecService {
	if p.specService == nil {
		p.specService = pmImpl.NewSpecService(p)
	}
	return p.specService
}

// 获取品牌服务
func (p *proModelRepo) BrandService() promodel.IBrandService {
	if p.brandService == nil {
		p.brandService = pmImpl.NewBrandService(p)
	}
	return p.brandService
}

// IsExistsBrand 是否存在相同名称的品牌
func (p *proModelRepo) IsExistsBrand(name string, id int) bool {
	var row int
	p.o.Connector().ExecScalar(`SELECT COUNT(1) FROM product_brand WHERE name = $1 AND id <> $2`,
		&row, name, id)
	return row > 0
}

// 获取模型的商品品牌
func (p *proModelRepo) GetModelBrands(proModel int) []*promodel.ProductBrand {
	return p.selectProBrandByQuery(`SELECT * FROM product_brand WHERE id IN (
	SELECT brand_id FROM product_model_brand WHERE prod_model= $1)`, proModel)
}

// Get ProductModel
func (p *proModelRepo) GetProModel(primary interface{}) *promodel.ProductModel {
	e := promodel.ProductModel{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductModel")
	}
	return nil
}

// Select ProductModel
func (p *proModelRepo) SelectProModel(where string, v ...interface{}) []*promodel.ProductModel {
	var list []*promodel.ProductModel
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductModel")
	}
	return list
}

// Save ProductModel
func (p *proModelRepo) SaveProModel(v *promodel.ProductModel) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductModel")
	}
	return id, err
}

// Delete ProductModel
func (p *proModelRepo) DeleteProModel(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.ProductModel{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductModel")
	}
	return err
}

// Get Attrs
func (p *proModelRepo) GetAttr(primary interface{}) *promodel.Attr {
	e := promodel.Attr{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attrs")
	}
	return nil
}

// Select Attrs
func (p *proModelRepo) SelectAttr(where string, v ...interface{}) []*promodel.Attr {
	list := []*promodel.Attr{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attrs")
	}
	return list
}

// Save Attrs
func (p *proModelRepo) SaveAttr(v *promodel.Attr) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attrs")
	}
	return id, err
}

// Delete Attrs
func (p *proModelRepo) DeleteAttr(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.Attr{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attrs")
	}
	return err
}

// Batch Delete Attrs
func (p *proModelRepo) BatchDeleteAttr(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.Attr{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attrs")
	}
	return r, err
}

// Get AttrItem
func (p *proModelRepo) GetAttrItem(primary interface{}) *promodel.AttrItem {
	e := promodel.AttrItem{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return nil
}

// Select AttrItem
func (p *proModelRepo) SelectAttrItem(where string, v ...interface{}) []*promodel.AttrItem {
	list := []*promodel.AttrItem{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return list
}

// Save AttrItem
func (p *proModelRepo) SaveAttrItem(v *promodel.AttrItem) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return id, err
}

// Delete AttrItem
func (p *proModelRepo) DeleteAttrItem(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.AttrItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return err
}

// Batch Delete AttrItem
func (p *proModelRepo) BatchDeleteAttrItem(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.AttrItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return r, err
}

// Get Spec
func (p *proModelRepo) GetSpec(primary interface{}) *promodel.Spec {
	e := promodel.Spec{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return nil
}

// Select Spec
func (p *proModelRepo) SelectSpec(where string, v ...interface{}) []*promodel.Spec {
	list := []*promodel.Spec{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return list
}

// Save Spec
func (p *proModelRepo) SaveSpec(v *promodel.Spec) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return id, err
}

// Delete Spec
func (p *proModelRepo) DeleteSpec(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.Spec{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return err
}

// Batch Delete Spec
func (p *proModelRepo) BatchDeleteSpec(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.Spec{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return r, err
}

// Get SpecItem
func (p *proModelRepo) GetSpecItem(primary interface{}) *promodel.SpecItem {
	e := promodel.SpecItem{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return nil
}

// Select SpecItem
func (p *proModelRepo) SelectSpecItem(where string, v ...interface{}) []*promodel.SpecItem {
	list := []*promodel.SpecItem{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return list
}

// Save SpecItem
func (p *proModelRepo) SaveSpecItem(v *promodel.SpecItem) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return id, err
}

// Delete SpecItem
func (p *proModelRepo) DeleteSpecItem(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.SpecItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return err
}

// Batch Delete SpecItem
func (p *proModelRepo) BatchDeleteSpecItem(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.SpecItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return r, err
}

// Get ProductBrand
func (p *proModelRepo) GetProBrand(primary interface{}) *promodel.ProductBrand {
	e := promodel.ProductBrand{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return nil
}

// Save ProductBrand
func (p *proModelRepo) SaveProBrand(v *promodel.ProductBrand) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return id, err
}

// Delete ProductBrand
func (p *proModelRepo) DeleteProBrand(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.ProductBrand{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return err
}

// Batch Delete ProductBrand
func (p *proModelRepo) BatchDeleteProBrand(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.ProductBrand{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return r, err
}

// Select ProductBrand
func (p *proModelRepo) SelectProBrand(where string, v ...interface{}) []*promodel.ProductBrand {
	list := []*promodel.ProductBrand{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return list
}

// Select ProductBrand
func (p *proModelRepo) selectProBrandByQuery(query string, v ...interface{}) []*promodel.ProductBrand {
	list := []*promodel.ProductBrand{}
	err := p.o.SelectByQuery(&list, query, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProductBrand")
	}
	return list
}

// Get ProModelBrand
func (p *proModelRepo) GetProModelBrand(primary interface{}) *promodel.ProModelBrand {
	e := promodel.ProModelBrand{}
	err := p.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return nil
}

// Save ProModelBrand
func (p *proModelRepo) SaveProModelBrand(v *promodel.ProModelBrand) (int, error) {
	id, err := orm.Save(p.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return id, err
}

// Delete ProModelBrand
func (p *proModelRepo) DeleteProModelBrand(primary interface{}) error {
	err := p.o.DeleteByPk(promodel.ProModelBrand{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return err
}

// Batch Delete ProModelBrand
func (p *proModelRepo) BatchDeleteProModelBrand(where string, v ...interface{}) (int64, error) {
	r, err := p.o.Delete(promodel.ProModelBrand{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return r, err
}

// Select ProModelBrand
func (p *proModelRepo) SelectProModelBrand(where string, v ...interface{}) []*promodel.ProModelBrand {
	list := []*promodel.ProModelBrand{}
	err := p.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return list
}
