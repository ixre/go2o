package api

import (
	"github.com/ixre/gof/api"
	"go2o/core/service/rsi"
)

var _ api.Handler = new(ArticleApi)

type ArticleApi struct {
}

func (a ArticleApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"list": a.list,
	})
}

// 文章列表
func (a ArticleApi) list(ctx api.Context) interface{} {
	form := ctx.Form()
	catStr := form.GetString("cat")
	if len(catStr) == 0 {
		return api.NewErrorResponse("缺少参数:cat")
	}
	page := form.GetInt("page")
	size := form.GetInt("size")
	begin := (page - 1) * size
	total, rows := rsi.ContentService.PagedArticleList(catStr, begin, size, "")
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}
