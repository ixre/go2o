/**
 * Copyright 2015 @ at3.net.
 * name : comm_dao.go
 * author : jarryliu
 * date : 2016-11-15 19:54
 * description :
 * history :
 */
package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core/dao/model"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/product"
	"gopkg.in/square/go-jose.v1/json"
	"log"
	"strings"
)

const (
	qrStoKey string = "go2o:comm:qr-templates"
)

type CommonDao struct {
	_orm    orm.Orm
	storage storage.Interface
	adRepo  ad.IAdRepo
	catRepo product.ICategoryRepo
}

func NewCommDao(o orm.Orm, sto storage.Interface,
	adRepo ad.IAdRepo, catRepo product.ICategoryRepo) *CommonDao {
	return &CommonDao{
		_orm:    o,
		storage: sto,
		adRepo:  adRepo,
		catRepo: catRepo,
	}
}

// 获取二维码所有模板
func (c *CommonDao) GetQrTemplates() []*model.CommQrTemplate {
	list := []*model.CommQrTemplate{}
	str, err := c.storage.GetString(qrStoKey)
	if err == nil {
		err = json.Unmarshal([]byte(str), &list)
	}
	if err != nil {
		err = c._orm.Select(&list, "")
		if err == nil {
			d, _ := json.Marshal(list)
			c.storage.Set(qrStoKey, string(d))
		}
	}
	return list
}

// 获取二维码模板
func (c *CommonDao) GetQrTemplate(id int32) *model.CommQrTemplate {
	for _, v := range c.GetQrTemplates() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 保存二维码模板
func (c *CommonDao) SaveQrTemplate(q *model.CommQrTemplate) error {
	q.Title = strings.TrimSpace(q.Title)
	q.Comment = strings.TrimSpace(q.Comment)
	q.BgImage = strings.TrimSpace(q.BgImage)
	if q.Title == "" {
		return errors.New("标题不能为空")
	}
	if q.BgImage == "" {
		return errors.New("二维码背景图片为空")
	}
	_, err := orm.Save(c._orm, q, int(q.Id))
	if err == nil {
		c.storage.Del(qrStoKey)
	}
	return err
}

// 删除二维码模板
func (c *CommonDao) DelQrTemplate(id int32) error {
	err := c._orm.DeleteByPk(model.CommQrTemplate{}, id)
	if err == nil {
		c.storage.Del(qrStoKey)
	}
	return err
}

// Get PortalNavType
func (p *CommonDao) GetPortalNavType(primary interface{}) *model.PortalNavType {
	e := model.PortalNavType{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNavType")
	}
	return nil
}

// Select PortalNavType
func (p *CommonDao) SelectPortalNavType(where string, v ...interface{}) []*model.PortalNavType {
	list := []*model.PortalNavType{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNavType")
	}
	return list
}

// Save PortalNavType
func (p *CommonDao) SavePortalNavType(v *model.PortalNavType) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNavType")
	}
	return id, err
}

// Delete PortalNavType
func (p *CommonDao) DeletePortalNavType(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PortalNavType{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNavType")
	}
	return err
}

// Batch Delete PortalNavType
func (p *CommonDao) BatchDeletePortalNavType(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PortalNavType{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNavType")
	}
	return r, err
}

// Get PortalNav
func (p *CommonDao) GetPortalNav(primary interface{}) *model.PortalNav {
	e := model.PortalNav{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNav")
	}
	return nil
}

// Select PortalNav
func (p *CommonDao) SelectPortalNav(where string, v ...interface{}) []*model.PortalNav {
	list := []*model.PortalNav{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNav")
	}
	return list
}

// Save PortalNav
func (p *CommonDao) SavePortalNav(v *model.PortalNav) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNav")
	}
	return id, err
}

// Delete PortalNav
func (p *CommonDao) DeletePortalNav(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PortalNav{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNav")
	}
	return err
}

// Batch Delete PortalNav
func (p *CommonDao) BatchDeletePortalNav(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PortalNav{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalNav")
	}
	return r, err
}

func (p *CommonDao) GetFloorAdPos(catId int32) string {
	key := fmt.Sprintf("go2o:portal:floor-ad:%d", catId)
	r, err := p.storage.GetString(key)
	if err != nil {
		e := model.PortalFloorAd{}
		err := p._orm.GetBy(&e, "cat_id=?", catId)
		if err == nil {
			pos := p.adRepo.GetAdPositionById(e.PosId)
			if pos != nil {
				r = pos.Key
			}
		}
		p.storage.Set(key, r)
	}
	return r
}

func (p *CommonDao) SetFloorAd(catId int32, posId int32) (err error) {
	cat := p.catRepo.GetCategory(0, catId)
	if cat == nil {
		err = product.ErrNoSuchCategory
	} else if cat.FloorShow != 1 {
		err = errors.New("商品分类设置楼层显示")
	}
	e := model.PortalFloorAd{}
	p._orm.GetBy(&e, "cat_id=?", catId)
	e.CatId = catId
	e.PosId = posId
	e.AdIndex = 0
	_, err = p.SavePortalFloorAd(&e)
	if err == nil {
		p.storage.Del(fmt.Sprintf("go2o:portal:floor-ad:%d", catId))
	}
	return err
}

// Get PortalFloorAd
func (p *CommonDao) GetPortalFloorAd(primary interface{}) *model.PortalFloorAd {
	e := model.PortalFloorAd{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorAd")
	}
	return nil
}

// Save PortalFloorAd
func (p *CommonDao) SavePortalFloorAd(v *model.PortalFloorAd) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorAd")
	}
	return id, err
}

// Batch Delete PortalFloorAd
func (p *CommonDao) BatchDeletePortalFloorAd(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PortalFloorAd{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorAd")
	}
	return r, err
}

// Get PortalFloorLink
func (p *CommonDao) GetPortalFloorLink(primary interface{}) *model.PortalFloorLink {
	e := model.PortalFloorLink{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorLink")
	}
	return nil
}

// Select PortalFloorLink
func (p *CommonDao) SelectPortalFloorLink(where string, v ...interface{}) []*model.PortalFloorLink {
	list := []*model.PortalFloorLink{}
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorLink")
	}
	return list
}

// Save PortalFloorLink
func (p *CommonDao) SavePortalFloorLink(v *model.PortalFloorLink) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorLink")
	}
	return id, err
}

// Delete PortalFloorLink
func (p *CommonDao) DeletePortalFloorLink(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PortalFloorLink{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorLink")
	}
	return err
}

// Batch Delete PortalFloorLink
func (p *CommonDao) BatchDeletePortalFloorLink(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PortalFloorLink{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PortalFloorLink")
	}
	return r, err
}
