/**
 * Copyright 2015 @ z3q.net.
 * name : content_test
 * author : jarryliu
 * date : 2016-07-06 10:23
 * description :
 * history :
 */
package testing

import (
	"go2o/core/repository"
	"go2o/core/testing/include"
	"testing"
)

func TestContentGetAllCategory(t *testing.T) {
	rep := repository.NewContentRepo(include.GetApp().Db())
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
