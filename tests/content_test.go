/**
 * Copyright 2015 @ to2.net.
 * name : content_test
 * author : jarryliu
 * date : 2016-07-06 10:23
 * description :
 * history :
 */
package tests

import (
	"go2o/core/repos"
	"go2o/tests/ti"
	"testing"
)

func TestContentGetAllCategory(t *testing.T) {
	rep := repos.NewContentRepo(ti.GetApp().Db())
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
