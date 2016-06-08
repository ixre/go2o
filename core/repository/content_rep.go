/**
 * Copyright 2015 @ z3q.net.
 * name : content_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	contentImpl "go2o/core/domain/content"
	"go2o/core/domain/interface/content"
)

var _ content.IContentRep = new(contentRep)

type contentRep struct {
	db.Connector
}

// 内容仓储
func NewContentRep(c db.Connector) content.IContentRep {
	return &contentRep{
		Connector: c,
	}
}

// 获取内容
func (this *contentRep) GetContent(merchantId int) content.IContent {
	return contentImpl.NewContent(merchantId, this)
}

// 根据编号获取页面
func (this *contentRep) GetPageById(merchantId, id int) *content.ValuePage {
	var e content.ValuePage
	if err := this.Connector.GetOrm().Get(id, &e); err == nil && e.MerchantId == merchantId {
		return &e
	}
	return nil
}

// 根据标识获取页面
func (this *contentRep) GetPageByStringIndent(merchantId int, indent string) *content.ValuePage {
	var e content.ValuePage
	if err := this.Connector.GetOrm().GetBy(&e, "merchant_id=? and str_indent=?", merchantId, indent); err == nil {
		return &e
	}
	return nil
}

// 删除页面
func (this *contentRep) DeletePage(merchantId, id int) error {
	_, err := this.Connector.GetOrm().Delete(content.ValuePage{}, "merchant_id=? AND id=?", merchantId, id)
	return err
}

// 保存页面
func (this *contentRep) SavePage(merchantId int, v *content.ValuePage) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM mch_page WHERE merchant_id=?", &v.Id, merchantId)
	}
	return v.Id, err
}
