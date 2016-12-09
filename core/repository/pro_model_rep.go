package repository

import (
	"database/sql"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/pro_model"
	pmImpl "go2o/core/domain/pro_model"
	"go2o/core/infrastructure/format"
	"log"
)

var _ promodel.IProModelRepo = new(proModelRepo)

type proModelRepo struct {
	_orm         orm.Orm
	conn         db.Connector
	brandService promodel.IBrandService
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
	return pmImpl.NewModel(v, p)
}

// 获取商品模型
func (p *proModelRepo) GetModel(id int32) promodel.IModel {
	v := p.GetProModel(id)
	if v != nil {
		return p.CreateModel(v)
	}
	return nil
}

//获取品牌服务
func (p *proModelRepo) BrandService() promodel.IBrandService {
	if p.brandService == nil {
		p.brandService = pmImpl.NewBrandService(p)
	}
	return p.brandService
}

// 设置产品模型的品牌
func (p *proModelRepo) SetModelBrands(proModel int32, brandIds []int32) error {
	idArrStr := format.IdArrJoinStr32(brandIds)
	//获取存在的品牌
	old := p.SelectProModelBrand("pro_model=?", proModel)
	//删除不包括的品牌
	if len(old) > 0 {
		p.BatchDeleteProModelBrand("pro_model = ? AND brand_id NOT IN(?)",
			proModel, idArrStr)
	}
	//写入品牌
	for _, v := range brandIds {
		isExist := false
		for _, vo := range old {
			if vo.BrandId == v {
				isExist = true
				break
			}
		}
		if isExist {
			e := &promodel.ProModelBrand{
				Id:       0,
				BrandId:  v,
				ProModel: proModel,
			}
			p.SaveProModelBrand(e)
		}
	}
	return nil
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
	id, err := orm.Save(p._orm, v, int(v.Id))
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

// Batch Delete ProModel
func (p *proModelRepo) BatchDeleteProModel(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(promodel.ProModel{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ProModel")
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
	id, err := orm.Save(p._orm, v, int(v.Id))
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
	id, err := orm.Save(p._orm, v, int(v.Id))
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
