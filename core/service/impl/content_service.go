package impl

/**
 * Copyright 2015 @ to2.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

import (
	"context"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/content"
	"go2o/core/query"
	"go2o/core/service/proto"
)

var _ proto.ContentServiceServer = new(contentService)

type contentService struct {
	_contentRepo content.IContentRepo
	_query       *query.ContentQuery
	_sysContent  content.IContent
	serviceUtil
}

func NewContentService(rep content.IContentRepo, q *query.ContentQuery) *contentService {
	return &contentService{
		_contentRepo: rep,
		_query:       q,
		_sysContent:  rep.GetContent(0),
	}
}



// 获取页面
func (c *contentService) GetPage(_ context.Context, id *proto.IdOrName) (*proto.SPage, error) {
	ic := c._contentRepo.GetContent(0)
	var ia content.IPage
	if id.Id >0{
		ia = ic.GetPage(int32(id.Id))
	}else{
		ia = ic.GetPageByStringIndent(id.Name)
	}
	if ia != nil {
		return c.parsePageDto(ia.GetValue()),nil
	}
	return nil,nil
}


// 保存页面
func (c *contentService) SavePage(_ context.Context, v *proto.SPage) (*proto.Result, error) {
	ic := c._contentRepo.GetContent(0)
	var ip content.IPage
	var err error
	iv := c.parsePage(v)
	if v.Id > 0 {
		ip = ic.GetPage(int32(v.Id))
	} else {
		ip = ic.CreatePage(iv)
	}
	err = ip.SetValue(iv)
	if err == nil {
		_,err = ip.Save()
	}
	return c.error(err),nil
}


// 删除页面
func (c *contentService) DeletePage(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	ic := c._contentRepo.GetContent(0)
	err := ic.DeletePage(int32(id.Value))
	return c.error(err),nil
}


// 获取所有栏目
func (c *contentService) GetArticleCategories(_ context.Context, empty *proto.Empty) (*proto.ArticleCategoriesResponse, error) {
	list := c._sysContent.ArticleManager().GetAllCategory()
	arr := make([]*proto.SArticleCategory, len(list))
	for i, v := range list {
		val := v.GetValue()
		arr[i] = c.parseArticleCategoryDto(val)
	}
	return &proto.ArticleCategoriesResponse{
		Value:arr,
	},nil
}

// 获取栏目
func (c *contentService) GetArticleCategory(_ context.Context, name *proto.IdOrName) (*proto.SArticleCategory, error) {
	mgr := c._sysContent.ArticleManager()
	var ic content.ICategory
	if name.Id >0 {
		ic = mgr.GetCategory(int32(name.Id))
	} else {
		ic = mgr.GetCategoryByAlias(name.Name)
	}
	if c != nil {
		return c.parseArticleCategoryDto(ic.GetValue()),nil
	}
	return nil,nil
}


// 保存文章栏目
func (c *contentService) SaveArticleCategory(_ context.Context, r *proto.SArticleCategory) (*proto.Result, error) {
	m := c._sysContent.ArticleManager()
	ic := m.GetCategory(int32(r.Id))
	v := c.parseArticleCategory(r)
	if ic == nil {
		ic = m.CreateCategory(v)
	}
	err := ic.SetValue(v)
	if err == nil {
		_,err = ic.Save()
	}
	return c.error(err),nil
}


// 删除文章分类
func (c *contentService) DeleteArticleCategory(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := c._sysContent.ArticleManager().DelCategory(int32(id.Value))
	return c.error(err),nil
}


// 获取文章
func (c *contentService) GetArticle(_ context.Context, id *proto.IdOrName) (*proto.SArticle, error) {
	m :=  c._sysContent.ArticleManager()
	var ia content.IArticle
	if id.Id > 0{
		ia = m.GetArticle(int32(id.Id))
	}else{

	}
	if ia != nil {
		v := ia.GetValue()
		return c.parseArticleDto(&v),nil
	}
	return nil,nil
}

// 删除文章
func (c *contentService) DeleteArticle(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := c._sysContent.ArticleManager().DeleteArticle(int32(id.Value))
	return c.error(err),nil
}


// 保存文章
func (c *contentService) SaveArticle(_ context.Context, r *proto.SArticle) (*proto.Result, error) {
	m := c._sysContent.ArticleManager()
	v := c.parseArticle(r)
	var ia content.IArticle
	if r.Id > 0{
		ia = m.GetArticle(r.Id)
	}else{
		ia = m.CreateArticle(v)
	}
	err := ia.SetValue(v)
	if err == nil {
		_,err = ia.Save()
	}
	return c.error(err),nil
}



func (c *contentService) QueryPagingArticles(_ context.Context, r *proto.PagingArticleRequest) (*proto.ArticleListResponse, error) {
	var total = 0
	var rows []*content.Article
	ic := c._sysContent.ArticleManager().GetCategoryByAlias(r.CategoryName)
	if ic != nil && ic.GetDomainId() > 0 {
		total, rows = c._query.PagedArticleList(ic.GetDomainId(), int(r.Begin), int(r.Size), "")
	}
	var arr = make([]*proto.SArticle, 0)
	for _, v := range rows {
		arr = append(arr, c.parseArticleDto(v))
	}
	return &proto.ArticleListResponse{
		Total: int64(total),
		Data:  arr,
	},nil
}

func (c *contentService) QueryTopArticles(_ context.Context, cat *proto.IdOrName) (*proto.ArticleListResponse, error) {
	var arr = make([]*proto.SArticle, 0)
	var ic content.ICategory
	var m = c._sysContent.ArticleManager()
	if cat.Id > 0{
		ic = m.GetCategory(int32(cat.Id))
	}else {
		ic = m.GetCategoryByAlias(cat.Name)
	}
	if ic != nil && ic.GetDomainId() > 0 {
		_, rows := c._query.PagedArticleList(ic.GetDomainId(), 0, 1, "")
		for _, v := range rows {
			arr = append(arr, c.parseArticleDto(v))
		}
	}
	return &proto.ArticleListResponse{
		Total: 0,
		Data:  arr,
	},nil
}






func (c *contentService) parsePageDto(src *content.Page) *proto.SPage {
	return &proto.SPage{
		Id:                   int64(src.Id),
		UserId:               src.UserId,
		Title:                src.Title,
		StrIndent:            src.StrIndent,
		PermFlag:             int32(src.PermFlag),
		AccessKey:            src.AccessKey,
		KeyWord:              src.KeyWord,
		Description:          src.Description,
		CssPath:              src.CssPath,
		Body:                 src.Body,
		UpdateTime:           src.UpdateTime,
		Enabled:              src.Enabled == 1,
	}
}

func (c *contentService) parsePage(v *proto.SPage) *content.Page{
	return &content.Page{
		Id:          int32(v.Id),
		UserId:      v.UserId,
		Title:       v.Title,
		StrIndent:   v.StrIndent,
		PermFlag:    int(v.PermFlag),
		AccessKey:   v.AccessKey,
		KeyWord:     v.KeyWord,
		Description: v.Description,
		CssPath:     v.CssPath,
		Body:        v.Body,
		UpdateTime:  v.UpdateTime,
		Enabled:     types.IntCond(v.Enabled,1,0),
	}
}

func (c *contentService) parseArticleCategoryDto(v content.ArticleCategory) *proto.SArticleCategory {
	return &proto.SArticleCategory{
		Id:                   int64(v.ID),
		ParentId:             int64(v.ParentId)	,
		PermFlag:             int32(v.PermFlag),
		Name:                 v.Name,
		Alias:                v.Alias,
		SortNum:              int32(v.SortNum),
		Location:             v.Location,
		Title:                v.Title,
		Keywords:             v.Keywords,
		Description:          v.Description,
	}
}

func (c *contentService) parseArticleCategory(r *proto.SArticleCategory) *content.ArticleCategory{
	return &content.ArticleCategory{
		ID:          int32(r.Id),
		ParentId:    int32(r.ParentId),
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
		Id:          v.ID,
		CategoryId:       v.CatId,
		Title:       v.Title,
		SmallTitle:  v.SmallTitle,
		Thumbnail:   v.Thumbnail,
		PublisherId: int32(v.PublisherId),
		Location:    v.Location,
		Priority:    int32(v.Priority),
		AccessKey:   v.AccessKey,
		Content:     v.Content,
		Tags:        v.Tags,
		ViewCount:   int32(v.ViewCount),
		SortNum:     int32(v.SortNum),
		CreateTime:  int32(v.CreateTime),
		UpdateTime:  int32(v.UpdateTime),
	}
}

func (c *contentService) parseArticle(v *proto.SArticle) *content.Article {
	return &content.Article{
		ID:          v.Id,
		CatId:       v.CategoryId,
		Title:       v.Title,
		SmallTitle:  v.SmallTitle,
		Thumbnail:   v.Thumbnail,
		Location:    v.Location,
		Priority:    int(v.Priority),
		AccessKey:   v.AccessKey,
		PublisherId: int(v.PublisherId),
		Content:     v.Content,
		Tags:        v.Tags,
		ViewCount:  int( v.ViewCount),
		SortNum:    int( v.SortNum),
		CreateTime: int64( v.CreateTime),
		UpdateTime: int64( v.UpdateTime),
	}
}


