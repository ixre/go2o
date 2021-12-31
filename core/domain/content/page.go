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
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/tmp"
	"time"
)

var _ content.IPage = new(pageImpl)

type pageImpl struct {
	contentRepo content.IArchiveRepo
	userId      int32
	value       *content.Page
}

func newPage(userId int32, rep content.IArchiveRepo,
	v *content.Page) content.IPage {
	return &pageImpl{
		contentRepo: rep,
		userId:      userId,
		value:       v,
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
	total := 0
	tmp.Db().ExecScalar("SELECT COUNT(0) FROM arc_page WHERE user_id= $1 AND str_indent= $2 AND id <> $3",
		&total, p.userId, alias, p.GetDomainId())
	return total == 0
}

// SetValue 设置值
func (p *pageImpl) SetValue(v *content.Page) error {
	v.Id = p.GetDomainId()
	if p.value.UserId != v.UserId {
		return content.ErrUserNotMatch
	}
	if p.value.Flag & content.FlagInternal == content.FlagInternal {
		if p.value.Code != v.Code{
			return content.ErrInternalPage
		}
	}
	if len(v.Code) > 0 && !p.checkAliasExists(v.Code) {
		return content.ErrAliasHasExists
	}
	p.value = v
	return nil
}

// Save 保存
func (p *pageImpl) Save() (int32, error) {
	p.value.UpdateTime = time.Now().Unix()
	return p.contentRepo.SavePage(p.userId, p.value)
}
