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
	"go2o/src/core/domain/interface/content"
	"time"
)

var _ content.IPage = new(Page)

type Page struct {
	_contentRep content.IContentRep
	_merchantId int
	_value      *content.ValuePage
}

func NewPage(merchantId int, rep content.IContentRep, v *content.ValuePage) content.IPage {
	return &Page{
		_contentRep: rep,
		_merchantId: merchantId,
		_value:      v,
	}
}

// 获取领域编号
func (this *Page) GetDomainId() int {
	return this._value.Id
}

// 获取值
func (this *Page) GetValue() *content.ValuePage {
	return this._value
}

// 设置值
func (this *Page) SetValue(v *content.ValuePage) error {
	v.Id = this.GetDomainId()
	this._value = v
	return nil
}

// 保存
func (this *Page) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	return this._contentRep.SavePage(this._merchantId, this._value)
}
