package repos

import (
	"database/sql"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/pro_model"
	pmImpl "go2o/core/domain/pro_model"
	"log"
)

var _ promodel.IProModelRepo = new(proModelRepo)

type proModelRepo struct {
	_orm         orm.Orm
	conn         db.Connector
	brandService promodel.IBrandService
	attrService  promodel.IAttrService
	specService  promodel.ISpecService
}

// Create new ProBrandRepo
func NewProModelRepo(conn db.Connector, o orm.Orm) promodel.IProModelRepo {
	return &proModelRepo{
		_orm: o,
		conn: conn,
	}
}

// 创建商品模型
func (p *proModelRepo) CreateModel(v *promodel.ProModel) promodel.IModel {
	return pmImpl.NewModel(v, p, p.AttrService(), p.SpecService(),
		p.BrandService())
}

// 获取商品模型
func (p *proModelRepo) GetModel(id int32) promodel.IModel {
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

//获取品牌服务
func (p *proModelRepo) BrandService() promodel.IBrandService {
	if p.brandService == nil {
		p.brandService = pmImpl.NewBrandService(p)
	}
	return p.brandService
}

// 获取模型的商品品牌
func (p *proModelRepo) GetModelBrands(proModel int32) []*promodel.ProBrand {
	return p.selectProBrandByQuery(`SELECT * FROM pro_brand WHERE id IN (
	SELECT brand_id FROM pro_model_brand WHERE pro_model=?)`, proModel)
}

// Get ProModel
func (p *proModelRepo) GetProModel(primary interface{}) *promodel.ProModel {
	e := promodel.ProModel{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModel")
	}
	return nil
}

// Select ProModel
func (p *proModelRepo) SelectProModel(where string, v ...interface{}) []*promodel.ProModel {
	list := []*promodel.ProModel{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModel")
	}
	return list
}

// Save ProModel
func (p *proModelRepo) SaveProModel(v *promodel.ProModel) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModel")
	}
	return id, err
}

// Delete ProModel
func (p *proModelRepo) DeleteProModel(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.ProModel{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModel")
	}
	return err
}

// Get Attr
func (p *proModelRepo) GetAttr(primary interface{}) *promodel.Attr {
	e := promodel.Attr{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attr")
	}
	return nil
}

// Select Attr
func (p *proModelRepo) SelectAttr(where string, v ...interface{}) []*promodel.Attr {
	list := []*promodel.Attr{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attr")
	}
	return list
}

// Save Attr
func (p *proModelRepo) SaveAttr(v *promodel.Attr) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attr")
	}
	return id, err
}

// Delete Attr
func (p *proModelRepo) DeleteAttr(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.Attr{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attr")
	}
	return err
}

// Batch Delete Attr
func (p *proModelRepo) BatchDeleteAttr(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.Attr{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Attr")
	}
	return r, err
}

// Get AttrItem
func (p *proModelRepo) GetAttrItem(primary interface{}) *promodel.AttrItem {
	e := promodel.AttrItem{}
	err := p._orm.Get(primary, &e)
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
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return list
}

// Save AttrItem
func (p *proModelRepo) SaveAttrItem(v *promodel.AttrItem) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return id, err
}

// Delete AttrItem
func (p *proModelRepo) DeleteAttrItem(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.AttrItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return err
}

// Batch Delete AttrItem
func (p *proModelRepo) BatchDeleteAttrItem(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.AttrItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AttrItem")
	}
	return r, err
}

// Get Spec
func (p *proModelRepo) GetSpec(primary interface{}) *promodel.Spec {
	e := promodel.Spec{}
	err := p._orm.Get(primary, &e)
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
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return list
}

// Save Spec
func (p *proModelRepo) SaveSpec(v *promodel.Spec) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return id, err
}

// Delete Spec
func (p *proModelRepo) DeleteSpec(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.Spec{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return err
}

// Batch Delete Spec
func (p *proModelRepo) BatchDeleteSpec(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.Spec{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Spec")
	}
	return r, err
}

// Get SpecItem
func (p *proModelRepo) GetSpecItem(primary interface{}) *promodel.SpecItem {
	e := promodel.SpecItem{}
	err := p._orm.Get(primary, &e)
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
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return list
}

// Save SpecItem
func (p *proModelRepo) SaveSpecItem(v *promodel.SpecItem) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return id, err
}

// Delete SpecItem
func (p *proModelRepo) DeleteSpecItem(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.SpecItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return err
}

// Batch Delete SpecItem
func (p *proModelRepo) BatchDeleteSpecItem(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.SpecItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SpecItem")
	}
	return r, err
}

// Get ProBrand
func (p *proModelRepo) GetProBrand(primary interface{}) *promodel.ProBrand {
	e := promodel.ProBrand{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return nil
}

// Save ProBrand
func (p *proModelRepo) SaveProBrand(v *promodel.ProBrand) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return id, err
}

// Delete ProBrand
func (p *proModelRepo) DeleteProBrand(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.ProBrand{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return err
}

// Batch Delete ProBrand
func (p *proModelRepo) BatchDeleteProBrand(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.ProBrand{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return r, err
}

// Select ProBrand
func (p *proModelRepo) SelectProBrand(where string, v ...interface{}) []*promodel.ProBrand {
	list := []*promodel.ProBrand{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return list
}

// Select ProBrand
func (p *proModelRepo) selectProBrandByQuery(query string, v ...interface{}) []*promodel.ProBrand {
	list := []*promodel.ProBrand{}
	err := p._orm.SelectByQuery(&list, query, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProBrand")
	}
	return list
}

// Get ProModelBrand
func (p *proModelRepo) GetProModelBrand(primary interface{}) *promodel.ProModelBrand {
	e := promodel.ProModelBrand{}
	err := p._orm.Get(primary, &e)
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
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return id, err
}

// Delete ProModelBrand
func (p *proModelRepo) DeleteProModelBrand(primary interface{}) error {
	err := p._orm.DeleteByPk(promodel.ProModelBrand{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return err
}

// Batch Delete ProModelBrand
func (p *proModelRepo) BatchDeleteProModelBrand(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.ProModelBrand{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return r, err
}

// Select ProModelBrand
func (p *proModelRepo) SelectProModelBrand(where string, v ...interface{}) []*promodel.ProModelBrand {
	list := []*promodel.ProModelBrand{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModelBrand")
	}
	return list
}
