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
	"errors"
	"github.com/jsix/gof/db/orm"
	"gopkg.in/square/go-jose.v1/json"
	"strings"
)

const (
	qrStoKey string = "go2o:comm:qr-templates"
)

type CommonDao struct {
}

// 获取二维码所有模板
func (c *CommonDao) GetQrTemplates() []*CommQrTemplate {
	list := []*CommQrTemplate{}
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
func (c *CommonDao) GetQrTemplate(id int) *CommQrTemplate {
	for _, v := range c.GetQrTemplates() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 保存二维码模板
func (c *CommonDao) SaveQrTemplate(q *CommQrTemplate) error {
	q.Title = strings.TrimSpace(q.Title)
	q.Comment = strings.TrimSpace(q.Comment)
	q.BgImage = strings.TrimSpace(q.BgImage)
	if q.Title == "" {
		return errors.New("标题不能为空")
	}
	if q.BgImage == "" {
		return errors.New("二维码背景图片为空")
	}
	_, err := orm.Save(dOrm, q, q.Id)
	if err == nil {
		dSto.Del(qrStoKey)
	}
	return err
}

// 删除二维码模板
func (c *CommonDao) DelQrTemplate(id int) error {
	err := dOrm.DeleteByPk(CommQrTemplate{}, id)
	if err == nil {
		dSto.Del(qrStoKey)
	}
	return err
}
