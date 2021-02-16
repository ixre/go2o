package api

import "testing"

/**
 * Copyright 2009-2019 @ to2.net
 * name : res_api_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-26 17:22
 * description :
 * history :
 */

func TestAdApi(t *testing.T) {
	mp := map[string]string{
		"text":        "共产党是中华人民共和国的执政党",
		"replacement": "*",
	}
	testGET(t, "/fd/replace_sensitive", mp)
}
