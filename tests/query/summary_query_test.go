/**
 * Copyright (C) 2007-2024 fze.NET,All rights reserved.
 *
 * name : summary_query_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2024-09-01 21:54
 * description : summary_query_test.go
 * history :
 */

package query

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

// 查询汇总信息
func TestQueryBoardSummary(t *testing.T) {
	qs := inject.GetStatisticsQueryService()
	s := qs.QuerySummary()
	t.Logf("%#v", s)
}
