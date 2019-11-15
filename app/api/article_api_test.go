package api

import (
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : article_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-09-04 18:13
 * description :
 * history :
 */

func TestPassportApi_Top_Article(t *testing.T) {
	mp := map[string]string{
		"cat": "news",
	}
	testApi(t, "article.top_article", mp, true)
}
