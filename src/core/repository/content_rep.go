/**
 * Copyright 2015 @ S1N1 Team.
 * name : content_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository
import (
	"go2o/src/core/domain/interface/content"
	"github.com/atnet/gof/db"
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