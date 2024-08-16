/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type ContentQuery struct {
	db.Connector
	o orm.Orm
	fw.BaseRepository[content.Article]
	categoryRepo fw.BaseRepository[content.Category]
	pageRepo     fw.BaseRepository[content.Page]
	mq           *MerchantQuery
	mmq          *MemberQuery
}

func NewContentQuery(o orm.Orm, fo fw.ORM, mq *MerchantQuery, mmq *MemberQuery) *ContentQuery {
	c := &ContentQuery{
		Connector: o.Connector(),
		o:         o,
		mq:        mq,
		mmq:       mmq,
	}
	c.ORM = fo
	c.categoryRepo.ORM = fo
	c.pageRepo.ORM = fo
	c.BaseRepository.ORM = fo
	return c
}

func (c *ContentQuery) PagedArticleList(p *fw.PagingParams) (ret *fw.PagingResult, err error) {
	ret, err = c.QueryPaging(p)
	var mchIds []int
	var memberIds []int
	for _, v := range ret.Rows {
		r := v.(*content.Article)
		if r.MchId > 0 {
			mchIds = append(mchIds, r.MchId)
		}
		if r.PublisherId > 0 {
			memberIds = append(memberIds, r.PublisherId)
		}
	}
	var mchMap map[int]*merchant.Merchant
	var mmMap map[int]*member.Member
	if len(mchIds) > 0 {
		mchMap = collections.ToMap(c.mq.FindList(nil, "id IN ?", mchIds), func(m *merchant.Merchant) (int, *merchant.Merchant) {
			return m.Id, m
		})
	}
	if len(memberIds) > 0 {
		mmMap = collections.ToMap(c.mmq.FindList(nil, "id IN ?", mchIds), func(m *member.Member) (int, *member.Member) {
			return int(m.Id), m
		})
	}
	retArray := make([]interface{}, len(ret.Rows))
	for i, v := range ret.Rows {
		r := v.(*content.Article)
		dst := &PagingArticleDto{
			Id:           r.Id,
			CatId:        r.CatId,
			Title:        r.Title,
			ShortTitle:   r.ShortTitle,
			Flag:         r.Flag,
			Thumbnail:    r.Thumbnail,
			PublisherId:  r.PublisherId,
			Location:     r.Location,
			Priority:     r.Priority,
			MchId:        r.MchId,
			Tags:         r.Tags,
			LikeCount:    r.LikeCount,
			DislikeCount: r.DislikeCount,
			ViewCount:    r.ViewCount,
			CreateTime:   r.CreateTime,
			UpdateTime:   r.UpdateTime,
			Ext:          map[string]interface{}{},
		}
		if m, ok := mchMap[r.MchId]; ok {
			dst.Ext["mchName"] = m.MchName
		}
		if m, ok := mmMap[r.PublisherId]; ok {
			dst.Ext["publisherName"] = m.Nickname
		}
		retArray[i] = dst
	}
	ret.Rows = retArray
	return ret, err
}

type PagingArticleDto struct {
	Id           int                    `json:"id"`
	CatId        int                    `json:"catId"`
	Title        string                 `json:"title"`
	ShortTitle   string                 `json:"shortTitle"`
	Flag         int                    `json:"flag"`
	Thumbnail    string                 `json:"thumbnail"`
	PublisherId  int                    `json:"publisherId"`
	Location     string                 `json:"location"`
	Priority     int                    `json:"priority"`
	MchId        int                    `json:"mchId"`
	Tags         string                 `json:"tags"`
	LikeCount    int                    `json:"likeCount"`
	DislikeCount int                    `json:"dislikeCount"`
	ViewCount    int                    `json:"viewCount"`
	SortNum      int                    `json:"sortNum"`
	CreateTime   int                    `json:"createTime"`
	UpdateTime   int                    `json:"updateTime"`
	Ext          map[string]interface{} `json:"ext"`
}

// QueryMerchantArticles 查询商户文章列表
func (c *ContentQuery) QueryMerchantArticles(mchId int, p *fw.PagingParams) (*fw.PagingResult, error) {
	p.Equal("mch_id", mchId)
	return c.BaseRepository.QueryPaging(p)
}

func (c *ContentQuery) QueryPagingArticleCategories(p *fw.PagingParams) (*fw.PagingResult, error) {
	return c.categoryRepo.QueryPaging(p)
}

// 系统查询文章列表
func (c *ContentQuery) QueryPagingArticles(p *fw.PagingParams) (*fw.PagingResult, error) {
	return c.BaseRepository.QueryPaging(p)
}

// 系统查询页面列表
func (c *ContentQuery) QueryPagingPages(p *fw.PagingParams) (*fw.PagingResult, error) {
	return c.pageRepo.QueryPaging(p)
}

// 查询分类编号列表
func (c *ContentQuery) QueryCategoryIdList(s []string) []int {
	arr := c.categoryRepo.FindList(nil, "alias IN ?", s)
	return collections.MapList(arr, func(v *content.Category) int {
		return v.Id
	})
}

func (c *ContentQuery) QueryRecommendArticles(p *fw.PagingParams) (*fw.PagingResult, error) {
	//.. 推荐条件
	row, err := c.BaseRepository.QueryPaging(p)
	for _, row := range row.Rows {
		art := row.(*content.Article)
		if art != nil {
			art.Content = ""
		}
	}
	return row, err
}
