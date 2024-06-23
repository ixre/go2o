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

	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/inject"
)

func TestContentGetAllCategory(t *testing.T) {
	rep := inject.GetContentRepo()
	u := rep.GetContent(0)
	list := u.ArticleManager().GetAllCategory()
	for i, v := range list {
		t.Log("--", i, v.Name, v.Alias)
	}

	c := u.ArticleManager().GetCategoryByAlias("news")
	if c == nil {
		t.Log("栏目不存在")
		t.Fail()
	}
}

// 测试创建文章
func TestCreateArticle(t *testing.T) {
	rep := inject.GetContentRepo()
	u := rep.GetContent(0)
	iam := u.ArticleManager()
	ia := iam.CreateArticle(&content.Article{
		CatId:        1,
		Title:        "2024年，中国将迎来“世界互联网”新的发展期",
		ShortTitle:   "",
		Flag:         0,
		Thumbnail:    "https://tse2-mm.cn.bing.net/th/id/OIP-C.WrqPjA540HIh4D2eXfunrAHaEo?rs=1&pid=ImgDetMain",
		AccessToken:  "",
		PublisherId:  0,
		Location:     "",
		Priority:     0,
		MchId:        0,
		Content:      "内容",
		Tags:         "",
		LikeCount:    0,
		DislikeCount: 0,
		ViewCount:    0,
		SortNum:      0,
	})
	ia.Save()
}
