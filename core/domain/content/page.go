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
func (this *pageImpl) GetDomainId() int {
	return this._value.Id
}

// 获取值
func (this *pageImpl) GetValue() *content.Page {
	return this._value
}

// 设置值
func (this *pageImpl) SetValue(v *content.Page) error {
	v.Id = this.GetDomainId()
	this._value = v
	return nil
}

// 保存
func (this *pageImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	return this._contentRep.SavePage(this._merchantId, this._value)
}
