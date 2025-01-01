/**
 * Copyright 2015 @ 56x.net.
 * name : pag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import (
	"time"

	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/infrastructure/format"
)

var _ content.IPageManager = new(pageManagerImpl)

type pageManagerImpl struct {
	tenantId int
	pageRepo content.IPageRepo
}

// CreatePage 创建页面
func (c *pageManagerImpl) CreatePage(v *content.Page) content.IPage {
	return NewPage(c.tenantId, c.pageRepo, v)
}

// GetPage 获取页面
func (c *pageManagerImpl) GetPage(id int) content.IPage {
	return c.pageRepo.GetPageById(c.tenantId, id)
}

// GetPageByCode 根据字符串标识获取页面
func (c *pageManagerImpl) GetPageByCode(indent string) content.IPage {
	return c.pageRepo.GetPageByCode(c.tenantId, indent)
}

// DeletePage 删除页面
func (c *pageManagerImpl) DeletePage(id int) error {
	ip := c.pageRepo.GetPageById(c.tenantId, id)
	if ip == nil {
		return content.ErrNoSuchPage
	}
	if ip.GetValue().Flag&content.FCategoryInternal == content.FCategoryInternal {
		return content.ErrInternalPage
	}
	return c.pageRepo.DeletePage(c.tenantId, id)
}

var _ content.IPage = new(pageImpl)

type pageImpl struct {
	repo     content.IPageRepo
	tenantId int
	value    *content.Page
}

func NewPage(tenantId int, repo content.IPageRepo,
	v *content.Page) content.IPage {
	return &pageImpl{
		repo:     repo,
		tenantId: tenantId,
		value:    v,
	}
}

// GetDomainId 获取领域编号
func (p *pageImpl) GetDomainId() int {
	return p.value.Id
}

// GetValue 获取值
func (p *pageImpl) GetValue() *content.Page {
	return p.value
}

// 检测别名是否可用
func (p *pageImpl) checkAliasExists(alias string) bool {
	v := p.repo.FindBy("user_id = ? AND code = ? AND id <> ?", p.tenantId, alias, p.GetDomainId())
	return v != nil
}

// SetValue 设置值
func (p *pageImpl) SetValue(v *content.Page) error {
	v.Id = p.GetDomainId()
	if p.value.UserId != v.UserId {
		return content.ErrUserNotMatch
	}
	if p.value.Flag&content.FCategoryInternal == content.FCategoryInternal {
		if p.value.Code != v.Code {
			return content.ErrInternalPage
		}
	}
	if len(v.Code) > 0 && p.checkAliasExists(v.Code) {
		return content.ErrAliasHasExists
	}

	v.Content = format.RemoveHtmlStyle(v.Content)
	p.value = v
	return nil
}

// Save 保存
func (p *pageImpl) Save() (int, error) {
	p.value.UpdateTime = int(time.Now().Unix())
	err := p.repo.SavePage(p.tenantId, p.value)
	return p.GetDomainId(), err
}
