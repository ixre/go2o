package api

import (
	"github.com/ixre/gof/api"
	"go2o/core/service/auto_gen/rpc/content_service"
	"go2o/core/service/thrift"
)

var _ api.Handler = new(ArticleApi)

type ArticleApi struct {
}

func (a ArticleApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"list":        a.list,
		"top_article": a.topArticle,
	})
}

/**
 * @api {post} /article/list 文章列表
 * @apiName list
 * @apiGroup article
 * @apiParam {String} cat 栏目别名
 * @apiParam {Int} page 页码
 * @apiParam {Int} size 数量
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (a ArticleApi) list(ctx api.Context) interface{} {
	form := ctx.Form()
	catStr := form.GetString("cat")
	if len(catStr) == 0 {
		return api.NewErrorResponse("缺少参数:cat")
	}
	page := form.GetInt("page")
	size := form.GetInt("size")
	begin := (page - 1) * size
	trans, cli, err := thrift.ContentServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.QueryPagingArticles(thrift.Context, catStr, int32(begin), int32(size))
		return r
	}
	return map[string]interface{}{
		"total": 0,
		"rows":  []*content_service.SArticle{},
	}
}


/**
 * @api {post} /article/top_article 获取置顶文章
 * @apiName top_article
 * @apiGroup article
 * @apiParam {String} cat 栏目别名
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (a ArticleApi) topArticle(ctx api.Context) interface{} {
	form := ctx.Form()
	catStr := form.GetString("cat")
	if len(catStr) == 0 {
		return api.NewErrorResponse("缺少参数:cat")
	}
	trans, cli, err := thrift.ContentServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.QueryTopArticles(thrift.Context, catStr)
		return r
	}
	return []*content_service.SArticle{}
}
