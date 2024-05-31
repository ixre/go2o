/**
 * Copyright 2015 @ 56x.net.
 * name : content_test
 * author : jarryliu
 * date : 2016-07-06 10:23
 * description :
 * history :
 */
package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

func TestContentGetAllCategory(t *testing.T) {
	rep := inject.GetContentRepo()
	u := rep.GetContent(0)
	list := u.ArticleManager().GetAllCategory()
	for i, c := range list {
		v := c.GetValue()
		t.Log("--", i, v.Name, v.Alias)
	}

	c := u.ArticleManager().GetCategoryByAlias("mall-activity")
	if c == nil {
		t.Log("栏目不存在")
		t.Fail()
	}
}
