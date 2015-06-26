/**
 * Copyright 2015 @ S1N1 Team.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps
import (
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/query"
)


type contentService struct {
	_contentRep content.IContentRep
	_query     *query.ContentQuery
}

func NewContentService(rep content.IContentRep, q *query.ContentQuery) *contentService {
	return &contentService{
		_contentRep: rep,
		_query:     q,
	}
}