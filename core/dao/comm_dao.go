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
	"github.com/jsix/gof/db/orm"
	"go2o/core/dao/model"
	"gopkg.in/square/go-jose.v1/json"
	"log"
	"strings"
)

const (
	qrStoKey string = "go2o:comm:qr-templates"
)

type CommonDao struct {
	_orm orm.Orm
}

func NewCommDao(o orm.Orm) *CommonDao {
	return &CommonDao{_orm: o}
}

// 获取二维码所有模板
func (c *CommonDao) GetQrTemplates() []*model.CommQrTemplate {
	list := []*model.CommQrTemplate{}
	str, err := dSto.GetString(qrStoKey)
	if err == nil {
		err = json.Unmarshal([]byte(str), &list)
	}
	if err != nil {
		err = dOrm.Select(&list, "")
		if err == nil {
			d, _ := json.Marshal(list)
			dSto.Set(qrStoKey, string(d))
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
	_, err := orm.Save(dOrm, q, int(q.Id))
	if err == nil {
		dSto.Del(qrStoKey)
	}
	return err
}

// 删除二维码模板
func (c *CommonDao) DelQrTemplate(id int32) error {
	err := dOrm.DeleteByPk(model.CommQrTemplate{}, id)
	if err == nil {
		dSto.Del(qrStoKey)
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
