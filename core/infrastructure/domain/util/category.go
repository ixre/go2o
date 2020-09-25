/**
 * Copyright 2015 @ to2.net.
 * name : category
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"github.com/ixre/gof/algorithm/iterator"
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/product"
	"go2o/core/service/proto"
)

type CategoryFormatFunc func(c *product.Category, level int)

// 迭代栏目
func IterateCategory(categories []*product.Category, c *product.Category,
	iterateFunc CategoryFormatFunc, finishFunc CategoryFormatFunc, level int) {
	if c.Id != 0 {
		iterateFunc(c, level)
	}

	for _, k := range categories {
		if k.ParentId == c.Id {
			IterateCategory(categories, k, iterateFunc, finishFunc, level+1)
		}
	}

	if finishFunc != nil {
		finishFunc(c, level)
	}

}

// 迭代栏目
func WalkSaleCategory(cs []*product.Category, v *product.Category,
	start iterator.WalkFunc, over iterator.WalkFunc) {
	var condition iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*product.Category).ParentId == v.(*product.Category).Id
	}
	var arr = make([]interface{}, len(cs))
	for i, v := range cs {
		arr[i] = v
	}
	iterator.Walk(arr, v, condition, start, over, 0)
}

// 迭代栏目
func WalkArticleCategory(cs *proto.ArticleCategoriesResponse, v *proto.SArticleCategory, start iterator.WalkFunc, over iterator.WalkFunc) {
	var condition iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*content.ArticleCategory).ParentId ==
			v.(*content.ArticleCategory).ID
	}
	var arr = make([]interface{}, len(cs.Value))
	for i, v := range cs.Value {
		arr[i] = v
	}
	iterator.Walk(arr, v, condition, start, over, 0)
}
