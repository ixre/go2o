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
	"go2o/core/domain/tmp"
	"time"
)

var _ content.IPage = new(pageImpl)

type pageImpl struct {
	contentRep content.IContentRep
	userId     int
	value      *content.Page
}

func newPage(userId int, rep content.IContentRep,
	v *content.Page) content.IPage {
	return &pageImpl{
		contentRep: rep,
		userId:     userId,
		value:      v,
	}
}

// 获取领域编号
func (p *pageImpl) GetDomainId() int {
	return p.value.Id
}

// 获取值
func (p *pageImpl) GetValue() *content.Page {
	return p.value
}

// 检测别名是否可用
func (a *pageImpl) checkAliasExists(alias string) bool {
	total := 0
	tmp.Db().ExecScalar("SELECT COUNT(0) FROM con_page WHERE user_id=? AND str_indent=? AND id<>?",
		&total, a.userId, alias, a.GetDomainId())
	return total == 0
}

// 设置值
func (p *pageImpl) SetValue(v *content.Page) error {
	v.Id = p.GetDomainId()
	if p.value.UserId != v.UserId {
		return content.ErrUserNotMatch
	}
	if len(v.StrIndent) > 0 && !p.checkAliasExists(v.StrIndent) {
		return content.ErrAliasHasExists
	}
	p.value = v
	return nil
}

// 保存
func (p *pageImpl) Save() (int, error) {
	p.value.UpdateTime = time.Now().Unix()
	return p.contentRep.SavePage(p.userId, p.value)
}
