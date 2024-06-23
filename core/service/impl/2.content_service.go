package impl

/**
 * Copyright 2015 @ 56x.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

import (
	"context"
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types"
)

var _ proto.ContentServiceServer = new(contentService)

type contentService struct {
	_contentRepo content.IArticleRepo
	_query       *query.ContentQuery
	_sysContent  content.IContentAggregateRoot
	serviceUtil
	proto.UnimplementedContentServiceServer
}

func NewContentService(rep content.IArticleRepo, q *query.ContentQuery) proto.ContentServiceServer {
	return &contentService{
		_contentRepo: rep,
		_query:       q,
		_sysContent:  rep.GetContent(0),
		serviceUtil:  serviceUtil{},
	}
}

// 获取页面
func (c *contentService) GetPage(_ context.Context, id *proto.IdOrName) (*proto.SPage, error) {
	ic := c._contentRepo.GetContent(0)
	var ia content.IPage
	if id.Id > 0 {
		ia = ic.PageManager().GetPage(int(id.Id))
	} else {
		ia = ic.PageManager().GetPageByCode(id.Name)
	}
	if ia != nil {
		return c.parsePageDto(ia.GetValue()), nil
	}
	return nil, fmt.Errorf("no such page: %v", id.Name)
}

// 保存页面
func (c *contentService) SavePage(_ context.Context, v *proto.SPage) (*proto.Result, error) {
	ic := c._contentRepo.GetContent(0)
	var ip content.IPage
	var err error
	iv := c.parsePage(v)
	if v.Id > 0 {
		ip = ic.PageManager().GetPage(int(v.Id))
	} else {
		ip = ic.PageManager().CreatePage(iv)
	}
	err = ip.SetValue(iv)
	if err == nil {
		_, err = ip.Save()
	}
	return c.error(err), nil
}

// 删除页面
func (c *contentService) DeletePage(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	ic := c._contentRepo.GetContent(0)
	err := ic.PageManager().DeletePage(int(id.Value))
	return c.error(err), nil
}

// 获取所有栏目
func (c *contentService) GetArticleCategories(_ context.Context, empty *proto.Empty) (*proto.ArticleCategoriesResponse, error) {
	list := c._sysContent.ArticleManager().GetAllCategory()
	arr := make([]*proto.SArticleCategory, len(list))
	for i, v := range list {
		arr[i] = c.parseArticleCategoryDto(v)
	}
	return &proto.ArticleCategoriesResponse{
		Value: arr,
	}, nil
}

// 获取栏目
func (c *contentService) GetArticleCategory(_ context.Context, name *proto.IdOrName) (*proto.SArticleCategory, error) {
	mgr := c._sysContent.ArticleManager()
	var ic *content.Category
	if name.Id > 0 {
		ic = mgr.GetCategory(int(name.Id))
	} else {
		ic = mgr.GetCategoryByAlias(name.Name)
	}
	if ic != nil {
		return c.parseArticleCategoryDto(*ic), nil
	}
	return nil, fmt.Errorf("no such category")
}

// 保存文章栏目
func (c *contentService) SaveArticleCategory(_ context.Context, r *proto.SArticleCategory) (*proto.Result, error) {
	m := c._sysContent.ArticleManager()
	v := c.parseArticleCategory(r)
	err := m.SaveCategory(v)
	return c.error(err), nil
}

// 删除文章分类
func (c *contentService) DeleteArticleCategory(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := c._sysContent.ArticleManager().DeleteCategory(int(id.Value))
	return c.error(err), nil
}

// GetArticle 获取文章
func (c *contentService) GetArticle(_ context.Context, id *proto.IdOrName) (*proto.SArticle, error) {
	m := c._sysContent.ArticleManager()
	var ia content.IArticle
	if id.Id > 0 {
		ia = m.GetArticle(int(id.Id))
	} else {

	}
	if ia != nil {
		v := ia.GetValue()
		return c.parseArticleDto(&v), nil
	}
	return nil, fmt.Errorf("no such article")
}

// DeleteArticle 删除文章
func (c *contentService) DeleteArticle(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	art := c._sysContent.ArticleManager().GetArticle(int(id.Value))
	if art == nil {
		return c.error(fmt.Errorf("no such article")), nil
	}
	err := art.Destory()
	return c.error(err), nil
}

// SaveArticle 保存文章
func (c *contentService) SaveArticle(_ context.Context, r *proto.SArticle) (*proto.Result, error) {
	m := c._sysContent.ArticleManager()
	v := c.parseArticle(r)
	var ia content.IArticle
	if r.Id > 0 {
		ia = m.GetArticle(int(r.Id))
	} else {
		ia = m.CreateArticle(v)
	}
	err := ia.SetValue(v)
	if err == nil {
		err = ia.Save()
	}
	return c.error(err), nil
}

// LikeArticle implements proto.ContentServiceServer.
func (c *contentService) LikeArticle(_ context.Context, req *proto.ArticleLikeRequest) (*proto.Result, error) {
	art := c._sysContent.ArticleManager().GetArticle(int(req.Id))
	if art == nil {
		return c.error(fmt.Errorf("no such article")), nil
	}
	var err error
	if req.IsDislike {
		err = art.Dislike(int(req.MemberId))
	}else{
		err = art.Like(int(req.MemberId))
	}
	return c.error(err), nil
}

// UpdateArticleViewsCount implements proto.ContentServiceServer.
func (c *contentService) UpdateArticleViewsCount(_ context.Context, req *proto.ArticleViewsRequest) (*proto.Result, error) {
	art := c._sysContent.ArticleManager().GetArticle(int(req.Id))
	if art == nil {
		return c.error(fmt.Errorf("no such article")), nil
	}
	err := art.IncreaseViewCount(int(req.MemberId), int(req.Count))
	return c.error(err), nil
}

func (c *contentService) QueryPagingArticles(_ context.Context, r *proto.PagingArticleRequest) (*proto.ArticleListResponse, error) {
	var total = 0
	var rows []*content.Article
	ic := c._sysContent.ArticleManager().GetCategoryByAlias(r.CategoryName)
	if ic != nil && ic.Id > 0 {
		total, rows = c._query.PagedArticleList(int32(ic.Id), int(r.Begin), int(r.Size), "")
	}
	var arr = make([]*proto.SArticle, 0)
	for _, v := range rows {
		arr = append(arr, c.parseArticleDto(v))
	}
	return &proto.ArticleListResponse{
		Total: int64(total),
		Data:  arr,
	}, nil
}

func (c *contentService) QueryTopArticles(_ context.Context, cat *proto.IdOrName) (*proto.ArticleListResponse, error) {
	var arr = make([]*proto.SArticle, 0)
	var ic *content.Category
	var m = c._sysContent.ArticleManager()
	if cat.Id > 0 {
		ic = m.GetCategory(int(cat.Id))
	} else {
		ic = m.GetCategoryByAlias(cat.Name)
	}
	if ic != nil && ic.Id > 0 {
		_, rows := c._query.PagedArticleList(int32(ic.Id), 0, 1, "")
		for _, v := range rows {
			arr = append(arr, c.parseArticleDto(v))
		}
	}
	return &proto.ArticleListResponse{
		Total: 0,
		Data:  arr,
	}, nil
}

func (c *contentService) parsePageDto(src *content.Page) *proto.SPage {
	return &proto.SPage{
		Id:          int64(src.Id),
		UserId:      int64(src.UserId),
		Title:       src.Title,
		Code:        src.Code,
		Flag:        int32(src.Flag),
		AccessKey:   src.AccessKey,
		KeyWord:     src.KeyWord,
		Description: src.Description,
		CssPath:     src.CssPath,
		Content:     src.Content,
		UpdateTime:  src.UpdateTime,
		Enabled:     src.Enabled == 1,
	}
}

func (c *contentService) parsePage(v *proto.SPage) *content.Page {
	return &content.Page{
		Id:          int(v.Id),
		UserId:      int(v.UserId),
		Title:       v.Title,
		Code:        v.Code,
		Flag:        int(v.Flag),
		AccessKey:   v.AccessKey,
		KeyWord:     v.KeyWord,
		Description: v.Description,
		CssPath:     v.CssPath,
		Content:     v.Content,
		UpdateTime:  v.UpdateTime,
		Enabled:     types.ElseInt(v.Enabled, 1, 0),
	}
}

func (c *contentService) parseArticleCategoryDto(v content.Category) *proto.SArticleCategory {
	return &proto.SArticleCategory{
		Id:          int64(v.Id),
		ParentId:    int64(v.ParentId),
		PermFlag:    int32(v.PermFlag),
		Name:        v.Name,
		Alias:       v.Alias,
		SortNum:     int32(v.SortNum),
		Location:    v.Location,
		Title:       v.Title,
		Keywords:    v.Keywords,
		Description: v.Description,
	}
}

func (c *contentService) parseArticleCategory(r *proto.SArticleCategory) *content.Category {
	return &content.Category{
		Id:          int(r.Id),
		ParentId:    int(r.ParentId),
		PermFlag:    int(r.PermFlag),
		Name:        r.Name,
		Alias:       r.Alias,
		SortNum:     int(r.SortNum),
		Location:    r.Location,
		Title:       r.Title,
		Keywords:    r.Keywords,
		Description: r.Description,
	}
}

func (c *contentService) parseArticleDto(v *content.Article) *proto.SArticle {
	return &proto.SArticle{
		Id:           int64(v.Id),
		CategoryId:   int64(v.CatId),
		Title:        v.Title,
		Flag:         int32(v.Flag),
		ShortTitle:   v.ShortTitle,
		Thumbnail:    v.Thumbnail,
		MchId:        int32(v.MchId),
		PublisherId:  int64(v.PublisherId),
		Location:     v.Location,
		Priority:     int32(v.Priority),
		AccessKey:    v.AccessToken,
		Content:      v.Content,
		Tags:         v.Tags,
		ViewCount:    int32(v.ViewCount),
		LikeCount:    int32(v.LikeCount),
		DislikeCount: int32(v.DislikeCount),
		SortNum:      int32(v.SortNum),
		CreateTime:   int64(v.CreateTime),
		UpdateTime:   int64(v.UpdateTime),
	}
}

func (c *contentService) parseArticle(v *proto.SArticle) *content.Article {
	return &content.Article{
		Id:           int(v.Id),
		CatId:        int(v.CategoryId),
		Title:        v.Title,
		ShortTitle:   v.ShortTitle,
		Flag:         int(v.Flag),
		Thumbnail:    v.Thumbnail,
		AccessToken:  v.AccessKey,
		PublisherId:  int(v.PublisherId),
		Location:     v.Location,
		Priority:     int(v.Priority),
		MchId:        int(v.MchId),
		Content:      v.Content,
		Tags:         v.Tags,
		LikeCount:    int(v.LikeCount),
		DislikeCount: int(v.DislikeCount),
		ViewCount:    int(v.ViewCount),
		SortNum:      int(v.SortNum),
		CreateTime:   int(v.CreateTime),
		UpdateTime:   int(v.UpdateTime),
	}
}
