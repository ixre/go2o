/**
 * Copyright 2015 @ z3q.net.
 * name : pag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import (
	"go2o/core/domain/interface/content"
	"time"
)

var _ content.IPage = new(pageImpl)

type pageImpl struct {
	_contentRep content.IContentRep
	_merchantId int
	_value      *content.Page
}

func newPage(merchantId int, rep content.IContentRep,
	v *content.Page) content.IPage {
	return &pageImpl{
		_contentRep: rep,
		_merchantId: merchantId,
		_value:      v,
	}
}

// 获取领域编号
func (p *pageImpl) GetDomainId() int {
	return p._value.Id
}

// 获取值
func (p *pageImpl) GetValue() *content.Page {
	return p._value
}

// 设置值
func (p *pageImpl) SetValue(v *content.Page) error {
	v.Id = p.GetDomainId()
	if p._value.UserId != v.UserId {
		return content.ErrUserNotMatch
	}
	p._value = v
	return nil
}

// 保存
func (p *pageImpl) Save() (int, error) {
	p._value.UpdateTime = time.Now().Unix()
	return p._contentRep.SavePage(p._merchantId, p._value)
}
