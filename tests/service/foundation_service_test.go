package service

import (
	"context"
	"go2o/core/service/impl"
	"go2o/core/service/proto"
	"testing"
)

/**
 * Copyright (C) 2007-2021 56X.NET,All rights reserved.
 *
 * name : foundation_service_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2021-02-13 23:13
 * description :
 * history :
 */

func TestCheckSensitive(t *testing.T) {
	ret, _ := impl.FoundationService.ReplaceSensitive(context.TODO(),
		&proto.ReplaceSensitiveRequest{
			Word:        "我自愿加入中国共产党,坚持党的领导,守护我们的长城",
			Replacement: "*",
		})
	t.Log(ret.Value)
}
