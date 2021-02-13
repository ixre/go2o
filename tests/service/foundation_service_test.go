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

func TestCheckSensitive (t *testing.T)  {
	ret,_ := impl.FoundationService.ReplaceSensitive(context.TODO(),
		&proto.ReplaceSensitiveRequest{
			Word:                 "共产党是我们唯一的领导政府",
			Replacement:          "",
		})
	t.Log(ret.Value)
}