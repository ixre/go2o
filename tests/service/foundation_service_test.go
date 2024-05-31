package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
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
	ret, _ := inject.GetFoundationService().ReplaceSensitive(context.TODO(),
		&proto.ReplaceSensitiveRequest{
			Text:        "我自愿加入中国共产党,坚持党的领导,守护我们的长城",
			Replacement: "*",
		})
	t.Log(ret.Value)
}
